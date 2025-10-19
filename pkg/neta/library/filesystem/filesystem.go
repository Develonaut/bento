// Package filesystem provides file system operations for the bento workflow system.
//
// The filesystem neta allows you to perform common file operations:
//   - read: Read file contents
//   - write: Write content to a file
//   - copy: Copy a file from source to destination
//   - move: Move/rename a file
//   - delete: Delete a file
//   - mkdir: Create a directory
//   - exists: Check if a file or directory exists
//
// Example usage:
//
//	// Read a file
//	params := map[string]interface{}{
//	    "operation": "read",
//	    "path": "/path/to/file.txt",
//	}
//
//	// Write a file
//	params := map[string]interface{}{
//	    "operation": "write",
//	    "path": "/path/to/file.txt",
//	    "content": "Hello, world!",
//	}
//
//	// Copy a file
//	params := map[string]interface{}{
//	    "operation": "copy",
//	    "source": "/path/to/source.txt",
//	    "dest": "/path/to/dest.txt",
//	}
//
// Learn more about Go's os and io packages:
// https://pkg.go.dev/os
// https://pkg.go.dev/io
package filesystem

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/Develonaut/bento/pkg/neta"
)

// FileSystemNeta implements file system operations.
type FileSystemNeta struct{}

// New creates a new filesystem neta instance.
func New() neta.Executable {
	return &FileSystemNeta{}
}

// Execute performs a file system operation based on the provided parameters.
//
// Parameters:
//   - operation (string, required): The operation to perform
//     (read, write, copy, move, delete, mkdir, exists)
//   - path (string, required for most operations): The file/directory path
//   - content (string, required for write): Content to write
//   - source (string, required for copy/move): Source path
//   - dest (string, required for copy/move): Destination path
//
// Returns a map with operation-specific results.
func (f *FileSystemNeta) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	// Extract operation
	operation, ok := params["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("operation parameter is required and must be a string")
	}

	// Route to appropriate operation
	switch operation {
	case "read":
		return f.read(params)
	case "write":
		return f.write(params)
	case "copy":
		return f.copy(params)
	case "move":
		return f.move(params)
	case "delete":
		return f.delete(params)
	case "mkdir":
		return f.mkdir(params)
	case "exists":
		return f.exists(params)
	default:
		return nil, fmt.Errorf("unsupported operation: %s", operation)
	}
}

// read reads the contents of a file.
func (f *FileSystemNeta) read(params map[string]interface{}) (interface{}, error) {
	path, ok := params["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path parameter is required and must be a string")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return map[string]interface{}{
		"content": string(data),
		"path":    path,
	}, nil
}

// write writes content to a file.
func (f *FileSystemNeta) write(params map[string]interface{}) (interface{}, error) {
	path, ok := params["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path parameter is required and must be a string")
	}

	content, ok := params["content"].(string)
	if !ok {
		return nil, fmt.Errorf("content parameter is required and must be a string")
	}

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	return map[string]interface{}{
		"path":    path,
		"written": true,
	}, nil
}

// copy copies a file from source to destination.
func (f *FileSystemNeta) copy(params map[string]interface{}) (interface{}, error) {
	source, ok := params["source"].(string)
	if !ok {
		return nil, fmt.Errorf("source parameter is required and must be a string")
	}

	dest, ok := params["dest"].(string)
	if !ok {
		return nil, fmt.Errorf("dest parameter is required and must be a string")
	}

	// Open source file
	srcFile, err := os.Open(source)
	if err != nil {
		return nil, fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Create destination file
	destFile, err := os.Create(dest)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	// Copy contents
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	return map[string]interface{}{
		"source": source,
		"dest":   dest,
		"copied": true,
	}, nil
}

// move moves/renames a file.
func (f *FileSystemNeta) move(params map[string]interface{}) (interface{}, error) {
	source, ok := params["source"].(string)
	if !ok {
		return nil, fmt.Errorf("source parameter is required and must be a string")
	}

	dest, ok := params["dest"].(string)
	if !ok {
		return nil, fmt.Errorf("dest parameter is required and must be a string")
	}

	err := os.Rename(source, dest)
	if err != nil {
		return nil, fmt.Errorf("failed to move file: %w", err)
	}

	return map[string]interface{}{
		"source": source,
		"dest":   dest,
		"moved":  true,
	}, nil
}

// delete deletes a file.
func (f *FileSystemNeta) delete(params map[string]interface{}) (interface{}, error) {
	path, ok := params["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path parameter is required and must be a string")
	}

	err := os.Remove(path)
	if err != nil {
		return nil, fmt.Errorf("failed to delete file: %w", err)
	}

	return map[string]interface{}{
		"path":    path,
		"deleted": true,
	}, nil
}

// mkdir creates a directory.
func (f *FileSystemNeta) mkdir(params map[string]interface{}) (interface{}, error) {
	path, ok := params["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path parameter is required and must be a string")
	}

	err := os.MkdirAll(path, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	return map[string]interface{}{
		"path":    path,
		"created": true,
	}, nil
}

// exists checks if a file or directory exists.
func (f *FileSystemNeta) exists(params map[string]interface{}) (interface{}, error) {
	path, ok := params["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path parameter is required and must be a string")
	}

	_, err := os.Stat(path)
	exists := !os.IsNotExist(err)

	return map[string]interface{}{
		"path":   path,
		"exists": exists,
	}, nil
}
