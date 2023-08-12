package env

import (
	"fmt"
	"os"
	"os/user"
	"runtime"
)

var (
	USER          *user.User
	CWD           string
	HOSTNAME      string
	SERVICE_TOKEN string
)

func init() {

	fmt.Println("runtime", runtime.GOOS)
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
