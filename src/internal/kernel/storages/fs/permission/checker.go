package permission

import (
	"fmt"
	"os"
)

// Check checks if the file at the specified path has the required permissions.
// It returns an error if the file does not have the required permissions.
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
func HasMode(current, need os.FileMode) bool {
	return current&need == need
}

// HasStrictMode checks if the file mode has exactly the specified permissions.
func HasStrictMode(current, need os.FileMode) bool {
	return current == need
}
