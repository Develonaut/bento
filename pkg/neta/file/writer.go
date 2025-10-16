// Package file provides file operations nodes.
package file

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"bento/pkg/neta"
)

// Writer writes content to a file.
type Writer struct{}

// NewWriter creates a new file writer node.
func NewWriter() *Writer {
	return &Writer{}
}

// Execute writes content to a file.
func (w *Writer) Execute(ctx context.Context, params map[string]interface{}) (neta.Result, error) {
	path := neta.GetStringParam(params, "path", "")
	if path == "" {
		return neta.Result{}, fmt.Errorf("path parameter required")
	}

	// Get content from params or use input from previous node
	content := neta.GetStringParam(params, "content", "")
	if content == "" {
		// Try to use input from previous node
		if input, ok := params["input"]; ok {
			if strInput, ok := input.(string); ok {
				content = strInput
			} else {
				// Convert to string if not already
				content = fmt.Sprintf("%v", input)
			}
		}
	}

	if content == "" {
		return neta.Result{}, fmt.Errorf("content parameter or input required")
	}

	// Expand home directory if needed
	if filepath.IsAbs(path) && len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return neta.Result{}, fmt.Errorf("failed to get home directory: %w", err)
		}
		path = filepath.Join(home, path[1:])
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return neta.Result{}, fmt.Errorf("failed to create directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return neta.Result{}, fmt.Errorf("failed to write file: %w", err)
	}

	result := map[string]interface{}{
		"path":    path,
		"bytes":   len(content),
		"message": fmt.Sprintf("Successfully wrote %d bytes to %s", len(content), path),
	}

	return neta.Result{Output: result}, nil
}
