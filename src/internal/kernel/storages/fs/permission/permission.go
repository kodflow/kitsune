package permission

import (
	"io/fs"
	"os"
)

// Constants representing file permissions.
const (
	read      = 0400 // Read permission
	write     = 0200 // Write permission
	execution = 0100 // Execution permission

	PERMS_R   = os.FileMode(read)                     // Read permission
	PERMS_W   = os.FileMode(write)                    // Write permission
	PERMS_X   = os.FileMode(execution)                // Execution permission
	PERMS_RW  = os.FileMode(read | write)             // Read and write permissions
	PERMS_RX  = os.FileMode(read | execution)         // Read and execution permissions
	PERMS_WX  = os.FileMode(write | execution)        // Write and execution permissions
	PERMS_RWX = os.FileMode(read | write | execution) // Read, write, and execution permissions
)

// Validate checks if the file at the specified filePath has the given permissions.
// It returns true if the file has the specified permissions, otherwise it returns false.
func Validate(filePath string, perms fs.FileMode) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	return info.Mode().Perm()&perms == perms
}
