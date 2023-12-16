package permission

import (
	"fmt"
	"os"
)

// Check checks if the file at the specified path has the required permissions.
// It retrieves the current permissions of the file and compares them with the needed permissions.
// An error is returned if the file does not have the required permissions.
//
// Parameters:
// - path: string The path of the file to check permissions for.
// - need: os.FileMode The required file mode (permissions) for the file.
//
// Returns:
// - error: An error if the file does not have the required permissions or if an error occurs in retrieving file information.
func Check(path string, need os.FileMode) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("cannot retrieve information about the file: %v", err)
	}

	current := fileInfo.Mode().Perm()

	if !HasStrictMode(current, need) {
		return fmt.Errorf("the path %s does not have the required permissions: %v, actual permissions: %v", path, need, current)
	}

	return nil
}

// HasMode checks if the file mode has the specified permissions.
// This function checks for the presence of all the permissions specified in 'need' within the 'current' mode.
//
// Parameters:
// - current: os.FileMode The current file mode.
// - need: os.FileMode The permissions to check for in the current mode.
//
// Returns:
// - bool: true if the current mode includes all the specified permissions, false otherwise.
func HasMode(current, need os.FileMode) bool {
	return current&need == need
}

// HasStrictMode checks if the file mode has exactly the specified permissions.
// This function checks if the current mode matches the specified mode exactly.
//
// Parameters:
// - current: os.FileMode The current file mode.
// - need: os.FileMode The exact permissions to match against the current mode.
//
// Returns:
// - bool: true if the current mode matches the specified mode exactly, false otherwise.
func HasStrictMode(current, need os.FileMode) bool {
	return current == need
}
