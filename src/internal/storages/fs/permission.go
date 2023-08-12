package fs

import "os"

const (
	read      = 0400
	write     = 0200
	execution = 0100

	PERMS_R   = os.FileMode(read)
	PERMS_W   = os.FileMode(write)
	PERMS_X   = os.FileMode(execution)
	PERMS_RW  = os.FileMode(read | write)
	PERMS_RX  = os.FileMode(read | execution)
	PERMS_WX  = os.FileMode(write | execution)
	PERMS_RWX = os.FileMode(read | write | execution)
)

func Permission(filePath string, perms os.FileMode) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	return info.Mode().Perm()&perms == perms
}
