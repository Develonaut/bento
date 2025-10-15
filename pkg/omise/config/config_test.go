package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	// Should end with .bento
	if !strings.HasSuffix(cfg.SaveDirectory, ".bento") {
		t.Errorf("Expected SaveDirectory to end with .bento, got %s", cfg.SaveDirectory)
	}

	// Should be an absolute path when home dir is available
	if !filepath.IsAbs(cfg.SaveDirectory) {
		t.Errorf("Expected SaveDirectory to be absolute, got %s", cfg.SaveDirectory)
	}
}

func TestExpandHome(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{
			name:     "expands tilde",
			input:    "~/.bento",
			contains: ".bento",
		},
		{
			name:     "expands tilde with subdirectory",
			input:    "~/.bento/bentos",
			contains: ".bento",
		},
		{
			name:     "leaves absolute path unchanged",
			input:    "/usr/local/bento",
			contains: "/usr/local/bento",
		},
		{
			name:     "leaves relative path unchanged",
			input:    "./bento",
			contains: "./bento",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandHome(tt.input)
			if !strings.Contains(result, tt.contains) {
				t.Errorf("expandHome(%s) = %s, should contain %s", tt.input, result, tt.contains)
			}
		})
	}
}

func TestContractHome(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Cannot get home directory")
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "contracts home to tilde",
			input:    filepath.Join(home, ".bento"),
			expected: "~/.bento",
		},
		{
			name:     "contracts home with bento root",
			input:    filepath.Join(home, ".bento"),
			expected: "~/.bento",
		},
		{
			name:     "leaves non-home path unchanged",
			input:    "/usr/local/bento",
			expected: "/usr/local/bento",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contractHome(tt.input)
			if result != tt.expected {
				t.Errorf("contractHome(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetSaveDirectory(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Cannot get home directory")
	}

	cfg := Config{
		SaveDirectory: filepath.Join(home, ".bento"),
	}

	result := cfg.GetSaveDirectory()
	expected := "~/.bento"

	if result != expected {
		t.Errorf("GetSaveDirectory() = %s, want %s", result, expected)
	}
}

func TestSaveAndLoad(t *testing.T) {
	// Create temp directory for test
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Test saving config
	cfg := Config{
		SaveDirectory: filepath.Join(tmpDir, "my-bentos"),
	}

	err := Save(cfg)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify config file exists
	configPath := filepath.Join(tmpDir, ".bento", "config")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("Config file not created at %s", configPath)
	}

	// Test loading config
	loaded := Load()
	if loaded.SaveDirectory != cfg.SaveDirectory {
		t.Errorf("Load() SaveDirectory = %s, want %s", loaded.SaveDirectory, cfg.SaveDirectory)
	}
}

func TestLoadNonExistent(t *testing.T) {
	// Create temp directory for test
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Load should return defaults when no config exists
	cfg := Load()
	expected := filepath.Join(tmpDir, ".bento")

	if cfg.SaveDirectory != expected {
		t.Errorf("Load() on non-existent config = %s, want default %s", cfg.SaveDirectory, expected)
	}
}
