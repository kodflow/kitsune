package logger_test

import (
	"errors"
	"os"
	"testing"

	"github.com/kodflow/kitsune/src/config"
	"github.com/kodflow/kitsune/src/internal/kernel/observability/logger"
	"github.com/kodflow/kitsune/src/internal/kernel/observability/logger/levels"
	"github.com/kodflow/kitsune/src/internal/kernel/observability/logger/writers"
	"github.com/kodflow/kitsune/src/internal/kernel/storages/fs"
	"github.com/stretchr/testify/assert"
)

// newTestLogger creates and returns a new logger instance for testing.
// This logger is configured with a file writer and trace level logging.
func newTestLogger() *logger.Logger {
	return logger.New(writers.FILE, levels.TRACE)
}

// TestLoggerErrorMethods tests the error handling methods of the logger.
// It verifies that Error, Fatal, and Panic methods return true when provided with an error.
func TestLoggerErrorMethods(t *testing.T) {
	logger := newTestLogger()

	err := errors.New("test error")
	assert.True(t, logger.Error(err), "Error should return true if an error is provided")
	assert.True(t, logger.Fatal(err), "Fatal should return true if an error is provided")
	assert.True(t, logger.Panic(err), "Panic should return true if an error is provided")
}

// TestLoggerLevels tests the logger's ability to handle different logging levels.
// It writes messages at various levels and checks if they are correctly written to the respective output files.
func TestLoggerLevels(t *testing.T) {
	defer os.RemoveAll(config.PATH_LOGS)
	logger := newTestLogger()

	tests := []struct {
		level    levels.TYPE
		message  string
		contains string
	}{
		{levels.DEBUG, "debug message", "debug message"},
		{levels.INFO, "info message", "info message"},
		{levels.WARN, "warn message", "warn message"},
		{levels.ERROR, "error message", "error message"},
		{levels.FATAL, "fatal message", "fatal message"},
		{levels.PANIC, "panic message", "panic message"},
		{levels.SUCCESS, "success message", "success message"},
		{levels.MESSAGE, "simple message", "simple message"},
		{levels.TRACE, "trace message", "trace message"},
	}

	for _, test := range tests {
		logger.Write(test.level, test.message)
		if test.level.Int() > levels.WARN.Int() {
			exists, err := fs.ExistsFile(writers.FILE_STDOUT)
			assert.NoError(t, err, "Error should be nil")
			assert.True(t, exists, "Stdout file should exist")
			ok, err := fs.Contains(writers.FILE_STDOUT, test.contains)
			assert.Nil(t, err, "Error should be nil")
			assert.True(t, ok, "File should contain the message")
		} else {
			exists, err := fs.ExistsFile(writers.FILE_STDERR)
			assert.NoError(t, err, "Error should be nil")
			assert.True(t, exists, "Stderr file should exist")
			ok, err := fs.Contains(writers.FILE_STDERR, test.contains)
			assert.Nil(t, err, "Error should be nil")
			assert.True(t, ok, "File should contain the message", writers.FILE_STDERR)
		}
	}
}
