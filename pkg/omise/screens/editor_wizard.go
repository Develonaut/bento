package screens

import (
	"fmt"
	"strings"

	"bento/pkg/neta/schemas"

	"github.com/charmbracelet/huh"
)

// NodeWizard guides parameter configuration using Huh forms
type NodeWizard struct {
	nodeType string
	schema   schemas.Schema
	form     *huh.Form
	values   map[string]interface{}
}

// NewNodeWizard creates wizard for node type
func NewNodeWizard(nodeType string, schema schemas.Schema) *NodeWizard {
	values := make(map[string]interface{})
	form := buildFormFromSchema(nodeType, schema, values)

	return &NodeWizard{
		nodeType: nodeType,
		schema:   schema,
		form:     form,
		values:   values,
	}
}

// Run executes the wizard and returns configured parameters
func (w *NodeWizard) Run() (map[string]interface{}, error) {
	if err := w.form.Run(); err != nil {
		return nil, err
	}
	return w.values, nil
}

// buildFormFromSchema creates a Huh form from schema fields
func buildFormFromSchema(nodeType string, schema schemas.Schema, values map[string]interface{}) *huh.Form {
	fields := schema.Fields()
	inputs := make([]huh.Field, 0, len(fields)+1)

	// Add node name field first
	var nodeName string
	inputs = append(inputs, createNameInput(&nodeName, nodeType))

	// Add fields from schema
	for _, field := range fields {
		input := createFieldInput(field, values)
		if input != nil {
			inputs = append(inputs, input)
		}
	}

	return huh.NewForm(huh.NewGroup(inputs...))
}

// createNameInput creates the node name input field
func createNameInput(target *string, nodeType string) huh.Field {
	return huh.NewInput().
		Title("Node Name").
		Placeholder(fmt.Sprintf("My %s Node", nodeType)).
		Value(target).
		Validate(validateRequired)
}

// createFieldInput creates an input field based on field type
func createFieldInput(field schemas.Field, values map[string]interface{}) huh.Field {
	switch field.Type {
	case schemas.FieldString:
		return createStringInput(field, values)
	case schemas.FieldBool:
		return createBoolInput(field, values)
	case schemas.FieldInt:
		return createIntInput(field, values)
	default:
		return nil // Skip unsupported types for now
	}
}

// createStringInput creates a string input field
func createStringInput(field schemas.Field, values map[string]interface{}) huh.Field {
	var value string
	if field.Default != nil {
		if defStr, ok := field.Default.(string); ok {
			value = defStr
		}
	}

	input := huh.NewInput().
		Title(field.Name).
		Description(field.Description).
		Value(&value)

	if len(field.Enum) > 0 {
		return createSelectInput(field, values)
	}

	if field.Required {
		input = input.Validate(validateRequired)
	}

	// Store value in map
	values[field.Name] = &value

	return input
}

// createSelectInput creates a select dropdown for enum fields
func createSelectInput(field schemas.Field, values map[string]interface{}) huh.Field {
	var value string
	if field.Default != nil {
		if defStr, ok := field.Default.(string); ok {
			value = defStr
		}
	}

	options := make([]huh.Option[string], len(field.Enum))
	for i, opt := range field.Enum {
		options[i] = huh.NewOption(opt, opt)
	}

	values[field.Name] = &value

	return huh.NewSelect[string]().
		Title(field.Name).
		Description(field.Description).
		Options(options...).
		Value(&value)
}

// createBoolInput creates a boolean confirm field
func createBoolInput(field schemas.Field, values map[string]interface{}) huh.Field {
	var value bool
	if field.Default != nil {
		if defBool, ok := field.Default.(bool); ok {
			value = defBool
		}
	}

	values[field.Name] = &value

	return huh.NewConfirm().
		Title(field.Name).
		Description(field.Description).
		Value(&value)
}

// createIntInput creates an integer input field
func createIntInput(field schemas.Field, values map[string]interface{}) huh.Field {
	var valueStr string
	if field.Default != nil {
		if defInt, ok := field.Default.(int); ok {
			valueStr = fmt.Sprintf("%d", defInt)
		}
	}

	input := huh.NewInput().
		Title(field.Name).
		Description(field.Description).
		Value(&valueStr)

	if field.Required {
		input = input.Validate(validateRequired)
	}

	values[field.Name] = &valueStr

	return input
}

// Validation functions

// validateRequired ensures field is not empty
func validateRequired(s string) error {
	if strings.TrimSpace(s) == "" {
		return fmt.Errorf("this field is required")
	}
	return nil
}
