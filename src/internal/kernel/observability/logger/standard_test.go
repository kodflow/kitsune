package logger_test

import (
	"errors"
	"testing"

	"github.com/kodmain/kitsune/src/internal/kernel/observability/logger"
	"github.com/stretchr/testify/assert"
)

func TestLoggerMethods(t *testing.T) {
	// Error simulation
	testError := errors.New("test error")

	// Test for Panic, Fatal, Error methods
	assert.True(t, logger.Panic(testError), "Panic should return true if an error is provided")
	assert.True(t, logger.Fatal(testError), "Fatal should return true if an error is provided")
	assert.True(t, logger.Error(testError), "Error should return true if an error is provided")

	// Test for other methods, ensuring they do not produce an error
	assert.NotPanics(t, func() { logger.Success("test success") }, "Success should not panic")
	assert.NotPanics(t, func() { logger.Message("test message") }, "Message should not panic")
	assert.NotPanics(t, func() { logger.Warn("test warn") }, "Warn should not panic")
	assert.NotPanics(t, func() { logger.Info("test info") }, "Info should not panic")
	assert.NotPanics(t, func() { logger.Debug("test debug") }, "Debug should not panic")
	assert.NotPanics(t, func() { logger.Trace() }, "Trace should not panic")
}
