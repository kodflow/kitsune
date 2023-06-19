package daemon

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kodmain/KitsuneFramework/internal/env"
	"github.com/kodmain/KitsuneFramework/internal/storages/fs"
)

var pidFilePath string = filepath.Join(env.PATH_PID, env.BUILD_APP_NAME+".pid")

func SetPID(pidFilePath string) error {
	return fs.WriteFile(pidFilePath, strconv.Itoa(os.Getpid()))
}

func GetPID(pidFilePath string) (*os.Process, error) {
	var process *os.Process = nil
	var err error = nil

	if !fs.ExistsFile(pidFilePath) {
		return process, err // PID file does not exist
	}

	pidBytes, err := fs.ReadFile(pidFilePath)
	if err != nil {
		return process, err
	}

	pidStr := strings.TrimSpace(string(pidBytes))
	pid, _ := strconv.Atoi(pidStr)
	process, _ = os.FindProcess(pid)

	if process == nil {
		if err := fs.DeleteFile(pidFilePath); err != nil {
			return process, fmt.Errorf("can't read process on pid file")
		}
		return process, err
	}

	return process, fmt.Errorf("process is already running")
}
