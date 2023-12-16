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

func TestDeleteDirectory(t *testing.T) {
	err := fs.CreateDirectory(VALID_DIR_PATH)
	assert.NoError(t, err)

	err = fs.DeleteDirectory(VALID_DIR_PATH)
	assert.NoError(t, err)

	exists := fs.ExistsDirectory(VALID_DIR_PATH)
	assert.False(t, exists)
}
