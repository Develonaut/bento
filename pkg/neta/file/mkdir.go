package file

import (
	"context"
	"fmt"
	"os"

	"bento/pkg/neta"
)

// Mkdir creates directories.
type Mkdir struct{}

// NewMkdir creates a new directory creation node.
func NewMkdir() *Mkdir {
	return &Mkdir{}
}

// Execute creates a directory with optional recursive creation.
func (m *Mkdir) Execute(ctx context.Context, params map[string]interface{}) (neta.Result, error) {
	path := neta.GetStringParam(params, "path", "")
	if path == "" {
		return neta.Result{}, fmt.Errorf("path parameter required")
	}

	// Get recursive flag (default: true for user convenience)
	recursive := true
	if val, ok := params["recursive"].(bool); ok {
		recursive = val
	}

	// Get mode (default: 0755)
	mode := os.FileMode(0755)
	if val, ok := params["mode"].(int); ok {
		mode = os.FileMode(val)
	}

	// Create directory
	var err error
	if recursive {
		err = os.MkdirAll(path, mode)
	} else {
		err = os.Mkdir(path, mode)
	}

	if err != nil {
		return neta.Result{}, fmt.Errorf("failed to create directory: %w", err)
	}

	result := map[string]interface{}{
		"path":    path,
		"created": true,
	}

	return neta.Result{Output: result}, nil
}
