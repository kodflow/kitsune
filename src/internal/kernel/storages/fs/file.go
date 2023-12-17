package fs

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// CreateFile creates a file at the specified file path with the given options.
// It wraps the createOrOpenFile function with specific flags for file creation.
//
// Parameters:
// - filePath: string The path of the file to create.
// - options: []*Options Optional parameters to specify file options.
//
// Returns:
// - *os.File: The created file.
// - error: An error if any occurred during the process.
func CreateFile(filePath string, options ...*Options) (*os.File, error) {
	return createOrOpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, options...)
}

// OpenFile opens a file at the specified file path with the given options.
// It wraps the createOrOpenFile function with specific flags for opening a file for appending and writing.
//
// Parameters:
// - filePath: string The path of the file to open.
// - options: []*Options Optional parameters to specify file options.
//
// Returns:
// - *os.File: The opened file.
// - error: An error if any occurred during the process.
func OpenFile(filePath string, options ...*Options) (*os.File, error) {
	return createOrOpenFile(filePath, os.O_APPEND|os.O_WRONLY, options...)
}

// Permissions sets the permissions of the file at the specified file path with the given options.
// It uses the resolveFileOptions function to determine the options and sets permissions accordingly.
//
// Parameters:
// - filePath: string The path of the file to set permissions for.
// - options: []*Options Optional parameters to specify file options.
//
// Returns:
// - error: An error if any occurred during the process.
func Permissions(filePath string, options ...*Options) error {
	opts := resolveFileOptions(options...)

	return perms(filePath, opts)
}

// createOrOpenFile creates or opens a file at the specified file path with the given flags and options.
// It resolves file options, creates the necessary directory, opens the file, sets permissions, and returns the file.
//
// Parameters:
// - filePath: string The path of the file to create or open.
// - flag: int The flags to use when opening the file.
// - options: []*Options Optional parameters to specify file options.
//
// Returns:
// - *os.File: The opened file.
// - error: An error if any occurred during the process.
func createOrOpenFile(filePath string, flag int, options ...*Options) (*os.File, error) {
	opts := resolveFileOptions(options...)

	err := CreateDirectory(filepath.Dir(filePath), opts)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(filePath, flag, opts.Perms)
	if err != nil {
		return nil, err
	}

	err = perms(filePath, opts)

	if err != nil {
		file.Close()
		return nil, err
	}

	return file, nil
}

// DeleteFile deletes the file at the specified file path.
// It simply calls os.Remove to delete the file.
//
// Parameters:
// - filePath: string The path of the file to delete.
//
// Returns:
// - error: An error if any occurred during the process.
func DeleteFile(filePath string) error {
	return os.Remove(filePath)
}

// SHA1File calculates the SHA1 hash of the file at the specified file path.
// It opens the file, calculates its SHA1 hash, and returns the hash as a hex string.
//
// Parameters:
// - filePath: string The path of the file to calculate the SHA1 hash for.
//
// Returns:
// - string: The SHA1 hash of the file as a hex string.
// - error: An error if any occurred during the process.
func SHA1File(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// ExistsFile checks if the file at the specified file path exists and is not a directory.
// It uses os.Stat to get file information and checks if the file is not a directory.
//
// Parameters:
// - filePath: string The path of the file to check.
//
// Returns:
// - bool: true if the file exists and is not a directory, false otherwise.
func ExistsFile(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// StatFile returns the file information for the file at the specified file path.
// It calls os.Stat to retrieve file information.
//
// Parameters:
// - filePath: string The path of the file to get information for.
//
// Returns:
// - fs.FileInfo: The file information.
// - error: An error if any occurred during the process.
func StatFile(filePath string) (fs.FileInfo, error) {
	return os.Stat(filePath)
}

// ReadFile reads the content of the file at the specified file path.
// It uses os.ReadFile to read and return the content of the file.
//
// Parameters:
// - filePath: string The path of the file to read.
//
// Returns:
// - []byte: The content of the file.
// - error: An error if any occurred during the process.
func ReadFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

// WriteFile writes the content to the file at the specified file path.
// It creates the necessary directory and then writes the content to the file using os.WriteFile.
//
// Parameters:
// - filePath: string The path of the file to write to.
// - content: string The content to write to the file.
//
// Returns:
// - error: An error if any occurred during the process.
func WriteFile(filePath string, content string) error {
	err := CreateDirectory(filepath.Dir(filePath))
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, []byte(content), 0644)
}

// Contains checks if the file at the specified file path contains the given substring.
// It opens the file, scans it line by line, and checks if any line contains the substring.
//
// Parameters:
// - filePath: string The path of the file to check.
// - substring: string The substring to search for in the file.
//
// Returns:
// - bool: true if the file contains the substring, false otherwise.
// - error: An error if any occurred during the process.
func Contains(filePath string, substring string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), substring) {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}
