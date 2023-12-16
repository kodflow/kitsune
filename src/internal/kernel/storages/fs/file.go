package fs

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// CreateFile creates a file at the specified file path with the given options.
//
// Parameters:
// - filePath: The path of the file to create.
// - options: The options to apply when creating the file.
//
// Returns:
// - *os.File: The created file.
// - error: An error if any occurred during the process.
func CreateFile(filePath string, options ...*Options) (*os.File, error) {
	return createOrOpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, options...)
}

// Permissions sets the permissions of the file at the specified file path with the given options.
//
// Parameters:
// - filePath: The path of the file to set permissions for.
// - options: The options to apply when setting permissions.
//
// Returns:
// - error: An error if any occurred during the process.
func Permissions(filePath string, options ...*Options) error {
	opts, err := resolveFileOptions(options...)
	if err != nil {
		return err
	}

	return perms(filePath, opts)
}

func OpenFile(filePath string, options ...*Options) (*os.File, error) {
	return createOrOpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, options...)
}

// createOrOpenFile creates or opens a file at the specified file path with the given flags and options.
// It resolves the file options, creates the necessary directory, opens the file, sets the permissions,
// and returns the opened file or an error if any occurred.
//
// Parameters:
// - filePath: The path of the file to create or open.
// - flag: The flags to use when opening the file.
// - options: The options to apply when creating or opening the file.
//
// Returns:
// - *os.File: The opened file.
// - error: An error if any occurred during the process.
func createOrOpenFile(filePath string, flag int, options ...*Options) (*os.File, error) {
	opts, err := resolveFileOptions(options...)
	fmt.Println(opts)

	if err != nil {
		return nil, err
	}

	err = CreateDirectory(filepath.Dir(filePath), opts)
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
//
// Parameters:
// - filePath: The path of the file to delete.
//
// Returns:
// - error: An error if any occurred during the process.
func DeleteFile(filePath string) error {
	return os.Remove(filePath)
}

// SHA1File calculates the SHA1 hash of the file at the specified file path.
//
// Parameters:
// - filePath: The path of the file to calculate the SHA1 hash for.
//
// Returns:
// - string: The SHA1 hash of the file.
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

// ExistsFile checks if the file at the specified file path exists.
//
// Parameters:
// - filePath: The path of the file to check.
//
// Returns:
// - bool: true if the file exists, false otherwise.
func ExistsFile(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// StatFile returns the file information for the file at the specified file path.
//
// Parameters:
// - filePath: The path of the file to get information for.
//
// Returns:
// - fs.FileInfo: The file information.
// - error: An error if any occurred during the process.
func StatFile(filePath string) (fs.FileInfo, error) {
	return os.Stat(filePath)
}

// ReadFile reads the content of the file at the specified file path.
//
// Parameters:
// - filePath: The path of the file to read.
//
// Returns:
// - []byte: The content of the file.
// - error: An error if any occurred during the process.
func ReadFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

// WriteFile writes the content to the file at the specified file path.
//
// Parameters:
// - filePath: The path of the file to write to.
// - content: The content to write to the file.
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
//
// Parameters:
// - filePath: The path of the file to check.
// - substring: The substring to search for in the file.
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
