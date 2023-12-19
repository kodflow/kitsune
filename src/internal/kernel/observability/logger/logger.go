package logger

import (
	"fmt"
	"log"
	"runtime/debug"

	"github.com/Code-Hex/dd"
	"github.com/kodflow/kitsune/src/internal/kernel/observability/logger/levels"
	"github.com/kodflow/kitsune/src/internal/kernel/observability/logger/writers"
)

// Constant for formatting log path with color.
const PATH = "| \033[38;5;%sm%s\033[39;49m |"

// Logger struct holds the logging configuration and loggers for different levels.
type Logger struct {
	level   levels.TYPE // The logging level of the logger.
	success *log.Logger // Logger for success messages.
	failure *log.Logger // Logger for error messages.
}

// New creates a new Logger instance with the specified writer type and log level.
// It initializes two loggers: one for success and another for failure messages.
//
// Parameters:
// - t: writers.TYPE The writer type (e.g., console, file) for the logger.
// - l: levels.TYPE The logging level for the logger.
//
// Returns:
// - *Logger: A new Logger instance.
func New(t writers.TYPE, l levels.TYPE) *Logger {
	return &Logger{
		level:   l,
		success: log.New(writers.Make(t, writers.SUCCESS), "", log.Ldate|log.Ltime),
		failure: log.New(writers.Make(t, writers.FAILURE), "", log.Ldate|log.Ltime),
	}
}

// Write writes the log message with the specified log level.
// It formats the message and decides which logger to use based on the level.
//
// Parameters:
// - level: levels.TYPE The log level for the message.
// - messages: ...any The messages or data to log.
func (l *Logger) Write(level levels.TYPE, messages ...any) {
	for _, message := range messages {
		var logger *log.Logger = nil
		if level <= levels.WARN {
			logger = l.failure
		} else {
			logger = l.success
		}

		if level == levels.DEBUG {
			logger.Println(fmt.Sprintf(PATH, level.Color(), level.String()), dd.Dump(message, dd.WithIndent(4)))
		} else if level <= l.level {
			logger.Println(fmt.Sprintf(PATH, level.Color(), level.String()), message)
		}
	}
}

// Panic logs the error and stack trace with the PANIC level.
// It logs the error and a stack trace for debugging.
//
// Parameters:
// - err: error The error to log.
//
// Returns:
// - bool: true if the error is not nil, false otherwise.
func (l *Logger) Panic(err error) bool {
	if err != nil {
		l.Write(levels.PANIC, err, string(debug.Stack()))
		return true
	}

	return false
}

// Fatal logs the error and stack trace with the FATAL level.
// It logs critical errors that might require the application to stop.
//
// Parameters:
// - err: error The error to log.
//
// Returns:
// - bool: true if the error is not nil, false otherwise.
func (l *Logger) Fatal(err error) bool {
	if err != nil {
		l.Write(levels.FATAL, err, string(debug.Stack()))
		return true
	}

	return false
}

// Error logs the error with the ERROR level.
// It is used for logging general errors.
//
// Parameters:
// - err: error The error to log.
//
// Returns:
// - bool: true if the error is not nil, false otherwise.
func (l *Logger) Error(err error) bool {
	if err != nil {
		l.Write(levels.ERROR, err)
		return true
	}

	return false
}

// Success logs the success message with the SUCCESS level.
// It is used for logging successful operations.
//
// Parameters:
// - v: ...any The success messages or data to log.
func (l *Logger) Success(v ...any) {
	l.Write(levels.SUCCESS, v...)
}

// Message logs the message with the MESSAGE level.
// It is used for general-purpose logging.
//
// Parameters:
// - v: ...any The messages or data to log.
func (l *Logger) Message(v ...any) {
	l.Write(levels.MESSAGE, v...)
}

// Warn logs the warning message with the WARN level.
// It is used for logging potential issues or warnings.
//
// Parameters:
// - v: ...any The warning messages or data to log.
func (l *Logger) Warn(v ...any) {
	l.Write(levels.WARN, v...)
}

// Info logs the info message with the INFO level.
// It is used for logging informational messages.
//
// Parameters:
// - v: ...any The informational messages or data to log.
func (l *Logger) Info(v ...any) {
	l.Write(levels.INFO, v...)
}

// Debug logs the debug message with the DEBUG level.
// It provides detailed debug information for troubleshooting.
//
// Parameters:
// - v: ...any The debug messages or data to log.
func (l *Logger) Debug(v ...any) {
	l.Write(levels.DEBUG, v...)
}

// Trace logs the stack trace with the TRACE level.
// It is used for logging detailed execution traces for in-depth debugging.
func (l *Logger) Trace() {
	l.Write(levels.TRACE, string(debug.Stack()))
}
