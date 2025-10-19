// Package filesystem provides file system operations for the bento workflow system.
//
// The filesystem neta allows you to perform common file operations:
//   - read: Read file contents
//   - write: Write content to a file
//   - copy: Copy a file from source to destination
//   - move: Move/rename a file
//   - delete: Delete a file or glob pattern
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
//	// Delete with glob pattern
//	params := map[string]interface{}{
//	    "operation": "delete",
//	    "path": "products/*/render-*.png",
//	}
//
// Learn more about Go's os and io packages:
// https://pkg.go.dev/os
// https://pkg.go.dev/io
package filesystem

import (
	"context"
	"fmt"

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
