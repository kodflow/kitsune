package permission

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCheckSuccess tests the Check function when the file has the required permissions.
func TestCheckSuccess(t *testing.T) {
	// Create a temporary directory
	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	// Set the permissions of the directory to 0700
	err = os.Chmod(dir, 0700)
	if err != nil {
		t.Fatalf("failed to set permissions: %v", err)
	}

	// Check if the directory has the required permissions
	err = Check(dir, 0700)
	if err != nil {
		t.Errorf("failed to check permissions: %v", err)
	}
}

// TestCheckFail tests the Check function when the file does not have the required permissions.
func TestCheckFail(t *testing.T) {
	// Check if a fake path has the required permissions
	if err := Check("/fake/path", 0700); err == nil {
		t.Errorf("expected to fail permission check, but passed")
	} else {
		t.Log("expected failure, error:", err)
	}
}

// TestCheckInvalidPermissions tests the Check function when the file has invalid permissions.
func TestCheckInvalidPermissions(t *testing.T) {
	// Create a temporary directory
	dir, err := os.MkdirTemp("", "test-permission")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	// Set the permissions of the directory to 0755
	err = os.Chmod(dir, 0755)
	if err != nil {
		t.Fatalf("failed to set permissions: %v", err)
	}

	// Check if the directory has the required permissions
	err = Check(dir, 0700)
	if err == nil {
		currentPerms, _ := os.Stat(dir)
		t.Errorf("expected to fail permission check for invalid permissions, but passed. Current permissions: %v", currentPerms.Mode().Perm())
	} else {
		t.Log("expected failure for invalid permissions, error:", err)
	}
}

// TestHasMode tests the HasMode function.
func TestHasMode(t *testing.T) {
	// Test case 1: current mode has the required mode
	currentMode := os.FileMode(0700)
	requiredMode := os.FileMode(0700)
	assert.True(t, HasMode(currentMode, requiredMode))

	// Test case 2: current mode does not have the required mode
	currentMode = os.FileMode(0655)
	requiredMode = os.FileMode(0700)
	assert.False(t, HasMode(currentMode, requiredMode))
}
