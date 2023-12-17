package logger_test

import (
	"errors"
	"testing"

	"github.com/kodflow/kitsune/src/internal/kernel/observability/logger"
	"github.com/stretchr/testify/assert"
)

// TestLoggerMethods tests various methods of the logger package.
// It checks the behavior of different logging methods, particularly focusing on error handling and ensuring no panics.
func TestLoggerMethods(t *testing.T) {
	// Error simulation
	testError := errors.New("test error")

	// Test for Panic, Fatal, Error methods
	// Ensures that Panic, Fatal, and Error methods return true when an error is provided.
	assert.True(t, logger.Panic(testError), "Panic should return true if an error is provided")
	assert.True(t, logger.Fatal(testError), "Fatal should return true if an error is provided")
	assert.True(t, logger.Error(testError), "Error should return true if an error is provided")

	assert.False(t, logger.Panic(nil), "Panic should return false if an error is provided")
	assert.False(t, logger.Fatal(nil), "Fatal should return false if an error is provided")
	assert.False(t, logger.Error(nil), "Error should return false if an error is provided")

	// Test for other methods, ensuring they do not produce an error
	// Verifies that methods such as Success, Message, Warn, Info, Debug, and Trace do not cause panics.
	assert.NotPanics(t, func() { logger.Success("test success") }, "Success should not panic")
	assert.NotPanics(t, func() { logger.Message("test message") }, "Message should not panic")
	assert.NotPanics(t, func() { logger.Warn("test warn") }, "Warn should not panic")
	assert.NotPanics(t, func() { logger.Info("test info") }, "Info should not panic")
	assert.NotPanics(t, func() { logger.Debug("test debug") }, "Debug should not panic")
	assert.NotPanics(t, func() { logger.Trace() }, "Trace should not panic")
}
