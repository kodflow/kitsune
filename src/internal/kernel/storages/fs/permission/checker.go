package permission

import (
	"fmt"
	"os"
)

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

// HasStrict checks if the file mode has exactly the specified permissions.
func HasStrictMode(current, need os.FileMode) bool {
	return current == need
}
