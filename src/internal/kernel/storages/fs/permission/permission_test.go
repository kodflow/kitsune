package permission

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	Kitsune               = "/tmp/kitsune"
	VALID_FILE_PATH       = Kitsune + "/file"
	INVALID_FILE_PATH     = Kitsune + "/invalid"
	NONEXISTENT_FILE_PATH = Kitsune + "/nonexistent"
)

// TestValidate tests the Validate function.
func TestValidate(t *testing.T) {
	// Clean up the temporary directory after the test.
	defer os.RemoveAll(Kitsune)
	os.MkdirAll(VALID_FILE_PATH, 0755)

	// Validate that the valid file has the specified permissions.
	valid := Validate(VALID_FILE_PATH, 0644)
	assert.True(t, valid)

	// Validate that the invalid file does not have the specified permissions.
	valid = Validate(INVALID_FILE_PATH, 0755)
	assert.False(t, valid)

	// Validate that a non-existent file does not have the specified permissions.
	valid = Validate(NONEXISTENT_FILE_PATH, 0644)
	assert.False(t, valid)
}
