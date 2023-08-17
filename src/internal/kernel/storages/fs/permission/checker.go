package permission

import (
	"fmt"
	"os"
)

func Check(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("cannot retrieve information about the file: %v", err)
	}

	perm := fileInfo.Mode().Perm()

	if fileInfo.IsDir() {
		if perm&0700 != 0700 {
			return fmt.Errorf("you don't have read, write and execute permission to this directory")
		}
	} else {
		if perm&0600 != 0600 {
			return fmt.Errorf("you don't have read and write permission to this file")
		}
	}

	return nil
}
