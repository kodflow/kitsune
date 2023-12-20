package daemon

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testPathRun = "/tmp"
const testProcessName = "testProcess"

func TestPIDHandler(t *testing.T) {
	handler := NewPIDHandler(testProcessName, testPathRun)

	// Test SetPID
	err := handler.SetPID()
	assert.Nil(t, err)

	// Test GetPID
	pid, err := handler.GetPID()
	assert.NotEqual(t, 0, pid)
	assert.Equal(t, os.Getpid(), pid)
	assert.Error(t, err, "Expected an error since process is already running")

	// Test ClearPID
	err = handler.ClearPID()
	assert.Nil(t, err)

	// Test GetPID again, now expecting 0 PID since it was cleared
	pid, err = handler.GetPID()
	assert.Equal(t, 0, pid)
	assert.Nil(t, err)
}

func TestIsProcessRunning(t *testing.T) {
	// Using PID of the current test process
	isRunning, err := isProcessRunning(os.Getpid())
	assert.True(t, isRunning)
	assert.NoError(t, err)

	// Using invalid PID
	isRunning, err = isProcessRunning(9999999999999)
	assert.False(t, isRunning)
	assert.Equal(t, "os: process already finished", err.Error())
}

// TestFindProcessByNameExisting tests finding an existing process.
func TestFindProcessByNameExisting(t *testing.T) {
	pid := os.Getpid()
	processName := findProcessName(pid)
	assert.NotEmpty(t, processName)
	processes := findProcessByName(processName)
	assert.NotNil(t, processes)
	assert.NotEmpty(t, processes)
}

// TestFindProcessByNameExisting tests finding an existing process.
func TestFindProcessByNameNotExisting(t *testing.T) {
	// Replace "process_name" with a process that is expected to be running on your system.
	processName := "process_name"
	processes := findProcessByName(processName)
	assert.Empty(t, processes)
}
func TestFindProcessName(t *testing.T) {
	pid := os.Getpid()
	processName := findProcessName(pid)
	assert.NotEmpty(t, processName)
}

func TestFindProcessNameError(t *testing.T) {
	invalidPid := 9999999999999
	processName := findProcessName(invalidPid)
	assert.Empty(t, processName)
}
