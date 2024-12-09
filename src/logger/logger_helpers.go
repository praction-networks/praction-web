package logger

import (
	"runtime/debug"

	"go.uber.org/zap"
)

// Info logs an info-level message with additional key-value pairs
func Info(msg string, args ...interface{}) {
	fields := append(withFields(args), zap.String("loglevel", "info"))
	log.Info(msg, fields...)
}

// Warn logs a warning-level message with additional key-value pairs
func Warn(msg string, args ...interface{}) {
	fields := append(withFields(args), zap.String("loglevel", "warn"))
	log.Warn(msg, fields...)
}

// Error logs an error-level message with additional key-value pairs
func Error(msg string, args ...interface{}) {
	fields := append(withFields(args), zap.String("loglevel", "error"))
	log.Error(msg, fields...)
}

// Debug logs a debug-level message with additional key-value pairs and stack trace
func Debug(msg string, args ...interface{}) {
	fields := append(withFields(args), zap.String("loglevel", "debug"))
	fields = append(fields, zap.String("stacktrace", string(debug.Stack())))
	log.Debug(msg, fields...)
}

// Fatal logs a fatal-level message with additional key-value pairs and exits
func Fatal(msg string, args ...interface{}) {
	fields := append(withFields(args), zap.String("loglevel", "fatal"))
	fields = append(fields, zap.String("stacktrace", string(debug.Stack())))
	log.Fatal(msg, fields...)
}

// Panic logs a panic-level message with additional key-value pairs and triggers panic
func Panic(msg string, args ...interface{}) {
	fields := append(withFields(args), zap.String("loglevel", "panic"))
	fields = append(fields, zap.String("stacktrace", string(debug.Stack())))
	log.Panic(msg, fields...)
}
