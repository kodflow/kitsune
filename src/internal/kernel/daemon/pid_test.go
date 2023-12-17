package daemon_test

import (
	"os"
	"testing"

	"github.com/kodflow/kitsune/src/internal/kernel/daemon"
	"github.com/stretchr/testify/assert"
)

const testPathRun = "/tmp"
const testProcessName = "testProcess"

func TestPIDHandler(t *testing.T) {
	handler := daemon.NewPIDHandler(testProcessName, testPathRun)

	// Test SetPID
	err := handler.SetPID(os.Getpid())
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
	assert.True(t, daemon.IsProcessRunning(os.Getpid()))

	// Using invalid PID
	assert.False(t, daemon.IsProcessRunning(999999))
}
