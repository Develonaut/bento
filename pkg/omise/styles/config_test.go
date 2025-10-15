package styles

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoadTheme(t *testing.T) {
	// Create a temporary config directory for testing
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")

	// Override HOME for testing
	t.Setenv("HOME", tempDir)
	defer func() {
		os.Setenv("HOME", originalHome)
	}()

	// Test saving a theme
	variant := VariantToro
	err := SaveTheme(variant)
	if err != nil {
		t.Fatalf("SaveTheme failed: %v", err)
	}

	// Verify the file was created
	configPath := filepath.Join(tempDir, ".bento", "theme")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Theme config file was not created")
	}

	// Test loading the saved theme
	loaded := LoadSavedTheme()
	if loaded != variant {
		t.Errorf("Expected loaded theme to be %s, got %s", variant, loaded)
	}
}

func TestLoadSavedTheme_NoFile(t *testing.T) {
	// Create a temporary directory with no theme file
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")

	t.Setenv("HOME", tempDir)
	defer func() {
		os.Setenv("HOME", originalHome)
	}()

	// Should return default variant when no file exists
	loaded := LoadSavedTheme()
	if loaded != VariantMaguro {
		t.Errorf("Expected default theme to be %s, got %s", VariantMaguro, loaded)
	}
}

func TestLoadSavedTheme_InvalidVariant(t *testing.T) {
	// Create a temporary config directory
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")

	t.Setenv("HOME", tempDir)
	defer func() {
		os.Setenv("HOME", originalHome)
	}()

	// Write an invalid variant to the config file
	configDir := filepath.Join(tempDir, ".bento")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	configPath := filepath.Join(configDir, "theme")
	if err := os.WriteFile(configPath, []byte("InvalidVariant"), 0644); err != nil {
		t.Fatalf("Failed to write invalid theme: %v", err)
	}

	// Should return default variant when invalid variant is saved
	loaded := LoadSavedTheme()
	if loaded != VariantMaguro {
		t.Errorf("Expected default theme for invalid variant, got %s", loaded)
	}
}

func TestSaveTheme_CreatesDirectory(t *testing.T) {
	// Create a temporary directory without .bento folder
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")

	t.Setenv("HOME", tempDir)
	defer func() {
		os.Setenv("HOME", originalHome)
	}()

	// Save a theme - should create the directory
	variant := VariantWasabi
	err := SaveTheme(variant)
	if err != nil {
		t.Fatalf("SaveTheme failed: %v", err)
	}

	// Verify the directory was created
	configDir := filepath.Join(tempDir, ".bento")
	info, err := os.Stat(configDir)
	if err != nil {
		t.Fatalf("Config directory was not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("Config path is not a directory")
	}
}

func TestSaveTheme_AllVariants(t *testing.T) {
	// Test that all variants can be saved and loaded
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")

	t.Setenv("HOME", tempDir)
	defer func() {
		os.Setenv("HOME", originalHome)
	}()

	variants := AllVariants()
	for _, variant := range variants {
		err := SaveTheme(variant)
		if err != nil {
			t.Errorf("SaveTheme failed for %s: %v", variant, err)
			continue
		}

		loaded := LoadSavedTheme()
		if loaded != variant {
			t.Errorf("Expected loaded theme to be %s, got %s", variant, loaded)
		}
	}
}
