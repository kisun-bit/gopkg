package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/kisunSea/gopkg/runtime/traceback"
)

type Level = zapcore.Level

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = iota - 1
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel
)

type Encoder string

const (
	// EncodeJson is a fast, low-allocation JSON encoder.
	// The encoder appropriately escapes all field keys and values.
	EncodeJson Encoder = "json"
	// EncodeConsole is an encoder whose output is designed for human
	EncodeConsole Encoder = "console"
)

type Logger struct {
	baseLogger *zap.Logger
	sLogger    *zap.SugaredLogger
	config_    zapcore.EncoderConfig
	level_     zapcore.Level
	format_    Encoder
	handlers   []handler
}

////////////////////////////////////////////

// InitConfig returns `EncoderConfig`
func InitConfig(format string) zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(format))
		},
	}
}

// NewLumberjackFileRotatingLogger returns instance of `*lumberjack.baseLogger`
func NewLumberjackFileRotatingLogger(level Level, file string, maxSize, maxBackups, maxAge int) *lumberjack.Logger {
	if level > DebugLevel {
		// rename log file name
		_ext := filepath.Ext(filepath.Base(file))
		_suffix := "_" + level.String() + _ext
		file = strings.TrimRight(file, _ext) + _suffix
	}
	return &lumberjack.Logger{Filename: file, MaxSize: maxSize, MaxBackups: maxBackups, MaxAge: maxAge, Compress: false}
}

// MkdirAllUtilSuccess. at the specified number of retries, the folder is created until it succeeds
func MkdirAllUtilSuccess(dir string, retryTimes int) (err error) {
	for i := 0; i < retryTimes; i++ {
		if _, err = os.Stat(dir); err != nil {
			// lack of dir. create it with `-p`
			if err = os.MkdirAll(dir, 0666); err != nil {
				continue
			}
		} else {
			return nil
		}
	}
	return fmt.Errorf("make dirs with `-p` arg, reach the max retry times(%d)", retryTimes)
}

// NewLogger returns logger and flushFunc for clearing, params description:
// `level`: Lowest level of log output
// `stackTrace`: Stack trace log level
// `prefix`: Prefix on each line to identify the logger
// `timeFormat`: The time format of each line in the log
// `color`: True when color is enabled
// `encoder`: Log encoding format, divided into `EncodeJson` and `EncodeConsole`, default `EncodeConsole`
// `writer`: The type is going to be either `*lumberjack.baseLogger` or `nil`,
//           when set to nil, it will be output to the stdout.
func NewLogger(
	level,
	stackTrace Level,
	prefix,
	timeFormat string,
	color bool,
	encoder Encoder,
	writers ...interface{}) (logger *Logger, err error) {

	logger = new(Logger)
	logger.level_ = level

	if timeFormat == "" {
		timeFormat = "2006/01/02 - 15:04:05.000"
	}
	if prefix != "" {
		prefix = " " + prefix
	}
	logger.config_ = InitConfig(timeFormat + prefix)

	if stackTrace == Level(-1) {
		logger.config_.StacktraceKey = ""
	}
	if !color {
		logger.config_.EncodeLevel = zapcore.CapitalLevelEncoder
	}
	if err = logger.setWriters(writers); err != nil {
		return nil, err
	}

	logger.format_ = encoder
	logger.baseLogger = zap.New(
		logger.GetAndBuildCore(logger.config_),
		zap.AddStacktrace(stackTrace))

	logger.baseLogger = logger.baseLogger.WithOptions(zap.AddCaller())
	logger = logger.sugared()

	return logger, nil
}

func (l *Logger) sugared() *Logger {
	l.baseLogger = l.baseLogger.WithOptions(zap.AddCallerSkip(1)) // Trace back 1 level of the call stack
	l.sLogger = l.baseLogger.Sugar()
	return l
}

func (l *Logger) initHandlers() {
	l.handlers = make([]handler, 0)
}

func (l *Logger) addHandler(sync zapcore.WriteSyncer, lowLevel zapcore.Level) {
	l.handlers = append(l.handlers, handler{
		Sync:       sync,
		EnableFunc: func(lev zapcore.Level) bool { return lev >= lowLevel },
	})
}

func (l *Logger) setWriters(writers []interface{}) (err error) {

	defer func() {
		if len(l.handlers) == 0 {
			l.addHandler(zapcore.AddSync(os.Stdout), DebugLevel)
		}
	}()

	for _, writer := range writers {
		if err = l.bindHandler(writer); err != nil {
			return err
		}
	}
	return nil
}

func (l *Logger) bindHandler(writer interface{}) (err error) {

	var (
		sync__     zapcore.WriteSyncer
		lowLevel__ Level
	)

	switch i := writer.(type) {
	case *rotateWriter:
		if err = MkdirAllUtilSuccess(filepath.Dir(i.LogSavePath), 10); err != nil {
			return err
		}
		// lumberjack.baseLogger is already safe for concurrent use, so we don't need to lock it.
		lumberJackLogger := NewLumberjackFileRotatingLogger(i.Level, i.LogSavePath, i.MaxSize, i.MaxBackups, i.MaxAge)
		sync__, lowLevel__ = zapcore.AddSync(lumberJackLogger), i.Level
		zapcore.Lock(sync__)
		break
	case *consoleWriter:
		sync__, lowLevel__ = zapcore.AddSync(os.Stdout), i.Level
		break
	default:
		return fmt.Errorf("unsupported writer: %T", i)
	}

	l.addHandler(sync__, lowLevel__)
	return nil
}

// GetAndBuildCore
func (l *Logger) GetAndBuildCore(config zapcore.EncoderConfig) zapcore.Core {

	if len(l.handlers) == 0 {
		panic("why handlers is nil ???")
	}

	var cores []zapcore.Core
	for _, handler := range l.handlers {
		var tmpEncoder zapcore.Encoder
		switch l.format_ {
		case EncodeJson:
			tmpEncoder = zapcore.NewJSONEncoder(config)
		case EncodeConsole:
			fallthrough
		default:
			tmpEncoder = zapcore.NewConsoleEncoder(config)
		}

		cores = append(cores,
			zapcore.NewCore(tmpEncoder, zapcore.NewMultiWriteSyncer(handler.Sync), handler.EnableFunc))
	}

	return zapcore.NewTee(cores...)
}

func (l *Logger) CloseStacktrace() *Logger {
	c := l.config_
	c.StacktraceKey = ""
	l.baseLogger = l.WrapCore(c)
	l.sugared()
	return l
}

func (l *Logger) SetStacktrace(level Level) *Logger {
	l.baseLogger = l.baseLogger.WithOptions(zap.AddStacktrace(level))
	return l
}

func (l *Logger) NoColor() *Logger {
	c := l.config_
	c.EncodeLevel = zapcore.CapitalLevelEncoder
	l.baseLogger = l.WrapCore(c)
	return l
}

// SetTimeFormat sets the log output format.
// default time format is `2006/01/02 - 15:04:05.000`,
func (l *Logger) SetTimeFormat(timeFormat string) *Logger {
	c := l.config_
	c.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(timeFormat))
	}

	l.baseLogger = l.WrapCore(c)
	l.sugared()
	return l
}

// WrapCore wraps or replaces the SugarLogger's underlying zapcore.Core.
func (l *Logger) WrapCore(ec zapcore.EncoderConfig) *zap.Logger {
	return l.baseLogger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return l.GetAndBuildCore(ec)
	}))
}

func (l *Logger) Sync() (err error) {
	return l.sLogger.Sync()
}

// Debug logs messages at DEBUG level
func (l *Logger) Debug(args ...interface{}) {
	l.sLogger.Debug(args...)
}

// Info logs messages at INFO level
func (l *Logger) Info(args ...interface{}) {
	l.sLogger.Info(args...)
}

// Warn logs messages at WARN level
func (l *Logger) Warn(args ...interface{}) {
	l.sLogger.Warn(args...)
}

// Error logs messages at ERROR level
func (l *Logger) Error(args ...interface{}) {
	l.sLogger.Error(args...)
}

// Fatal logs messages at FATAL level
func (l *Logger) Fatal(args ...interface{}) {
	l.sLogger.Fatal(args...)
}

// Panic logs messages at Panic level
func (l *Logger) Panic(args ...interface{}) {
	l.sLogger.Panic(args...)
}

// DPanic logs messages at DPanic level
func (l *Logger) DPanic(args ...interface{}) {
	l.sLogger.DPanic(args...)
}

// DebugW logs messages at DEBUG level
func (l *Logger) DebugW(msg string, keysAndValues ...interface{}) {
	l.sLogger.Debugw(msg, keysAndValues...)
}

// InfoW logs messages at INFO level
func (l *Logger) InfoW(msg string, keysAndValues ...interface{}) {
	l.sLogger.Infow(msg, keysAndValues...)
}

// WarnW logs messages at WARN level
func (l *Logger) WarnW(msg string, keysAndValues ...interface{}) {
	l.sLogger.Warnw(msg, keysAndValues...)
}

// ErrorW logs messages at ERROR level
func (l *Logger) ErrorW(msg string, keysAndValues ...interface{}) {
	l.sLogger.Errorw(msg, keysAndValues...)
}

// FatalW logs messages at FATAL level
func (l *Logger) FatalW(msg string, keysAndValues ...interface{}) {
	l.sLogger.Fatalw(msg, keysAndValues...)
}

// PanicW logs messages at Panic level
func (l *Logger) PanicW(msg string, keysAndValues ...interface{}) {
	l.sLogger.Panicw(msg, keysAndValues...)
}

// DPanicW logs messages at DPanic level
func (l *Logger) DPanicW(msg string, keysAndValues ...interface{}) {
	l.sLogger.DPanicw(msg, keysAndValues...)
}

// DebugF logs messages at DEBUG level
func (l *Logger) DebugF(format string, args ...interface{}) {
	l.sLogger.Debugf(format, args...)
}

// InfoF logs messages at INFO level
func (l *Logger) InfoF(format string, args ...interface{}) {
	l.sLogger.Infof(format, args...)
}

// WarnF logs messages at WARN level
func (l *Logger) WarnF(format string, args ...interface{}) {
	l.sLogger.Warnf(format, args...)
}

// ErrorF logs messages at ERROR level
func (l *Logger) ErrorF(format string, args ...interface{}) {
	l.sLogger.Errorf(format, args...)
}

// Fatalf logs messages at FATAL level
func (l *Logger) FatalF(format string, args ...interface{}) {
	l.sLogger.Fatalf(format, args...)
}

// PanicF logs messages at Panic level
func (l *Logger) PanicF(format string, args ...interface{}) {
	l.sLogger.Panicf(format, args...)
}

// DPanicF logs messages at DPanic level
func (l *Logger) DPanicF(format string, args ...interface{}) {
	l.sLogger.DPanicf(format, args...)
}

// Stack returns the function call stack
func (l *Logger) Stack() string {
	return traceback.TakeStacktrace(1)
}
