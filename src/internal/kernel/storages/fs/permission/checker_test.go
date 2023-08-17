package permission

import (
	"os"
	"testing"
)

func TestCheck_Positive(t *testing.T) {
	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	err = os.Chmod(dir, 0700)
	if err != nil {
		t.Fatalf("failed to set permissions: %v", err)
	}

	err = Check(dir)
	if err != nil {
		t.Errorf("failed to check permissions: %v", err)
	}
}

func TestCheck_Negative(t *testing.T) {
	if err := Check("/fake/path"); err != nil {
		t.Log("expected failure, error:", err)
	} else {
		t.Errorf("expected to fail permission check, but passed")
	}
}
