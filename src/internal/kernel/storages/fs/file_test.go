package fs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	VALID_FILE_PATH       = "/tmp/kitsune/file"
	INVALID_FILE_PATH     = "/path/to/file-invalid"
	NONEXISTENT_FILE_PATH = "/tmp/kitsune/file-nonexistent"
)

// TestCreateFile tests the file creation functionality in the fs package.
// It ensures that a file is successfully created without errors.
func TestCreateFile(t *testing.T) {
	defer os.RemoveAll("/tmp/kitsune")

	fv, err := CreateFile(VALID_FILE_PATH)
	assert.NoError(t, err)
	assert.NotNil(t, fv)

	fi, err := CreateFile(INVALID_FILE_PATH)
	assert.Error(t, err)
	assert.Nil(t, fi)
}

// TestPermissions tests the setting of file permissions.
// It verifies that the file permissions are correctly set and can be retrieved.
func TestPermissions(t *testing.T) {
	defer os.RemoveAll("/tmp/kitsune")

	_, err := CreateFile(VALID_FILE_PATH)
	assert.NoError(t, err)

	err = Permissions(VALID_FILE_PATH, &Options{Perms: 0600})
	assert.NoError(t, err)

	info, err := os.Stat(VALID_FILE_PATH)
	assert.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
}

// TestOpenFile tests the file opening functionality.
// It checks if a file can be opened and then verifies its existence.
func TestOpenFile(t *testing.T) {
	defer os.Remove(VALID_FILE_PATH)

	file, err := OpenFile(INVALID_FILE_PATH, &Options{Perms: 0644})
	assert.Error(t, err)
	assert.Nil(t, file)
	file.Close()

	file, err = CreateFile(VALID_FILE_PATH, &Options{Perms: 0644})
	assert.NoError(t, err)
	assert.NotNil(t, file)
	file.Close()

	exists := ExistsFile(VALID_FILE_PATH)
	assert.True(t, exists)
}

// TestDeleteFile tests the file deletion functionality.
// It ensures that a file is deleted and confirms its non-existence.
func TestDeleteFile(t *testing.T) {
	_, err := CreateFile(VALID_FILE_PATH, &Options{Perms: 0644})
	assert.NoError(t, err)

	err = DeleteFile(VALID_FILE_PATH)
	assert.NoError(t, err)

	exists := ExistsFile(VALID_FILE_PATH)
	assert.False(t, exists)
}

// TestSHA1File tests the SHA1 hashing functionality for files.
// It verifies that the hash is generated correctly and is not empty.
func TestSHA1File(t *testing.T) {
	defer os.Remove(VALID_FILE_PATH)

	err := WriteFile(VALID_FILE_PATH, "Hello, World!")
	assert.NoError(t, err)

	hash, err := SHA1File(VALID_FILE_PATH)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	hash, err = SHA1File(INVALID_FILE_PATH)
	assert.Error(t, err)
	assert.Empty(t, hash)
	assert.Contains(t, err.Error(), "no such file or directory")

	err = WriteFile(VALID_FILE_PATH, "Hello, World!")
	assert.NoError(t, err)
}

// TestExistsFile checks if the file existence checking functionality works.
// It verifies both the existence and non-existence cases.
func TestExistsFile(t *testing.T) {
	filePath := "test_exists.txt"

	_, err := CreateFile(filePath, &Options{Perms: 0644})
	defer os.Remove(filePath)
	assert.NoError(t, err)

	exists := ExistsFile(filePath)
	assert.True(t, exists)

	exists = ExistsFile("nonexistent.txt")
	assert.False(t, exists)
}

// TestStatFile tests the file status retrieval functionality.
// It confirms that the file status information is correctly obtained.
func TestStatFile(t *testing.T) {
	filePath := "test_stat.txt"
	defer os.Remove(filePath)

	_, err := CreateFile(filePath, &Options{Perms: 0644})
	assert.NoError(t, err)

	info, err := StatFile(filePath)
	assert.NoError(t, err)
	assert.NotNil(t, info)
}

// TestReadFile tests reading content from a file.
// It ensures that the content read is the same as what was written.
func TestReadFile(t *testing.T) {
	filePath := "test_read.txt"
	content := "Test content"
	defer os.Remove(filePath)

	err := WriteFile(filePath, content)
	assert.NoError(t, err)

	readContent, err := ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, content, string(readContent))
}

// TestWriteFile tests the file writing functionality.
// It verifies that the content written to a file can be read back accurately.
func TestWriteFile(t *testing.T) {
	filePath := "test_write.txt"
	content := "New content"
	defer os.Remove(filePath)

	err := WriteFile(filePath, content)
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

	err := WriteFile(filePath, content)
	assert.NoError(t, err)

	contains, err := Contains(filePath, substring)
	assert.NoError(t, err)
	assert.True(t, contains)

	contains, err = Contains(filePath, "Nonexistent")
	assert.NoError(t, err)
	assert.False(t, contains)
}

// TestCreateOrOpenFile tests the createOrOpenFile function in the fs package.
// It ensures that a file is successfully created or opened without errors.
func TestCreateOrOpenFile(t *testing.T) {
	defer os.RemoveAll("/tmp/kitsune")

	defer os.Remove(VALID_FILE_PATH)

	file, err := createOrOpenFile(VALID_FILE_PATH, os.O_RDWR|os.O_CREATE, &Options{Perms: 0644})
	assert.NoError(t, err)
	assert.NotNil(t, file)
	file.Close()

	file, err = createOrOpenFile(VALID_FILE_PATH, os.O_RDWR, &Options{Perms: 0644})
	assert.NoError(t, err)
	assert.NotNil(t, file)
	file.Close()
}

// TestCreateOrOpenFileError tests error handling in the createOrOpenFile function.
// It verifies that the function returns an error when file creation or opening fails.
func TestCreateOrOpenFileError(t *testing.T) {
	defer os.RemoveAll("/tmp/kitsune")

	file, err := createOrOpenFile(INVALID_FILE_PATH, os.O_RDWR|os.O_CREATE, &Options{Perms: 0644})
	assert.Error(t, err)
	assert.Nil(t, file)
}
