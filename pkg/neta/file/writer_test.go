package file

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestWriter_Execute(t *testing.T) {
	tests := []struct {
		name    string
		params  map[string]interface{}
		wantErr bool
	}{
		{
			name: "successful write",
			params: map[string]interface{}{
				"path":    filepath.Join(t.TempDir(), "test.txt"),
				"content": "Hello, World!",
			},
			wantErr: false,
		},
		{
			name: "missing path parameter",
			params: map[string]interface{}{
				"content": "Hello, World!",
			},
			wantErr: true,
		},
		{
			name: "missing content parameter",
			params: map[string]interface{}{
				"path": filepath.Join(t.TempDir(), "test.txt"),
			},
			wantErr: true,
		},
		{
			name: "creates parent directories",
			params: map[string]interface{}{
				"path":    filepath.Join(t.TempDir(), "subdir", "test.txt"),
				"content": "Hello, World!",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := NewWriter()
			result, err := writer.Execute(context.Background(), tt.params)

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
			}

			// Verify file was written
			if path, ok := tt.params["path"].(string); ok {
				content, err := os.ReadFile(path)
				if err != nil {
					t.Errorf("Failed to read written file: %v", err)
					return
				}

				expectedContent := tt.params["content"].(string)
				if string(content) != expectedContent {
					t.Errorf("File content = %s, want %s", string(content), expectedContent)
				}
			}
		})
	}
}
