package components

import (
	"testing"
)

func TestNewDirPicker(t *testing.T) {
	startDir := "/home/user/documents"
	defaultDir := "/home/user/.bento/bentos"

	dp := NewDirPicker(startDir, defaultDir)

	// Verify ShowHidden is enabled
	if !dp.ShowHidden {
		t.Error("ShowHidden should be true to show hidden directories like .bento")
	}

	// Verify DirAllowed is true
	if !dp.DirAllowed {
		t.Error("DirAllowed should be true for directory picker")
	}

	// Verify FileAllowed is false
	if dp.FileAllowed {
		t.Error("FileAllowed should be false for directory picker")
	}

	// Verify current directory is set
	if dp.CurrentDirectory != startDir {
		t.Errorf("CurrentDirectory = %s, want %s", dp.CurrentDirectory, startDir)
	}

	// Verify default directory is stored
	if dp.defaultDir != defaultDir {
		t.Errorf("defaultDir = %s, want %s", dp.defaultDir, defaultDir)
	}

	// Note: Height is set via WithHeight() method, no public field to verify
}

func TestResetToDefault(t *testing.T) {
	startDir := "/some/random/path"
	defaultDir := "/home/user/.bento"

	dp := NewDirPicker(startDir, defaultDir)

	// Verify starts at startDir
	if dp.CurrentDirectory != startDir {
		t.Errorf("Initial CurrentDirectory = %s, want %s", dp.CurrentDirectory, startDir)
	}

	// Reset to default
	dp = dp.ResetToDefault()

	// Verify now at defaultDir
	if dp.CurrentDirectory != defaultDir {
		t.Errorf("After reset CurrentDirectory = %s, want %s", dp.CurrentDirectory, defaultDir)
	}
}

func TestDirPickerSeparatesStartAndDefault(t *testing.T) {
	// This tests the fix where we separate "where to start" from "where to reset"
	userCurrentDir := "/home/user/Desktop/Desktop/Desktop"
	appDefaultDir := "/home/user/.bento"

	dp := NewDirPicker(userCurrentDir, appDefaultDir)

	// Should start at user's current directory
	if dp.CurrentDirectory != userCurrentDir {
		t.Errorf("Should start at user dir %s, got %s", userCurrentDir, dp.CurrentDirectory)
	}

	// But reset should go to app default, not user's current
	dp = dp.ResetToDefault()
	if dp.CurrentDirectory != appDefaultDir {
		t.Errorf("Reset should go to app default %s, got %s", appDefaultDir, dp.CurrentDirectory)
	}
}
