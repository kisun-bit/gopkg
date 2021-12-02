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

// Log is used for logging messages
type Log interface {
	// Debug logs messages at DEBUG level
	Debug(args ...interface{})
	// Info logs messages at INFO level
	Info(args ...interface{})
	// Warn logs messages at WARN level
	Warn(args ...interface{})
	// Error logs messages at ERROR level
	Error(args ...interface{})
	// Fatal logs messages at FATAL level
	Fatal(args ...interface{})
	// Panic logs messages at Panic level
	Panic(args ...interface{})
	// DPanic logs messages at DPanic level
	DPanic(args ...interface{})
	// DebugF logs messages at DEBUG level
	DebugF(format string, args ...interface{})
	// InfoF logs messages at INFO level
	InfoF(format string, args ...interface{})
	// WarnF logs messages at WARN level
	WarnF(format string, args ...interface{})
	// ErrorF logs messages at ERROR level
	ErrorF(format string, args ...interface{})
	// FatalF logs messages at FATAL level
	FatalF(format string, args ...interface{})
	// PanicF logs messages at Panic level
	PanicF(format string, args ...interface{})
	// DPanicF logs messages at DPanic level
	DPanicF(format string, args ...interface{})
	// DebugW logs messages at DEBUG level
	DebugW(msg string, keysAndValues ...interface{})
	// InfoW logs messages at INFO level
	InfoW(msg string, keysAndValues ...interface{})
	// WarnW logs messages at WARN level
	WarnW(msg string, keysAndValues ...interface{})
	// ErrorW logs messages at ERROR level
	ErrorW(msg string, keysAndValues ...interface{})
	// FatalW logs messages at FATAL level
	FatalW(msg string, keysAndValues ...interface{})
	// PanicW logs messages at Panic level
	PanicW(msg string, keysAndValues ...interface{})
	// DPanicW logs messages at DPanic level
	DPanicW(msg string, keysAndValues ...interface{})
	// stack
	Stack() string
}

type Logger struct {
	Logger   *zap.Logger
	SLogger  *zap.SugaredLogger
	Config_  zapcore.EncoderConfig
	Level_   zapcore.Level
	Format_  Encoder
	Handlers []Handler
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

// NewLumberjackFileRotatingLogger returns instance of `*lumberjack.Logger`
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
// `writer`: The type is going to be either `*lumberjack.Logger` or `nil`,
//           when set to nil, it will be output to the stdout.
func NewLogger(
	level,
	stackTrace Level,
	prefix,
	timeFormat string,
	color bool,
	encoder Encoder,
	writers ...interface{}) (logger Log, flushFunc func() error, err error) {

	l_ := new(Logger)
	l_.Level_ = level

	if timeFormat == "" {
		timeFormat = "2006/01/02 - 15:04:05.000"
	}
	if prefix != "" {
		prefix = " " + prefix
	}
	l_.Config_ = InitConfig(timeFormat + prefix)

	if stackTrace == Level(-1) {
		l_.Config_.StacktraceKey = ""
	}
	if !color {
		l_.Config_.EncodeLevel = zapcore.CapitalLevelEncoder
	}
	if err = l_.SetWriters(writers); err != nil {
		return nil, nil, err
	}

	l_.Format_ = encoder
	l_.Logger = zap.New(
		l_.GetAndBuildCore(l_.Config_),
		zap.AddStacktrace(stackTrace))

	l_.Logger = l_.Logger.WithOptions(zap.AddCaller())
	l_ = l_.sugared()

	return l_, l_.Logger.Sync, nil
}

func (l *Logger) sugared() *Logger {
	l.Logger = l.Logger.WithOptions(zap.AddCallerSkip(1)) // Trace back 1 level of the call stack
	l.SLogger = l.Logger.Sugar()
	return l
}

func (l *Logger) InitHandlers() {
	l.Handlers = make([]Handler, 0)
}

func (l *Logger) AddHandler(sync zapcore.WriteSyncer, lowLevel zapcore.Level) {
	l.Handlers = append(l.Handlers, Handler{
		Sync:       sync,
		EnableFunc: func(lev zapcore.Level) bool { return lev >= lowLevel },
	})
}

func (l *Logger) SetWriters(writers []interface{}) (err error) {

	defer func() {
		if len(l.Handlers) == 0 {
			l.AddHandler(zapcore.AddSync(os.Stdout), DebugLevel)
		}
	}()

	for _, writer := range writers {
		if err = l.BindHandler(writer); err != nil {
			return err
		}
	}
	return nil
}

func (l *Logger) BindHandler(writer interface{}) (err error) {

	var (
		sync__     zapcore.WriteSyncer
		lowLevel__ Level
	)

	switch i := writer.(type) {
	case *RotateWriter:
		if err = MkdirAllUtilSuccess(filepath.Dir(i.LogSavePath), 10); err != nil {
			return err
		}
		// lumberjack.Logger is already safe for concurrent use, so we don't need to lock it.
		lumberJackLogger := NewLumberjackFileRotatingLogger(i.Level, i.LogSavePath, i.MaxSize, i.MaxBackups, i.MaxAge)
		sync__, lowLevel__ = zapcore.AddSync(lumberJackLogger), i.Level
		zapcore.Lock(sync__)
		break
	case *ConsoleWriter:
		sync__, lowLevel__ = zapcore.AddSync(os.Stdout), i.Level
		break
	default:
		return fmt.Errorf("unsupported writer: %T", i)
	}

	l.AddHandler(sync__, lowLevel__)
	return nil
}

// GetAndBuildCore
func (l *Logger) GetAndBuildCore(config zapcore.EncoderConfig) zapcore.Core {

	if len(l.Handlers) == 0 {
		panic("why handlers is nil ???")
	}

	var cores []zapcore.Core
	for _, handler := range l.Handlers {
		var tmpEncoder zapcore.Encoder
		switch l.Format_ {
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
	c := l.Config_
	c.StacktraceKey = ""
	l.Logger = l.WrapCore(c)
	return l
}

func (l *Logger) SetStacktrace(level Level) *Logger {
	l.Logger = l.Logger.WithOptions(zap.AddStacktrace(level))
	return l
}

func (l *Logger) NoColor() *Logger {
	c := l.Config_
	c.EncodeLevel = zapcore.CapitalLevelEncoder
	l.Logger = l.WrapCore(c)
	return l
}

// SetTimeFormat sets the log output format.
// default time format is `2006/01/02 - 15:04:05.000`,
func (l *Logger) SetTimeFormat(timeFormat string) *Logger {
	c := l.Config_
	c.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(timeFormat))
	}

	l.Logger = l.WrapCore(c)
	return l
}

// WrapCore wraps or replaces the SugarLogger's underlying zapcore.Core.
func (l *Logger) WrapCore(ec zapcore.EncoderConfig) *zap.Logger {
	return l.Logger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return l.GetAndBuildCore(ec)
	}))
}

// Debug logs messages at DEBUG level
func (l *Logger) Debug(args ...interface{}) {
	l.SLogger.Debug(args...)
}

// Info logs messages at INFO level
func (l *Logger) Info(args ...interface{}) {
	l.SLogger.Info(args...)
}

// Warn logs messages at WARN level
func (l *Logger) Warn(args ...interface{}) {
	l.SLogger.Warn(args...)
}

// Error logs messages at ERROR level
func (l *Logger) Error(args ...interface{}) {
	l.SLogger.Error(args...)
}

// Fatal logs messages at FATAL level
func (l *Logger) Fatal(args ...interface{}) {
	l.SLogger.Fatal(args...)
}

// Panic logs messages at Panic level
func (l *Logger) Panic(args ...interface{}) {
	l.SLogger.Panic(args...)
}

// DPanic logs messages at DPanic level
func (l *Logger) DPanic(args ...interface{}) {
	l.SLogger.DPanic(args...)
}

// DebugW logs messages at DEBUG level
func (l *Logger) DebugW(msg string, keysAndValues ...interface{}) {
	l.SLogger.Debugw(msg, keysAndValues...)
}

// InfoW logs messages at INFO level
func (l *Logger) InfoW(msg string, keysAndValues ...interface{}) {
	l.SLogger.Infow(msg, keysAndValues...)
}

// WarnW logs messages at WARN level
func (l *Logger) WarnW(msg string, keysAndValues ...interface{}) {
	l.SLogger.Warnw(msg, keysAndValues...)
}

// ErrorW logs messages at ERROR level
func (l *Logger) ErrorW(msg string, keysAndValues ...interface{}) {
	l.SLogger.Errorw(msg, keysAndValues...)
}

// FatalW logs messages at FATAL level
func (l *Logger) FatalW(msg string, keysAndValues ...interface{}) {
	l.SLogger.Fatalw(msg, keysAndValues...)
}

// PanicW logs messages at Panic level
func (l *Logger) PanicW(msg string, keysAndValues ...interface{}) {
	l.SLogger.Panicw(msg, keysAndValues...)
}

// DPanicW logs messages at DPanic level
func (l *Logger) DPanicW(msg string, keysAndValues ...interface{}) {
	l.SLogger.DPanicw(msg, keysAndValues...)
}

// DebugF logs messages at DEBUG level
func (l *Logger) DebugF(format string, args ...interface{}) {
	l.SLogger.Debugf(format, args...)
}

// InfoF logs messages at INFO level
func (l *Logger) InfoF(format string, args ...interface{}) {
	l.SLogger.Infof(format, args...)
}

// WarnF logs messages at WARN level
func (l *Logger) WarnF(format string, args ...interface{}) {
	l.SLogger.Warnf(format, args...)
}

// ErrorF logs messages at ERROR level
func (l *Logger) ErrorF(format string, args ...interface{}) {
	l.SLogger.Errorf(format, args...)
}

// Fatalf logs messages at FATAL level
func (l *Logger) FatalF(format string, args ...interface{}) {
	l.SLogger.Fatalf(format, args...)
}

// PanicF logs messages at Panic level
func (l *Logger) PanicF(format string, args ...interface{}) {
	l.SLogger.Panicf(format, args...)
}

// DPanicF logs messages at DPanic level
func (l *Logger) DPanicF(format string, args ...interface{}) {
	l.SLogger.DPanicf(format, args...)
}

// Stack returns the function call stack
func (l *Logger) Stack() string {
	return traceback.TakeStacktrace(1)
}
