package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type RotateWriter struct {
	Level       Level
	LogSavePath string // path for saving logs
	LogFileExt  string // Log file suffix
	MaxSize     int    // size of backup
	MaxBackups  int    // Maximum backup number
	MaxAge      int    // Maximum backup days
	Compress    bool   // Whether to compress expiration logs
}

type ConsoleWriter struct {
	Level Level
}

type Handler struct {
	Sync       zapcore.WriteSyncer
	EnableFunc zap.LevelEnablerFunc
}

// NewRotateWriter returns rotate logs configuration
func NewRotateWriter(level Level, file string, maxSize, maxBackups, maxAges int) *RotateWriter {
	r := new(RotateWriter)
	r.Level = level
	r.LogSavePath = file
	r.MaxSize = maxSize
	r.MaxAge = maxAges
	r.MaxBackups = maxBackups
	return r
}

// NewConsoleWriter returns console los
func NewConsoleWriter(level Level) *ConsoleWriter {
	c := new(ConsoleWriter)
	c.Level = level
	return c
}
