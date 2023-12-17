package fs_test

import (
	"os"
	"testing"

	"github.com/kodflow/kitsune/src/internal/kernel/storages/fs"
	"github.com/stretchr/testify/assert"
)

const (
	VALID_FILE_PATH       = "/tmp/kitsune/file"
	INVALID_FILE_PATH     = "/tmp/kitsune/file-invalid"
	NONEXISTENT_FILE_PATH = "/tmp/kitsune/file-nonexistent"
)

// TestCreateFile tests the file creation functionality in the fs package.
// It ensures that a file is successfully created without errors.
func TestCreateFile(t *testing.T) {
	defer os.RemoveAll("/tmp/kitsune")

	fs, err := fs.CreateFile(VALID_FILE_PATH)
	assert.NoError(t, err)
	assert.NotNil(t, fs)
}

// TestPermissions tests the setting of file permissions.
// It verifies that the file permissions are correctly set and can be retrieved.
func TestPermissions(t *testing.T) {
	defer os.RemoveAll("/tmp/kitsune")

	_, err := fs.CreateFile(VALID_FILE_PATH)
	assert.NoError(t, err)

	err = fs.Permissions(VALID_FILE_PATH, &fs.Options{Perms: 0600})
	assert.NoError(t, err)

	info, err := os.Stat(VALID_FILE_PATH)
	assert.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
}

// TestOpenFile tests the file opening functionality.
// It checks if a file can be opened and then verifies its existence.
func TestOpenFile(t *testing.T) {
	defer os.Remove(VALID_FILE_PATH)

	file, err := fs.OpenFile(VALID_FILE_PATH, &fs.Options{Perms: 0644})
	assert.NoError(t, err)
	assert.NotNil(t, file)
	file.Close()

	exists := fs.ExistsFile(VALID_FILE_PATH)
	assert.True(t, exists)
}

// TestDeleteFile tests the file deletion functionality.
// It ensures that a file is deleted and confirms its non-existence.
func TestDeleteFile(t *testing.T) {
	_, err := fs.CreateFile(VALID_FILE_PATH, &fs.Options{Perms: 0644})
	assert.NoError(t, err)

	err = fs.DeleteFile(VALID_FILE_PATH)
	assert.NoError(t, err)

	exists := fs.ExistsFile(VALID_FILE_PATH)
	assert.False(t, exists)
}

// TestSHA1File tests the SHA1 hashing functionality for files.
// It verifies that the hash is generated correctly and is not empty.
func TestSHA1File(t *testing.T) {
	filePath := "test_sha1.txt"
	content := "Hello, World!"
	defer os.Remove(filePath)

	err := fs.WriteFile(filePath, content)
	assert.NoError(t, err)

	hash, err := fs.SHA1File(filePath)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
}

// TestExistsFile checks if the file existence checking functionality works.
// It verifies both the existence and non-existence cases.
func TestExistsFile(t *testing.T) {
	filePath := "test_exists.txt"

	_, err := fs.CreateFile(filePath, &fs.Options{Perms: 0644})
	defer os.Remove(filePath)
	assert.NoError(t, err)

	exists := fs.ExistsFile(filePath)
	assert.True(t, exists)

	exists = fs.ExistsFile("nonexistent.txt")
	assert.False(t, exists)
}

// TestStatFile tests the file status retrieval functionality.
// It confirms that the file status information is correctly obtained.
func TestStatFile(t *testing.T) {
	filePath := "test_stat.txt"
	defer os.Remove(filePath)

	_, err := fs.CreateFile(filePath, &fs.Options{Perms: 0644})
	assert.NoError(t, err)

	info, err := fs.StatFile(filePath)
	assert.NoError(t, err)
	assert.NotNil(t, info)
}

// TestReadFile tests reading content from a file.
// It ensures that the content read is the same as what was written.
func TestReadFile(t *testing.T) {
	filePath := "test_read.txt"
	content := "Test content"
	defer os.Remove(filePath)

	err := fs.WriteFile(filePath, content)
	assert.NoError(t, err)

	readContent, err := fs.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, content, string(readContent))
}

// TestWriteFile tests the file writing functionality.
// It verifies that the content written to a file can be read back accurately.
func TestWriteFile(t *testing.T) {
	filePath := "test_write.txt"
	content := "New content"
	defer os.Remove(filePath)

	err := fs.WriteFile(filePath, content)
	assert.NoError(t, err)

	readContent, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, content, string(readContent))
}

// TestContains tests the substring checking functionality in a file.
// It checks both the presence and absence of specified substrings.
func TestContains(t *testing.T) {
	filePath := "test_contains.txt"
	content := "Hello, World!"
	substring := "World"
	defer os.Remove(filePath)

	err := fs.WriteFile(filePath, content)
	assert.NoError(t, err)

	contains, err := fs.Contains(filePath, substring)
	assert.NoError(t, err)
	assert.True(t, contains)

	contains, err = fs.Contains(filePath, "Nonexistent")
	assert.NoError(t, err)
	assert.False(t, contains)
}
