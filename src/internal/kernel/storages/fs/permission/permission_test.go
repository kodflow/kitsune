package permission_test

import (
	"os"
	"testing"

	"github.com/kodflow/kitsune/src/internal/kernel/storages/fs"
	"github.com/kodflow/kitsune/src/internal/kernel/storages/fs/permission"
	"github.com/stretchr/testify/assert"
)

const (
	VALID_FILE_PATH       = "/tmp/kitsune/file"
	INVALID_FILE_PATH     = "/tmp/kitsune/invalid"
	NONEXISTENT_FILE_PATH = "/tmp/kitsune/nonexistent"
)

// TestValidate tests the Validate function.
func TestValidate(t *testing.T) {
	// Clean up the temporary directory after the test.
	defer os.RemoveAll("/tmp/kitsune")

	// Create a valid file.
	fv, err := fs.CreateFile(VALID_FILE_PATH)
	assert.NoError(t, err)
	assert.NotNil(t, fv)

	// Create an invalid file.
	fi, err := fs.CreateFile(INVALID_FILE_PATH)
	assert.NoError(t, err)
	assert.NotNil(t, fi)

	// Validate that the valid file has the specified permissions.
	valid := permission.Validate(VALID_FILE_PATH, 0644)
	assert.True(t, valid)

	// Validate that the invalid file does not have the specified permissions.
	valid = permission.Validate(INVALID_FILE_PATH, 0755)
	assert.False(t, valid)

	// Validate that a non-existent file does not have the specified permissions.
	valid = permission.Validate(NONEXISTENT_FILE_PATH, 0644)
	assert.False(t, valid)
}
