package config

import (
	"fmt"
	"os"
	"os/user"
)

var (
	USER          *user.User
	CWD           string
	HOSTNAME      string
	SERVICE_TOKEN string
)

func init() {

	user, err := user.Current()
	if err == nil {
		USER = user
	} else {
		fmt.Println(err)
	}

	if hostname, err := os.Hostname(); err == nil {
		HOSTNAME = hostname
	}

	if cwd, err := os.Getwd(); err == nil {
		CWD = cwd
	}

	SERVICE_TOKEN = HOSTNAME + "/" + BUILD_APP_NAME
}
