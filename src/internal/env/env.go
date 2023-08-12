package env

import (
	"os"
	"os/user"
)

var (
	USER          *user.User
	GROUP         *user.Group
	HOME          string
	CWD           string
	HOSTNAME      string
	SERVICE_TOKEN string
)

func init() {
	if user, err := user.Current(); err == nil {
		USER = user
	}

	if group, err := user.LookupGroupId(USER.Gid); err == nil {
		GROUP = group
	}

	if hostname, err := os.Hostname(); err == nil {
		HOSTNAME = hostname
	}

	if home, err := os.UserHomeDir(); err == nil {
		HOME = home
	}

	if cwd, err := os.Getwd(); err == nil {
		CWD = cwd
	}

	SERVICE_TOKEN = HOSTNAME + "/" + BUILD_APP_NAME
}
