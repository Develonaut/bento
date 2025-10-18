package file

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestMkdir_Execute(t *testing.T) {
	tests := []struct {
		name    string
		params  map[string]interface{}
		wantErr bool
		verify  func(t *testing.T, path string)
	}{
		{
			name: "create single directory",
			params: map[string]interface{}{
				"path": filepath.Join(t.TempDir(), "testdir"),
			},
			wantErr: false,
			verify: func(t *testing.T, path string) {
				info, err := os.Stat(path)
				if err != nil {
					t.Fatalf("Directory not created: %v", err)
				}
				if !info.IsDir() {
					t.Error("Path is not a directory")
				}
			},
		},
		{
			name: "create nested directories recursively",
			params: map[string]interface{}{
				"path":      filepath.Join(t.TempDir(), "parent", "child", "grandchild"),
				"recursive": true,
			},
			wantErr: false,
			verify: func(t *testing.T, path string) {
				info, err := os.Stat(path)
				if err != nil {
					t.Fatalf("Nested directories not created: %v", err)
				}
				if !info.IsDir() {
					t.Error("Path is not a directory")
				}
			},
		},
		{
			name: "create directory with custom mode",
			params: map[string]interface{}{
				"path": filepath.Join(t.TempDir(), "testdir"),
				"mode": 0700,
			},
			wantErr: false,
			verify: func(t *testing.T, path string) {
				info, err := os.Stat(path)
				if err != nil {
					t.Fatalf("Directory not created: %v", err)
				}
				// Check permissions (mask off type bits)
				if info.Mode().Perm() != 0700 {
					t.Errorf("Directory mode = %o, want 0700", info.Mode().Perm())
				}
			},
		},
		{
			name: "missing path parameter",
			params: map[string]interface{}{
				"recursive": true,
			},
			wantErr: true,
		},
		{
			name: "directory already exists - no error",
			params: map[string]interface{}{
				"path": filepath.Join(t.TempDir(), "existing"),
			},
			wantErr: false,
			verify: func(t *testing.T, path string) {
				// Create it first
				if err := os.MkdirAll(path, 0755); err != nil {
					t.Fatalf("Failed to pre-create directory: %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mkdir := NewMkdir()

			// Pre-verify if needed
			if tt.verify != nil && tt.name == "directory already exists - no error" {
				if path, ok := tt.params["path"].(string); ok {
					tt.verify(t, path)
				}
			}

			result, err := mkdir.Execute(context.Background(), tt.params)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Execute() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Execute() unexpected error: %v", err)
				return
			}

			if result.Output == nil {
				t.Error("Execute() output is nil")
				return
			}

			// Verify the directory was created
			if path, ok := tt.params["path"].(string); ok {
				if tt.verify != nil && tt.name != "directory already exists - no error" {
					tt.verify(t, path)
				}

				// Check result contains path
				output, ok := result.Output.(map[string]interface{})
				if !ok {
					t.Error("Output is not a map")
					return
				}

				if output["path"] != path {
					t.Errorf("Output path = %v, want %v", output["path"], path)
				}

				if created, ok := output["created"].(bool); !ok || !created {
					t.Error("Output should indicate directory was created")
				}
			}
		})
	}
}
