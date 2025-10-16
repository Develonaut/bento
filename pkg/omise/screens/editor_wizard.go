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
func NewNodeWizard(nodeType string, schema schemas.Schema, values map[string]interface{}) *NodeWizard {
	form := buildFormFromSchema(nodeType, schema, values)

	return &NodeWizard{
		nodeType: nodeType,
		schema:   schema,
		form:     form,
		values:   values,
	}
}

// Form returns the Huh form for Bubble Tea integration
func (w *NodeWizard) Form() *huh.Form {
	return w.form
}

// Values returns the configured parameters
func (w *NodeWizard) Values() map[string]interface{} {
	return w.values
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

// createBentoNameForm creates a form for entering bento name
func createBentoNameForm(nameTarget *string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Bento Name").
				Placeholder("my-awesome-bento").
				Description("Enter a name for your bento").
				Value(nameTarget).
				Validate(validateBentoName),
		),
	)
}

// createNodeTypeForm creates a form for selecting node type
func createNodeTypeForm(nodeTypes []string, typeTarget *string) *huh.Form {
	options := make([]huh.Option[string], len(nodeTypes))
	for i, nodeType := range nodeTypes {
		options[i] = huh.NewOption(nodeType, nodeType)
	}

	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Node Type").
				Description("Select the type of node to add").
				Options(options...).
				Value(typeTarget),
		),
	)
}

// validateBentoName validates bento name format
func validateBentoName(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return fmt.Errorf("bento name is required")
	}
	if len(s) < 3 {
		return fmt.Errorf("bento name must be at least 3 characters")
	}
	// Check for invalid characters (only allow alphanumeric, dash, underscore)
	for _, ch := range s {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') || ch == '-' || ch == '_') {
			return fmt.Errorf("bento name can only contain letters, numbers, dash, and underscore")
		}
	}
	return nil
}
