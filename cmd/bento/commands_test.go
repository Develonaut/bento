package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPrepareCommand(t *testing.T) {
	// Create temp .bento.yaml file
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.bento.yaml")

	content := []byte(`type: http
name: Test
parameters:
  url: https://example.com
  method: GET
`)
	if err := os.WriteFile(file, content, 0644); err != nil {
		t.Fatal(err)
	}

	// Test prepare command
	rootCmd.SetArgs([]string{"prepare", file})
	if err := rootCmd.Execute(); err != nil {
		t.Errorf("prepare failed: %v", err)
	}
}

func TestPrepareCommandInvalidFile(t *testing.T) {
	// Test with non-existent file
	rootCmd.SetArgs([]string{"prepare", "nonexistent.bento.yaml"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}

func TestPantryCommand(t *testing.T) {
	// Test pantry command
	rootCmd.SetArgs([]string{"pantry"})
	if err := rootCmd.Execute(); err != nil {
		t.Errorf("pantry failed: %v", err)
	}
}

func TestPantryCommandWithSearch(t *testing.T) {
	// Test pantry command with search
	rootCmd.SetArgs([]string{"pantry", "http"})
	if err := rootCmd.Execute(); err != nil {
		t.Errorf("pantry with search failed: %v", err)
	}
}

func TestTasteCommand(t *testing.T) {
	// Create temp .bento.yaml file
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.bento.yaml")

	content := []byte(`type: http
name: Test
parameters:
  url: https://example.com
  method: GET
`)
	if err := os.WriteFile(file, content, 0644); err != nil {
		t.Fatal(err)
	}

	// Test taste command
	rootCmd.SetArgs([]string{"taste", file})
	if err := rootCmd.Execute(); err != nil {
		t.Errorf("taste failed: %v", err)
	}
}

func TestPackCommandDryRun(t *testing.T) {
	// Create temp .bento.yaml file
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.bento.yaml")

	content := []byte(`type: http
name: Test
parameters:
  url: https://example.com
  method: GET
`)
	if err := os.WriteFile(file, content, 0644); err != nil {
		t.Fatal(err)
	}

	// Test pack command with dry-run
	rootCmd.SetArgs([]string{"pack", "--dry-run", file})
	if err := rootCmd.Execute(); err != nil {
		t.Errorf("pack dry-run failed: %v", err)
	}
}

func TestValidateGroup(t *testing.T) {
	// Create temp group .bento.yaml file
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "group.bento.yaml")

	content := []byte(`type: sequence
name: Test Group
nodes:
  - type: http
    name: Request 1
    parameters:
      url: https://example.com
  - type: http
    name: Request 2
    parameters:
      url: https://example.com
`)
	if err := os.WriteFile(file, content, 0644); err != nil {
		t.Fatal(err)
	}

	// Test prepare command on group
	rootCmd.SetArgs([]string{"prepare", file})
	if err := rootCmd.Execute(); err != nil {
		t.Errorf("prepare group failed: %v", err)
	}
}

func TestValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name: "missing type",
			content: `name: Test
parameters:
  url: https://example.com`,
			wantErr: true,
		},
		{
			name: "empty group",
			content: `type: sequence
name: Empty Group
nodes: []`,
			wantErr: true,
		},
		{
			name: "valid single node",
			content: `type: http
name: Valid
parameters:
  url: https://example.com`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			file := filepath.Join(tmpDir, "test.bento.yaml")

			if err := os.WriteFile(file, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			rootCmd.SetArgs([]string{"prepare", file})
			err := rootCmd.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("prepare() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
