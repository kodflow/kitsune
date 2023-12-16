package fs

import "os"

// CreateDirectory creates a directory at the specified path.
// If the directory already exists, it does nothing.
// The options parameter allows specifying additional options for creating the directory.
// Returns an error if the directory creation fails.
func CreateDirectory(dirPath string, options ...*Options) error {
	if ExistsDirectory(dirPath) {
		return nil
	}

	opts, err := resolveDirectoryOptions(options...)
	if err != nil {
		return err
	}

	opts.AddPerms(0111)

	err = os.MkdirAll(dirPath, opts.Perms)
	if err != nil {
		return err
	}

	return perms(dirPath, opts)
}

// ExistsDirectory checks if a directory exists at the specified path.
// Returns true if the directory exists, false otherwise.
func ExistsDirectory(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}

// DeleteDirectory deletes a directory at the specified path.
// Returns an error if the directory deletion fails.
func DeleteDirectory(dirPath string) error {
	return os.RemoveAll(dirPath)
}
