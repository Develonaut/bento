// Package template provides template processing operations for text-based files.
//
// The template neta supports:
//   - XML/SVG element replacement by ID attribute
//   - Simple placeholder replacement ({{key}})
//   - Pure Go implementation (no external dependencies)
//
// This is useful for:
// - Personalizing SVG graphics exported from Figma
// - Generating configuration files from templates
// - Batch processing text/XML files with dynamic content
//
// Example SVG personalization:
//
//	params := map[string]interface{}{
//	    "operation": "replace",
//	    "input": "template.svg",
//	    "output": "personalized.svg",
//	    "replacements": map[string]interface{}{
//	        "product_name": "Combat Dog",
//	        "product_scale": "32mm - Heroic",
//	    },
//	    "mode": "id",
//	}
//	result, err := templateNeta.Execute(ctx, params)
package template

import (
	"context"
	"fmt"
)

// Template implements the template neta for text/XML template processing.
type Template struct{}

// New creates a new template neta instance.
func New() *Template {
	return &Template{}
}

// Execute runs template processing operations.
//
// Parameters:
//   - operation: "replace" (currently only supported operation)
//   - input: input file path
//   - output: output file path
//   - replacements: map[string]interface{} of ID/placeholder to replacement value
//   - mode: "id" (default), "xpath", or "placeholder"
//
// Returns:
//   - path: output file path
//   - replacements_made: number of replacements performed
func (t *Template) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	operation, ok := params["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("operation parameter is required")
	}

	switch operation {
	case "replace":
		return t.replace(ctx, params)
	default:
		return nil, fmt.Errorf("invalid operation: %s", operation)
	}
}
