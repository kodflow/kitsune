package fs_test

import (
	"os"
	"testing"

	"github.com/kodmain/kitsune/src/internal/kernel/storages/fs"
	"github.com/stretchr/testify/assert"
)

const (
	VALID_DIR_PATH       = "/tmp/kitsune/dir"
	INVALID_DIR_PATH     = "/tmp/kitsune/dir-invalid"
	NONEXISTENT_DIR_PATH = "/tmp/kitsune/dir-nonexistent"
)

// TestCreateDirectory tests the directory creation functionality in the fs package.
// It ensures that a directory is successfully created without errors.
// This test also verifies that creating an already existing directory does not cause an error.
func TestCreateDirectory(t *testing.T) {
	defer os.RemoveAll("/tmp/kitsune")

	err := fs.CreateDirectory(VALID_DIR_PATH)
	assert.NoError(t, err)

	err = fs.CreateDirectory(VALID_DIR_PATH)
	assert.NoError(t, err)

	_, err = os.Stat(VALID_DIR_PATH)
	assert.NoError(t, err)

	err = os.RemoveAll(VALID_DIR_PATH)
	assert.NoError(t, err)
}

// TestExistsDirectory tests the directory existence checking functionality.
// It verifies that the function correctly identifies the existence and non-existence of a directory.
func TestExistsDirectory(t *testing.T) {
	err := fs.CreateDirectory(VALID_DIR_PATH)
	assert.NoError(t, err)

	exists := fs.ExistsDirectory(VALID_DIR_PATH)
	assert.True(t, exists)

	err = os.RemoveAll(VALID_DIR_PATH)
	assert.NoError(t, err)

	exists = fs.ExistsDirectory(VALID_DIR_PATH)
	assert.False(t, exists)
}

// TestDeleteDirectory tests the directory deletion functionality.
// It ensures that a directory is successfully deleted and confirms its non-existence post-deletion.
func TestDeleteDirectory(t *testing.T) {
	err := fs.CreateDirectory(VALID_DIR_PATH)
	assert.NoError(t, err)

	err = fs.DeleteDirectory(VALID_DIR_PATH)
	assert.NoError(t, err)

	exists := fs.ExistsDirectory(VALID_DIR_PATH)
	assert.False(t, exists)
}
