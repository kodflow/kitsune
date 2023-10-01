// Package daemon provides functionalities for handling process identifiers (PID).
package daemon

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/kodmain/kitsune/src/internal/kernel/storages/fs"
)

// PIDHandler handles the creation, retrieval, and deletion of PID files.
type PIDHandler struct {
	// processName represents the name of the process for which the PID file is managed.
	processName string

	// pathRun represents the path to the directory containing the PID file.
	pathRun string
}

// NewPIDHandler creates and returns a new instance of PIDHandler for a given process name and path.
func NewPIDHandler(processName, pathRun string) *PIDHandler {
	return &PIDHandler{processName: processName, pathRun: pathRun}
}

// getPIDFilePath constructs the full path of the PID file based on the process name and path.
func (p *PIDHandler) getPIDFilePath() string {
	return filepath.Join(p.pathRun, p.processName+".pid")
}

// SetPID creates or updates the PID file with the provided PID.
func (p *PIDHandler) SetPID(pid int) error {
	return fs.WriteFile(p.getPIDFilePath(), strconv.Itoa(pid))
}

// ClearPID deletes the PID file associated with the process.
func (p *PIDHandler) ClearPID() error {
	return fs.DeleteFile(p.getPIDFilePath())
}

// GetPID retrieves the PID from the PID file.
// If the associated process is running, it returns the PID and an error indicating the process is already running.
// If the PID file exists but the process isn't running, it deletes the PID file and returns 0 without error.
func (p *PIDHandler) GetPID() (int, error) {
	if !fs.ExistsFile(p.getPIDFilePath()) {
		return 0, nil
	}

	pidBytes, err := fs.ReadFile(p.getPIDFilePath())
	if err != nil {
		return 0, err
	}

	pidStr := strings.TrimSpace(string(pidBytes))
	pid, _ := strconv.Atoi(pidStr)

	// Check if the process is effectively running
	if IsProcessRunning(pid) {
		return pid, fmt.Errorf("process is already running")
	}

	// If the process isn't running but the PID file exists, delete the PID file
	if err := p.ClearPID(); err != nil {
		return 0, fmt.Errorf("can't read process on pid file")
	}

	return 0, nil
}

// IsProcessRunning checks if the process with the given PID is running.
func IsProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return process.Signal(syscall.Signal(0)) == nil
}
