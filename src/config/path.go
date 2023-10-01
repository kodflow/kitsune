package config

import (
	"os"
	"path/filepath"

	"github.com/kodmain/kitsune/src/internal/kernel/storages/fs/permission"
)

var (
	PATH_SERVICES = "/etc/kitsune/"
	PATH_BIN      = "/usr/local/bin/"
	PATH_RUN      = "/var/run/kitsune/"
	PATH_LOGS     = "/var/log/kitsune/"
)

var PATHS = []*string{
	&PATH_SERVICES,
	&PATH_RUN,
	&PATH_BIN,
	&PATH_LOGS,
}

func init() {
	var hasPerms = true
	for _, path := range PATHS {
		if err := permission.Check(*path, 0755); err != nil {
			hasPerms = false
			break
		}
	}

	if !hasPerms {
		homeDir, _ := os.UserHomeDir()
		kitsunePath := filepath.Join(homeDir, ".kitsune")
		for _, path := range PATHS {
			*path = filepath.Join(kitsunePath, *path)
		}
	}
}
