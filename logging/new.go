package logging

import (
	"fmt"
)

var (
	defaultLogger Log
	defaultFlush  func() error
)

// Instantiate `defaultLogger` and `defaultFlush` while waiting for package initialization
func init() {
	defaultLogger, defaultFlush, _ = NewConsoleStreamingLogger("go-pkg", DebugLevel, WarnLevel, DebugLevel)
}

// GLogger returns the global logger in `go-pkg`
func GLogger() (logger_ Log, flush func() error) {
	return defaultLogger, defaultFlush
}

// ReplaceGlobalLogger can replace the default logger with custom logger
func ReplaceGlobalLogger(newLogger Log, flush func() error) {
	defaultLogger = newLogger
	flush = defaultFlush
}

// NewFileRotatingLogger returns logging instance by rotating file.
func NewFileRotatingLogger(
	name, filePath string,
	logLevel, stackLevel, handlerLevel Level,
	maxSize, maxBackups, maxAge int) (logger_ Log, flush func() error, err error) {
	writer := NewRotateWriter(handlerLevel, filePath, maxSize, maxBackups, maxAge)
	return NewLogger(logLevel, stackLevel, name, "", false, EncodeConsole, writer)
}

// NewConsoleStreamingLogger returns logging instance which outputs message to stdout.
func NewConsoleStreamingLogger(
	name string,
	outputLevel, stackLevel, handlerLevel Level) (logger_ Log, flush func() error, err error) {
	writer := NewConsoleWriter(handlerLevel)
	return NewLogger(outputLevel, stackLevel, name, "", false, EncodeConsole, writer)
}

// GetLoggerFromConf returns logging instance by configuration file
func GetLoggerFromConf(name, conf string) (logger_ Log, flush func() error, err error) {
	__loggerContainersOnce.Do(func() {
		if err = initLoggersContainers(conf); err != nil {
			return
		}
	})
	if __k, ok := __containers[name]; ok && __k._logger != nil && __k._flushFunc != nil {
		return __k._logger, __k._flushFunc, nil
	}
	return nil, nil, fmt.Errorf("failed to get logger named `%s` from `%s`", name, conf)
}
