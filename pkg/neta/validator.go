package neta

import (
	"fmt"
	"sort"

	"bento/pkg/neta/schemas"
)

// Validator validates definitions against registered schemas.
type Validator struct {
	schemas map[string]schemas.Schema
}

// NewValidator creates a validator with all built-in schemas registered.
func NewValidator() *Validator {
	v := &Validator{
		schemas: make(map[string]schemas.Schema),
	}

	// Register built-in schemas with both short and full names
	// Short names (for pantry compatibility)
	v.Register("http", schemas.NewHTTPSchema())
	v.Register("jq", schemas.NewJQSchema())
	v.Register("sequence", schemas.NewSequenceSchema())
	v.Register("parallel", schemas.NewParallelSchema())
	v.Register("for", schemas.NewForLoopSchema())
	v.Register("if", schemas.NewIfSchema())

	// Full names (for YAML compatibility)
	v.Register("transform.jq", schemas.NewJQSchema())
	v.Register("group.sequence", schemas.NewSequenceSchema())
	v.Register("group.parallel", schemas.NewParallelSchema())
	v.Register("loop.for", schemas.NewForLoopSchema())
	v.Register("conditional.if", schemas.NewIfSchema())

	return v
}

// Register adds a schema for a node type.
func (v *Validator) Register(nodeType string, schema schemas.Schema) {
	v.schemas[nodeType] = schema
}

// Validate checks if a definition's parameters are valid.
func (v *Validator) Validate(def Definition) error {
	schema, ok := v.schemas[def.Type]
	if !ok {
		return fmt.Errorf("unknown node type: %s", def.Type)
	}

	return schema.Validate(def.Parameters)
}

// ValidateRecursive validates a definition and all its children.
func (v *Validator) ValidateRecursive(def Definition) error {
	if err := v.Validate(def); err != nil {
		return formatNodeError(def.Name, err)
	}

	for i, child := range def.Nodes {
		if err := v.ValidateRecursive(child); err != nil {
			return fmt.Errorf("node %d: %w", i, err)
		}
	}

	return nil
}

// formatNodeError formats validation error with node name.
func formatNodeError(nodeName string, err error) error {
	if nodeName == "" {
		return err
	}
	return fmt.Errorf("node %q: %w", nodeName, err)
}

// GetSchema returns the schema for a node type.
func (v *Validator) GetSchema(nodeType string) (schemas.Schema, bool) {
	schema, ok := v.schemas[nodeType]
	return schema, ok
}

// ListTypes returns all registered node types sorted alphabetically.
func (v *Validator) ListTypes() []string {
	types := make([]string, 0, len(v.schemas))
	for t := range v.schemas {
		types = append(types, t)
	}
	sort.Strings(types)
	return types
}
