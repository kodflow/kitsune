package logger

import (
	"github.com/kodmain/kitsune/src/internal/kernel/observability/logger/levels"
	"github.com/kodmain/kitsune/src/internal/kernel/observability/logger/writers"
)

var instance *Logger = nil

// standard returns an instance of the standard logger.
// If no instance exists, it is created with default parameters using New function.
// This function implements the singleton pattern to ensure only one instance of Logger exists.
//
// Returns:
// - *Logger: The singleton instance of the Logger.
func standard() *Logger {
	if instance == nil {
		instance = New(writers.DEFAULT, levels.DEFAULT)
	}
	return instance
}

// SetLevel sets the log level of the standard logger.
// This function allows dynamic adjustment of the logging level.
//
// Parameters:
// - l: levels.TYPE The log level to set for the standard logger.
func SetLevel(l levels.TYPE) {
	standard().level = l
}

// Panic logs a message with Panic level.
// It logs a message and then panics.
//
// Parameters:
// - err: error The error to log.
//
// Returns:
// - bool: true if the message was successfully logged, otherwise false.
func Panic(err error) bool {
	return standard().Panic(err)
}

// Fatal logs a message with Fatal level.
// It logs a critical error message and typically exits the program.
//
// Parameters:
// - err: error The error to log.
//
// Returns:
// - bool: true if the message was successfully logged, otherwise false.
func Fatal(err error) bool {
	return standard().Fatal(err)
}

// Error logs a message with Error level.
// It logs an error message, used for non-critical failures.
//
// Parameters:
// - err: error The error to log.
//
// Returns:
// - bool: true if the message was successfully logged, otherwise false.
func Error(err error) bool {
	return standard().Error(err)
}

// Success logs a message with Success level.
// It logs a success or completion message.
//
// Parameters:
// - v: ...any The message or variables to log.
func Success(v ...any) {
	standard().Success(v...)
}

// Message logs a message with Message level.
// It logs a general, informational message.
//
// Parameters:
// - v: ...any The message or variables to log.
func Message(v ...any) {
	standard().Message(v...)
}

// Warn logs a message with Warn level.
// It logs a warning message, indicating potential issues.
//
// Parameters:
// - v: ...any The message or variables to log.
func Warn(v ...any) {
	standard().Warn(v...)
}

// Info logs a message with Info level.
// It logs informational messages, useful for tracking the flow of the application.
//
// Parameters:
// - v: ...any The message or variables to log.
func Info(v ...any) {
	standard().Info(v...)
}

// Debug logs a message with Debug level.
// It logs detailed debugging information.
//
// Parameters:
// - v: ...any The message or variables to log.
func Debug(v ...any) {
	standard().Debug(v...)
}

// Trace logs a message with Trace level.
// It logs the most detailed information, often for tracing code execution paths.
func Trace() {
	standard().Trace()
}
