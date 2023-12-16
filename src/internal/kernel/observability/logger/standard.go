package logger

import (
	"github.com/kodmain/kitsune/src/internal/kernel/observability/logger/levels"
	"github.com/kodmain/kitsune/src/internal/kernel/observability/logger/writers"
)

var instance *Logger = nil

// standard returns an instance of the standard logger.
// If no instance exists, it is created with the default parameters.
func standard() *Logger {
	if instance == nil {
		instance = New(writers.DEFAULT, levels.DEFAULT)
	}
	return instance
}

// SetLevel sets the log level of the standard logger.
func SetLevel(l levels.TYPE) {
	standard().level = l
}

// Panic logs a message with Panic level.
// Returns true if the message was successfully logged, otherwise false.
func Panic(err error) bool {
	return standard().Panic(err)
}

// Fatal logs a message with Fatal level.
// Returns true if the message was successfully logged, otherwise false.
func Fatal(err error) bool {
	return standard().Fatal(err)
}

// Error logs a message with Error level.
// Returns true if the message was successfully logged, otherwise false.
func Error(err error) bool {
	return standard().Error(err)
}

// Success logs a message with Success level.
func Success(v ...any) {
	standard().Success(v...)
}

// Message logs a message with Message level.
func Message(v ...any) {
	standard().Message(v...)
}

// Warn logs a message with Warn level.
func Warn(v ...any) {
	standard().Warn(v...)
}

// Info logs a message with Info level.
func Info(v ...any) {
	standard().Info(v...)
}

// Debug logs a message with Debug level.
func Debug(v ...any) {
	standard().Debug(v...)
}

// Trace logs a message with Trace level.
func Trace() {
	standard().Trace()
}
