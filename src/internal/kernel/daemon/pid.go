package daemon

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/kodmain/kitsune/src/config"
	"github.com/kodmain/kitsune/src/internal/kernel/storages/fs"
)

func getPIDFilePath(processName string) string {
	return filepath.Join(config.PATH_RUN, processName+".pid")
}

// SetPID writes the PID of the current process to a file
func SetPID() error {
	return fs.WriteFile(getPIDFilePath(config.BUILD_APP_NAME), strconv.Itoa(os.Getpid()))
}

// ClearProcess kills the given process and clears its PID file
func ClearProcess(process *os.Process, name string) error {
	if err := process.Kill(); err != nil {
		return err
	}

	return ClearPID(name)
}

// ClearPID deletes the PID file associated with the given process name
func ClearPID(name string) error {
	return fs.DeleteFile(getPIDFilePath(name))
}

// GetPID returns the os.Process if the process is running, otherwise it returns an error
func GetPID(processName string) (*os.Process, error) {
	if !fs.ExistsFile(getPIDFilePath(processName)) {
		return nil, nil
	}

	pidBytes, err := fs.ReadFile(getPIDFilePath(processName))
	if err != nil {
		return nil, err
	}

	pidStr := strings.TrimSpace(string(pidBytes))
	pid, _ := strconv.Atoi(pidStr)
	process, _ := os.FindProcess(pid)

	// Check if the process is actually running
	if isProcessRunning(pid) {
		return process, fmt.Errorf("process is already running")
	}

	// If process isn't running but PID file exists, delete the PID file
	if err := fs.DeleteFile(getPIDFilePath(processName)); err != nil {
		return nil, fmt.Errorf("can't read process on pid file")
	}

	return nil, nil
}

// isProcessRunning checks if a process with the given PID is currently running
func isProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return process.Signal(syscall.Signal(0)) == nil
}
