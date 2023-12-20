// Package daemon provides functionalities for handling process identifiers (PID).
package daemon

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/kodflow/kitsune/src/internal/kernel/storages/fs"
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
func (p *PIDHandler) SetPID() error {
	return fs.WriteFile(p.getPIDFilePath(), strconv.Itoa(os.Getpid()))
}

// ClearPID deletes the PID file associated with the process.
// It removes the PID file, indicating that the process is no longer running.
//
// Returns:
// - error: An error if the file operation fails.
func (p *PIDHandler) ClearPID() error {
	return fs.DeleteFile(p.getPIDFilePath())
}

func (p *PIDHandler) GetPID() (int, error) {
	pidBytes, err := fs.ReadFile(p.getPIDFilePath())
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(strings.TrimSpace(string(pidBytes)))
}

// isProcessRunning checks if the process with the given PID is running.
// It verifies the existence of a process by its PID.
//
// Parameters:
// - pid: int The PID of the process to check.
//
// Returns:
// - bool: true if the process is running, false otherwise.
func isProcessRunning(pid int) (bool, error) {
	// Find the process by PID
	process, _ := findProcessByID(pid)
	// Send a signal to check if the process is still running
	err := process.Signal(syscall.Signal(0))
	if err != nil && err.Error() == "os: process already finished" {
		return false, err
	}

	return err == nil, nil
}

// findProcessByID searches and returns a pointer to the process corresponding to the given ID.
// If no matching process is found, an error is returned.
//
// Parameters:
// - id: int The ID of the process to find.
//
// Returns:
// - process: *Process A pointer to the process if found, otherwise nil.
// - error An error if no matching process is found.
func findProcessByID(id int) (*os.Process, error) {
	return os.FindProcess(id)
}

// findProcessByName finds processes by their name.
//
// This function searches for processes by their name using the 'ps' command,
// which is available on most Unix-like systems. It returns all processes
// that match the name. The function uses 'ps' to obtain the process IDs (PIDs)
// and then uses os.FindProcess to create process objects for each PID.
//
// Parameters:
// - processName: string The name of the process to search for.
//
// Returns:
// - []*os.Process: A slice of process objects matching the given name.
func findProcessByName(processName string) []*os.Process {
	// Execute the 'ps' command to get details of processes.
	cmd := exec.Command("ps", "-e", "-o", "pid,comm")
	var out bytes.Buffer
	cmd.Stdout = &out
	_ = cmd.Run() // Ignoring error for simplicity

	var processes []*os.Process
	// Parse the output to find the processes.
	for _, line := range strings.Split(out.String(), "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[1] == processName {
			pid, err := strconv.Atoi(fields[0])
			if err != nil {
				continue // Skip lines that don't have a valid PID.
			}
			process, err := os.FindProcess(pid)
			if err != nil {
				continue // Skip if the process cannot be found.
			}
			processes = append(processes, process)
		}
	}

	return processes
}

// findProcessName retrieves the name of the process given its PID.
//
// This function uses the 'ps' command, available on most Unix-like systems,
// to find the name of the process associated with the provided PID. The function
// executes 'ps' with the specified PID and parses the output to extract the process name.
//
// Parameters:
// - pid: int The Process ID of the process whose name is to be found.
//
// Returns:
// - string: The name of the process associated with the given PID.
// - error: An error object if any error occurs during the process name retrieval.
func findProcessName(pid int) string {
	// Execute the 'ps' command to get the name of the process for the given PID.
	cmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "comm=")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return ""
	}

	processName := strings.TrimSpace(out.String())
	return processName
}
