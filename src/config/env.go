package config

import (
	"os"
	"os/user"
)

// Global variables representing various system and user-specific information.
var (
	USER          *user.User // The current user running the application.
	CWD           string     // The current working directory of the application.
	HOSTNAME      string     // The hostname of the system where the application is running.
	SERVICE_TOKEN string     // A token generated using the hostname and application build name.
)

// init is automatically called when the package is imported.
// It initializes the global variables with the current user, current working directory,
// and hostname of the system. It also constructs the SERVICE_TOKEN.
func init() {
	// Retrieve and set the current user
	user, err := user.Current()
	if err == nil {
		USER = user
	}

	// Retrieve and set the hostname of the system
	if hostname, err := os.Hostname(); err == nil {
		HOSTNAME = hostname
	}

	// Retrieve and set the current working directory
	if cwd, err := os.Getwd(); err == nil {
		CWD = cwd
	}

	// Construct the SERVICE_TOKEN using the hostname and BUILD_APP_NAME.
	// Note: BUILD_APP_NAME needs to be defined elsewhere in your application.
	SERVICE_TOKEN = HOSTNAME + "/" + BUILD_APP_NAME
}
