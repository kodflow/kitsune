package daemon_test

import (
	"os"
	"testing"

	"github.com/kodmain/kitsune/src/config"
	"github.com/kodmain/kitsune/src/internal/kernel/daemon"
)

func TestSetAndGetPID(t *testing.T) {
	err := daemon.SetPID()
	if err != nil {
		t.Fatalf("Failed to set PID: %s", err)
	}

	process, err := daemon.GetPID(config.BUILD_APP_NAME)
	if err != nil {
		t.Fatalf("Error getting PID: %s", err)
	}

	if process == nil {
		t.Fatal("Expected a process, got nil")
	}

	currentPID := os.Getpid()
	if process.Pid != currentPID {
		t.Fatalf("Expected PID %d, got %d", currentPID, process.Pid)
	}
}

func TestClearPID(t *testing.T) {
	daemon.SetPID()
	err := daemon.ClearPID(config.BUILD_APP_NAME)
	if err != nil {
		t.Fatalf("Failed to clear PID: %s", err)
	}

	_, err = daemon.GetPID(config.BUILD_APP_NAME)
	if err == nil {
		t.Fatal("Expected error due to missing PID, got nil")
	}
}
