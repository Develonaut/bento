// glob.go provides file deletion operations with support for glob patterns.
package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// delete deletes a file or files matching a glob pattern.
// Supports both single file paths and glob patterns (e.g., "*.png", "render-*.png").
func (f *FileSystemNeta) delete(params map[string]interface{}) (interface{}, error) {
	path, ok := params["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path parameter is required and must be a string")
	}

	// Check if path contains glob pattern (* or ?)
	if strings.ContainsAny(path, "*?") {
		return f.deleteGlobPattern(path)
	}

	// Single file deletion
	return f.deleteSingleFile(path)
}

// deleteGlobPattern deletes all files matching a glob pattern.
func (f *FileSystemNeta) deleteGlobPattern(pattern string) (interface{}, error) {
	// Expand glob pattern
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to expand glob pattern: %w", err)
	}

	// Delete all matching files
	deletedCount := 0
	for _, match := range matches {
		if err := os.Remove(match); err != nil {
			return nil, fmt.Errorf("failed to delete file %s: %w", match, err)
		}
		deletedCount++
	}

	return map[string]interface{}{
		"path":    pattern,
		"deleted": deletedCount,
		"files":   matches,
	}, nil
}

// deleteSingleFile deletes a single file.
func (f *FileSystemNeta) deleteSingleFile(path string) (interface{}, error) {
	err := os.Remove(path)
	if err != nil {
		return nil, fmt.Errorf("failed to delete file: %w", err)
	}

	return map[string]interface{}{
		"path":    path,
		"deleted": true,
	}, nil
}
