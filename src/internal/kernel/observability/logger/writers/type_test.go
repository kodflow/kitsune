package writers_test

import (
	"io"
	"testing"

	"github.com/kodflow/kitsune/src/internal/kernel/observability/logger/writers"
	"github.com/stretchr/testify/assert"
)

// TestMake tests the Make function in the writers package.
// This function verifies that the Make function correctly creates writer instances
// for different configurations (console and file) and buffer settings (true or false).
// It also checks if the created writers implement the io.Writer interface.
func TestMake(t *testing.T) {
	// Test creation of a console writer with buffering enabled
	w := writers.Make(writers.CONSOLE, true)
	assert.NotNil(t, w, "Console writer (buffered) should not be nil")
	assert.Implements(t, (*io.Writer)(nil), w, "Console writer (buffered) should implement io.Writer")

	// Test creation of a console writer with buffering disabled
	w = writers.Make(writers.CONSOLE, false)
	assert.NotNil(t, w, "Console writer (unbuffered) should not be nil")
	assert.Implements(t, (*io.Writer)(nil), w, "Console writer (unbuffered) should implement io.Writer")

	// Test creation of a file writer with buffering enabled
	w = writers.Make(writers.FILE, true)
	assert.NotNil(t, w, "File writer (buffered) should not be nil")
	assert.Implements(t, (*io.Writer)(nil), w, "File writer (buffered) should implement io.Writer")

	// Test creation of a file writer with buffering disabled
	w = writers.Make(writers.FILE, false)
	assert.NotNil(t, w, "File writer (unbuffered) should not be nil")
	assert.Implements(t, (*io.Writer)(nil), w, "File writer (unbuffered) should implement io.Writer")
}
