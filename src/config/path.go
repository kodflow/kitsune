package config

import (
	"os"
	"path/filepath"

	"github.com/kodflow/kitsune/src/internal/kernel/storages/fs/permission"
)

// PATH_SERVICES, PATH_BIN, PATH_RUN, and PATH_LOGS are constants representing
// default filesystem paths used by the Kitsune application.
var (
	PATH_SERVICES = "/etc/kitsune/"     // Path for storing service-related configurations.
	PATH_BIN      = "/usr/local/bin/"   // Path for binary executables.
	PATH_RUN      = "/var/run/kitsune/" // Path for runtime files, like sockets and PID files.
	PATH_LOGS     = "/var/log/kitsune/" // Path for storing log files.
)

// PATHS is a slice of pointers to the path variables.
// This allows for iterating over and modifying these paths in a unified manner.
var PATHS = []*string{
	&PATH_SERVICES,
	&PATH_RUN,
	&PATH_BIN,
	&PATH_LOGS,
}

// init is called when the package is initialized.
// It checks the permissions of the default paths and, if necessary, changes them
// to a user-specific directory under the user's home directory.
func init() {
	// hasPerms flags if the default paths have the required permissions.
	var hasPerms = true

	// Check permissions for each path. If any path lacks the required permissions,
	// hasPerms is set to false.
	for _, path := range PATHS {
		if err := permission.Check(*path, 0755); err != nil {
			hasPerms = false
			break
		}
	}

	// If the default paths don't have required permissions, set the paths to
	// a user-specific directory under the user's home directory.
	if !hasPerms {
		homeDir, _ := os.UserHomeDir()
		kitsunePath := filepath.Join(homeDir, ".kitsune")
		for _, path := range PATHS {
			*path = filepath.Join(kitsunePath, *path)
		}
	}
}
