package logging

var (
	// Instantiate `defaultLogger` while waiting for package initialization
	defaultLogger, _ = NewConsoleStreamingLogger("go-pkg", DebugLevel, WarnLevel, DebugLevel)
)

// GLogger returns the global logger in `go-pkg`
func GLogger() (logger_ *Logger) {
	return defaultLogger
}

// ReplaceGlobalLogger can replace the default logger with custom logger
func ReplaceGlobalLogger(newLogger *Logger) {
	defaultLogger = newLogger
}

// NewFileRotatingLogger returns logging instance by rotating file.
func NewFileRotatingLogger(
	name, filePath string,
	logLevel, stackLevel, handlerLevel Level,
	maxSize, maxBackups, maxAge int) (logger_ *Logger, err error) {
	writer := NewRotateWriter(handlerLevel, filePath, maxSize, maxBackups, maxAge)
	return NewLogger(logLevel, stackLevel, name, "", false, EncodeConsole, writer)
}

// NewConsoleStreamingLogger returns logging instance which outputs message to stdout.
func NewConsoleStreamingLogger(
	name string,
	outputLevel, stackLevel, handlerLevel Level) (logger_ *Logger, err error) {
	writer := NewConsoleWriter(handlerLevel)
	return NewLogger(outputLevel, stackLevel, name, "", false, EncodeConsole, writer)
}
