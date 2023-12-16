package logger_test

import (
	"errors"
	"testing"

	"github.com/kodmain/kitsune/src/internal/kernel/observability/logger"
	"github.com/kodmain/kitsune/src/internal/kernel/observability/logger/levels"
	"github.com/kodmain/kitsune/src/internal/kernel/observability/logger/writers"
	"github.com/kodmain/kitsune/src/internal/kernel/storages/fs"
	"github.com/stretchr/testify/assert"
)

func newTestLogger() *logger.Logger {
	return logger.New(writers.FILE, levels.TRACE)
}

func TestLoggerErrorMethods(t *testing.T) {
	logger := newTestLogger()

	err := errors.New("test error")
	assert.True(t, logger.Error(err), "Error should return true if an error is provided")
	assert.True(t, logger.Fatal(err), "Fatal should return true if an error is provided")
	assert.True(t, logger.Panic(err), "Panic should return true if an error is provided")
}

func TestLoggerLevels(t *testing.T) {
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
			ok, err := fs.Contains(writers.FILE_STDOUT, test.contains)
			assert.Nil(t, err, "Error should be nil")
			assert.True(t, ok, "File should contain the message")
		} else {
			ok, err := fs.Contains(writers.FILE_STDERR, test.contains)
			assert.Nil(t, err, "Error should be nil")
			assert.True(t, ok, "File should contain the message", writers.FILE_STDERR)
		}
	}
}
