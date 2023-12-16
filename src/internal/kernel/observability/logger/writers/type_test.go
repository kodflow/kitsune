package writers_test

import (
	"io"
	"testing"

	"github.com/kodmain/kitsune/src/internal/kernel/observability/logger/writers"
	"github.com/stretchr/testify/assert"
)

func TestMake(t *testing.T) {
	// Test case 1: CONSOLE and SOF true
	w := writers.Make(writers.CONSOLE, true)
	assert.NotNil(t, w)
	assert.Implements(t, (*io.Writer)(nil), w)

	// Test case 2: CONSOLE and SOF false
	w = writers.Make(writers.CONSOLE, false)
	assert.NotNil(t, w)
	assert.Implements(t, (*io.Writer)(nil), w)

	// Test case 3: FILE and SOF true
	w = writers.Make(writers.FILE, true)
	assert.NotNil(t, w)
	assert.Implements(t, (*io.Writer)(nil), w)

	// Test case 4: FILE and SOF false
	w = writers.Make(writers.FILE, false)
	assert.NotNil(t, w)
	assert.Implements(t, (*io.Writer)(nil), w)
}
