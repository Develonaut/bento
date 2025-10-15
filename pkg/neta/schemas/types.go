// Package schemas provides validation schemas and types for node validation.
package schemas

import "fmt"

// Schema defines the validation rules for a node type.
// Each node type (http, transform, etc.) implements this interface
// to specify its required parameters and field metadata.
type Schema interface {
	// Validate checks if parameters are valid
	Validate(params map[string]interface{}) error

	// Fields returns field definitions for UI forms
	Fields() []Field
}

// Field describes a parameter field for a node type.
// Fields power the Huh form generation in the editor UI.
type Field struct {
	Name        string
	Type        FieldType
	Required    bool
	Description string
	Default     interface{}
	Enum        []string // For string enums like HTTP methods
	Min         *int     // For int/duration fields
	Max         *int     // For int/duration fields
}

// FieldType represents the type of a field.
type FieldType string

const (
	FieldString   FieldType = "string"
	FieldInt      FieldType = "int"
	FieldBool     FieldType = "bool"
	FieldDuration FieldType = "duration"
	FieldMap      FieldType = "map"
	FieldArray    FieldType = "array"
)

// ValidationError represents a single validation failure.
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors represents multiple validation failures.
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return "no errors"
	}

	msg := "validation failed:\n"
	for _, err := range e {
		msg += fmt.Sprintf("  - %s\n", err.Error())
	}
	return msg
}
