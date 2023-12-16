package fs

import "os"

// CreateDirectory creates a directory at the specified path.
// If the directory already exists, it does nothing.
// This function allows specifying additional options for creating the directory.
// It also adds execute permissions to allow directory traversal.
//
// Parameters:
// - dirPath: string The path where the directory will be created.
// - options: []*Options Optional parameters to specify directory options.
//
// Returns:
// - error: An error if the directory creation fails.
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
// It uses os.Stat to get file information and determines if the path is a directory.
//
// Parameters:
// - dirPath: string The path of the directory to check.
//
// Returns:
// - bool: true if the directory exists, false otherwise.
func ExistsDirectory(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}

// DeleteDirectory deletes a directory at the specified path.
// It uses os.RemoveAll to delete the directory and all its contents.
//
// Parameters:
// - dirPath: string The path of the directory to delete.
//
// Returns:
// - error: An error if the directory deletion fails.
func DeleteDirectory(dirPath string) error {
	return os.RemoveAll(dirPath)
}
