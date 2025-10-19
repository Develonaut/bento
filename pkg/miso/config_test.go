// Package miso provides terminal output "seasoning" - themed styling and progress display.
//
// Tests for theme configuration persistence.
package miso

import (
	"os"
	"path/filepath"
	"testing"
)

// TestConfigDir verifies config directory path.
func TestConfigDir(t *testing.T) {
	dir, err := configDir()
	if err != nil {
		t.Fatalf("configDir() failed: %v", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("UserHomeDir() failed: %v", err)
	}

	expected := filepath.Join(home, ".bento")
	if dir != expected {
		t.Errorf("configDir() = %s, want %s", dir, expected)
	}
}

// TestThemeConfigPath verifies theme config file path.
func TestThemeConfigPath(t *testing.T) {
	path, err := themeConfigPath()
	if err != nil {
		t.Fatalf("themeConfigPath() failed: %v", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("UserHomeDir() failed: %v", err)
	}

	expected := filepath.Join(home, ".bento", "theme")
	if path != expected {
		t.Errorf("themeConfigPath() = %s, want %s", path, expected)
	}
}

// TestSaveAndLoadTheme verifies theme persistence round-trip.
func TestSaveAndLoadTheme(t *testing.T) {
	// Use temp directory to avoid polluting user's actual config
	tmpDir := t.TempDir()
	originalConfigDir := configDir

	// Mock configDir to use temp directory
	configDir = func() (string, error) {
		return tmpDir, nil
	}
	t.Cleanup(func() {
		configDir = originalConfigDir
	})

	// Test each variant
	for _, variant := range AllVariants() {
		t.Run(string(variant), func(t *testing.T) {
			// Save theme
			if err := SaveTheme(variant); err != nil {
				t.Fatalf("SaveTheme(%s) failed: %v", variant, err)
			}

			// Load theme
			loaded := LoadSavedTheme()
			if loaded != variant {
				t.Errorf("LoadSavedTheme() = %s, want %s", loaded, variant)
			}
		})
	}
}

// TestLoadSavedTheme_NoFile verifies default when no theme file exists.
func TestLoadSavedTheme_NoFile(t *testing.T) {
	tmpDir := t.TempDir()
	originalConfigDir := configDir

	configDir = func() (string, error) {
		return tmpDir, nil
	}
	t.Cleanup(func() {
		configDir = originalConfigDir
	})

	// Should return Tonkotsu default when no file exists
	variant := LoadSavedTheme()
	if variant != VariantTonkotsu {
		t.Errorf("LoadSavedTheme() with no file = %s, want %s", variant, VariantTonkotsu)
	}
}

// TestLoadSavedTheme_InvalidContent verifies default for invalid content.
func TestLoadSavedTheme_InvalidContent(t *testing.T) {
	tmpDir := t.TempDir()
	originalConfigDir := configDir

	configDir = func() (string, error) {
		return tmpDir, nil
	}
	t.Cleanup(func() {
		configDir = originalConfigDir
	})

	// Write invalid content
	path := filepath.Join(tmpDir, "theme")
	if err := os.WriteFile(path, []byte("InvalidVariant"), 0644); err != nil {
		t.Fatal(err)
	}

	// Should return Tonkotsu default
	variant := LoadSavedTheme()
	if variant != VariantTonkotsu {
		t.Errorf("LoadSavedTheme() with invalid content = %s, want %s", variant, VariantTonkotsu)
	}
}

// TestSaveTheme_CreatesDirectory verifies directory creation.
func TestSaveTheme_CreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	tmpSubDir := filepath.Join(tmpDir, "nonexistent")
	originalConfigDir := configDir

	configDir = func() (string, error) {
		return tmpSubDir, nil
	}
	t.Cleanup(func() {
		configDir = originalConfigDir
	})

	// Directory should not exist yet
	if _, err := os.Stat(tmpSubDir); err == nil {
		t.Fatal("Directory should not exist yet")
	}

	// SaveTheme should create it
	if err := SaveTheme(VariantNasu); err != nil {
		t.Fatalf("SaveTheme() failed: %v", err)
	}

	// Directory should now exist
	if _, err := os.Stat(tmpSubDir); err != nil {
		t.Fatalf("Directory was not created: %v", err)
	}
}
