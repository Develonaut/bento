package shell

import (
	"context"
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestCommand_Execute(t *testing.T) {
	tests := []struct {
		name    string
		params  map[string]interface{}
		wantErr bool
		verify  func(t *testing.T, output interface{})
	}{
		{
			name: "simple echo command",
			params: map[string]interface{}{
				"command": "echo",
				"args":    []interface{}{"Hello, World!"},
			},
			wantErr: false,
			verify: func(t *testing.T, output interface{}) {
				result, ok := output.(map[string]interface{})
				if !ok {
					t.Fatal("Output is not a map")
				}
				stdout, ok := result["stdout"].(string)
				if !ok {
					t.Fatal("stdout is not a string")
				}
				if !strings.Contains(stdout, "Hello, World!") {
					t.Errorf("stdout = %v, want to contain 'Hello, World!'", stdout)
				}
				exitCode, ok := result["exit_code"].(int)
				if !ok {
					t.Fatal("exit_code is not an int")
				}
				if exitCode != 0 {
					t.Errorf("exit_code = %d, want 0", exitCode)
				}
			},
		},
		{
			name: "command without args",
			params: map[string]interface{}{
				"command": "pwd",
			},
			wantErr: false,
			verify: func(t *testing.T, output interface{}) {
				result, ok := output.(map[string]interface{})
				if !ok {
					t.Fatal("Output is not a map")
				}
				stdout, ok := result["stdout"].(string)
				if !ok {
					t.Fatal("stdout is not a string")
				}
				if stdout == "" {
					t.Error("stdout should not be empty for pwd command")
				}
			},
		},
		{
			name: "command with working directory",
			params: map[string]interface{}{
				"command":     "pwd",
				"working_dir": os.TempDir(),
			},
			wantErr: false,
			verify: func(t *testing.T, output interface{}) {
				result, ok := output.(map[string]interface{})
				if !ok {
					t.Fatal("Output is not a map")
				}
				stdout, ok := result["stdout"].(string)
				if !ok {
					t.Fatal("stdout is not a string")
				}
				// Just verify we got some output (pwd returns something)
				if stdout == "" {
					t.Error("stdout should not be empty for pwd command")
				}
			},
		},
		{
			name: "list files in directory",
			params: map[string]interface{}{
				"command": "ls",
				"args":    []interface{}{"-la"},
			},
			wantErr: runtime.GOOS == "windows", // ls doesn't exist on Windows
			verify: func(t *testing.T, output interface{}) {
				if runtime.GOOS == "windows" {
					return // Skip verification on Windows
				}
				result, ok := output.(map[string]interface{})
				if !ok {
					t.Fatal("Output is not a map")
				}
				stdout, ok := result["stdout"].(string)
				if !ok {
					t.Fatal("stdout is not a string")
				}
				if stdout == "" {
					t.Error("stdout should not be empty for ls command")
				}
			},
		},
		{
			name: "missing command parameter",
			params: map[string]interface{}{
				"args": []interface{}{"test"},
			},
			wantErr: true,
		},
		{
			name: "command not found",
			params: map[string]interface{}{
				"command": "nonexistent_command_12345",
			},
			wantErr: true,
		},
		{
			name: "command fails with non-zero exit",
			params: map[string]interface{}{
				"command": "ls",
				"args":    []interface{}{"/nonexistent/directory/path"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewCommand()
			result, err := cmd.Execute(context.Background(), tt.params)

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

			// Verify output
			if tt.verify != nil {
				tt.verify(t, result.Output)
			}
		})
	}
}

func TestCommand_ExecuteWithContext(t *testing.T) {
	// Test that context cancellation works
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	cmd := NewCommand()
	params := map[string]interface{}{
		"command": "sleep",
		"args":    []interface{}{"10"},
	}

	_, err := cmd.Execute(ctx, params)
	if err == nil {
		t.Error("Expected error when context is cancelled")
	}
}
