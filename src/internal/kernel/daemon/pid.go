package daemon

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kodmain/kitsune/src/internal/env"
	"github.com/kodmain/kitsune/src/internal/storages/fs"
)

func getPIDFilePath(processName string) string {
	return filepath.Join(env.PATH_RUN, processName+".pid")
}

func SetPID() error {
	return fs.WriteFile(getPIDFilePath(env.BUILD_APP_NAME), strconv.Itoa(os.Getpid()))
}

func ClearProcess(process *os.Process, name string) error {
	if err := process.Kill(); err != nil {
		return err
	}

	return ClearPID(name)
}

func ClearPID(name string) error {
	return fs.DeleteFile(getPIDFilePath(name))
}

func GetPID(processName string) (*os.Process, error) {
	var process *os.Process = nil
	var err error = nil

	if !fs.ExistsFile(getPIDFilePath(processName)) {
		return process, err
	}

	pidBytes, err := fs.ReadFile(getPIDFilePath(processName))
	if err != nil {
		return process, err
	}

	pidStr := strings.TrimSpace(string(pidBytes))
	pid, _ := strconv.Atoi(pidStr)
	process, _ = os.FindProcess(pid)

	if process == nil {
		if err := fs.DeleteFile(getPIDFilePath(processName)); err != nil {
			return process, fmt.Errorf("can't read process on pid file")
		}
		return process, err
	}

	return process, fmt.Errorf("process is already running")
}
