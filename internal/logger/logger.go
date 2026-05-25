// File: internal/logger/logger.go
// Version: v0.1
// Purpose: Initialize and provide a structured logger (zap) for the entire application.
// Security: Never log Discord token, MongoDB URI, passwords, or any secret values.
// Notes: Use logger.With(fields...) to attach context (userId, guildId, command) to logs.

package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

// Init initializes the global logger with the given log level.
// Call this once at application startup before any logging.
func Init(level string) error {
	zapLevel, err := parseLevel(level)
	if err != nil {
		zapLevel = zapcore.InfoLevel
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zapLevel)
	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	log, err := cfg.Build(zap.AddCallerSkip(0))
	if err != nil {
		return err
	}

	globalLogger = log
	return nil
}

// L returns the global logger. Panics if Init has not been called.
func L() *zap.Logger {
	if globalLogger == nil {
		panic("logger: Init() must be called before using L()")
	}
	return globalLogger
}

// S returns the global SugaredLogger for printf-style logging.
func S() *zap.SugaredLogger {
	return L().Sugar()
}

// Sync flushes any buffered log entries. Call on application shutdown.
func Sync() {
	if globalLogger != nil {
		_ = globalLogger.Sync()
	}
}

func parseLevel(level string) (zapcore.Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn", "warning":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InfoLevel, nil
	}
}
