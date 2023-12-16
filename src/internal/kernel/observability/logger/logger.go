package logger

import (
	"fmt"
	"log"
	"runtime/debug"

	"github.com/Code-Hex/dd"
	"github.com/kodmain/kitsune/src/internal/kernel/observability/logger/levels"
	"github.com/kodmain/kitsune/src/internal/kernel/observability/logger/writers"
)

const PATH = "| \033[38;5;%sm%s\033[39;49m |"

type Logger struct {
	level   levels.TYPE
	success *log.Logger
	failure *log.Logger
}

// New creates a new Logger instance with the specified writer type and log level.
func New(t writers.TYPE, l levels.TYPE) *Logger {
	return &Logger{
		level:   l,
		success: log.New(writers.Make(t, writers.SUCCESS), "", log.Ldate|log.Ltime),
		failure: log.New(writers.Make(t, writers.FAILURE), "", log.Ldate|log.Ltime),
	}
}

// Write writes the log message with the specified log level.
func (l *Logger) Write(level levels.TYPE, messages ...interface{}) {
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
// It returns true if the error is not nil, false otherwise.
func (l *Logger) Panic(err error) bool {
	if err != nil {
		l.Write(levels.PANIC, err, string(debug.Stack()))
		return true
	}

	return false
}

// Fatal logs the error and stack trace with the FATAL level.
// It returns true if the error is not nil, false otherwise.
func (l *Logger) Fatal(err error) bool {
	if err != nil {
		l.Write(levels.FATAL, err, string(debug.Stack()))
		return true
	}

	return false
}

// Error logs the error with the ERROR level.
// It returns true if the error is not nil, false otherwise.
func (l *Logger) Error(err error) bool {
	if err != nil {
		l.Write(levels.ERROR, err)
		return true
	}

	return false
}

// Success logs the success message with the SUCCESS level.
func (l *Logger) Success(v ...interface{}) {
	l.Write(levels.SUCCESS, v...)
}

// Message logs the message with the MESSAGE level.
func (l *Logger) Message(v ...interface{}) {
	l.Write(levels.MESSAGE, v...)
}

// Warn logs the warning message with the WARN level.
func (l *Logger) Warn(v ...interface{}) {
	l.Write(levels.WARN, v...)
}

// Info logs the info message with the INFO level.
func (l *Logger) Info(v ...interface{}) {
	l.Write(levels.INFO, v...)
}

// Debug logs the debug message with the DEBUG level.
func (l *Logger) Debug(v ...interface{}) {
	l.Write(levels.DEBUG, v...)
}

// Trace logs the stack trace with the TRACE level.
func (l *Logger) Trace() {
	l.Write(levels.TRACE, string(debug.Stack()))
}
