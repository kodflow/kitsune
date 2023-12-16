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
// It is used to manage PID files for processes, typically in a Unix-like environment.
type PIDHandler struct {
	// processName represents the name of the process for which the PID file is managed.
	processName string

	// pathRun represents the path to the directory containing the PID file.
	pathRun string
}

// NewPIDHandler creates and returns a new instance of PIDHandler for a given process name and path.
// It initializes a PIDHandler with a specified process name and the path for the PID file.
//
// Parameters:
// - processName: string The name of the process.
// - pathRun: string The path to the directory containing the PID file.
//
// Returns:
// - *PIDHandler: A new instance of PIDHandler.
func NewPIDHandler(processName, pathRun string) *PIDHandler {
	return &PIDHandler{processName: processName, pathRun: pathRun}
}

// getPIDFilePath constructs the full path of the PID file based on the process name and path.
// This is an internal helper method to build the path for the PID file.
//
// Returns:
// - string: The full file path for the PID file.
func (p *PIDHandler) getPIDFilePath() string {
	return filepath.Join(p.pathRun, p.processName+".pid")
}

// SetPID creates or updates the PID file with the provided PID.
// It writes the given PID to the PID file for the process.
//
// Parameters:
// - pid: int The process ID to set in the PID file.
//
// Returns:
// - error: An error if the file operation fails.
func (p *PIDHandler) SetPID(pid int) error {
	return fs.WriteFile(p.getPIDFilePath(), strconv.Itoa(pid))
}

// ClearPID deletes the PID file associated with the process.
// It removes the PID file, indicating that the process is no longer running.
//
// Returns:
// - error: An error if the file operation fails.
func (p *PIDHandler) ClearPID() error {
	return fs.DeleteFile(p.getPIDFilePath())
}

// GetPID retrieves the PID from the PID file.
// If the associated process is running, it returns the PID and an error.
// If the PID file exists but the process isn't running, it deletes the PID file and returns 0.
//
// Returns:
// - int: The PID of the running process or 0 if not running.
// - error: An error if the process is already running or if there's an issue reading the PID file.
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
// It verifies the existence of a process by its PID.
//
// Parameters:
// - pid: int The PID of the process to check.
//
// Returns:
// - bool: true if the process is running, false otherwise.
func IsProcessRunning(pid int) bool {
	// Find the process by PID
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// Send a signal to check if the process is still running
	err = process.Signal(syscall.Signal(0))
	if err != nil && err.Error() == "os: process already finished" {
		return false
	}

	return err == nil
}
