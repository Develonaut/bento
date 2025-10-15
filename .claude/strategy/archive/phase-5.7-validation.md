# Phase 5.7: Node and Bento Validation Framework

**Status**: Pending
**Duration**: 3-4 hours
**Prerequisites**: Phase 5.5 complete (version validation working)

## Overview

Add structured validation for node types and bento definitions to ensure users create well-formed bentos. This validation framework will power the Huh forms in the editor and provide clear guidance on required/optional parameters.

**Inspiration**: Similar to Zod validation in TypeScript/JavaScript, but idiomatic Go.

## Pre-Work Checklist

Before starting, you MUST:

1. ✅ Read [BENTO_BOX_PRINCIPLE.md](../BENTO_BOX_PRINCIPLE.md)
2. ✅ Confirm: "I understand the Bento Box Principle and will follow it"
3. ✅ Use TodoWrite to track all tasks
4. ✅ Phase 5.5 approved by Karen

## Goals

1. Define validation schemas for all node types
2. Implement validation framework (no external validation libs)
3. Integrate with jubako parser for automatic validation
4. Provide clear, actionable error messages
5. Power Huh forms in editor with field metadata
6. Prevent invalid "compartments" in the bento box

## Why This Matters

**Current State:**
```yaml
version: "1.0"
type: http
# Missing required url parameter - only caught at execution time!
parameters:
  method: GET
```

**With Validation:**
```bash
$ bento prepare bad-http.bento.yaml
Error: validation failed: node "Fetch Data":
  - parameters.url is required
  - parameters.method must be one of: GET, POST, PUT, DELETE, PATCH
```

**In Editor:**
The validation framework will tell Huh forms:
- Which fields are required
- Which fields are optional
- Valid values for enums (method, etc.)
- Field types (string, int, bool)
- Help text for each field

## Architecture

```
pkg/neta/
  ├── definition.go      # Core Definition type (exists)
  ├── version.go         # Version validation (exists)
  ├── validator.go       # NEW - Validation framework
  ├── schema.go          # NEW - Schema definitions
  └── schemas/           # NEW - Per-type schemas
      ├── http.go        # HTTP node schema
      ├── transform.go   # Transform node schema
      ├── group.go       # Group node schema
      ├── loop.go        # Loop node schema
      └── conditional.go # Conditional node schema
```

## Deliverables

### 1. Validation Framework

**File**: `pkg/neta/validator.go` (NEW)
**Target Size**: < 150 lines

```go
package neta

import (
	"fmt"
)

// Validator validates definitions against schemas.
type Validator struct {
	schemas map[string]Schema
}

// NewValidator creates a validator with registered schemas.
func NewValidator() *Validator {
	v := &Validator{
		schemas: make(map[string]Schema),
	}

	// Register built-in schemas
	v.Register("http", NewHTTPSchema())
	v.Register("transform.jq", NewJQSchema())
	v.Register("sequence", NewSequenceSchema())
	v.Register("parallel", NewParallelSchema())
	v.Register("loop.for", NewForLoopSchema())
	v.Register("conditional.if", NewIfSchema())

	return v
}

// Register adds a schema for a node type.
func (v *Validator) Register(nodeType string, schema Schema) {
	v.schemas[nodeType] = schema
}

// Validate checks if a definition is valid.
func (v *Validator) Validate(def Definition) error {
	schema, ok := v.schemas[def.Type]
	if !ok {
		return fmt.Errorf("unknown node type: %s", def.Type)
	}

	return schema.Validate(def.Parameters)
}

// ValidateRecursive validates definition and all children.
func (v *Validator) ValidateRecursive(def Definition) error {
	if err := v.Validate(def); err != nil {
		return fmt.Errorf("node %q: %w", def.Name, err)
	}

	for i, child := range def.Nodes {
		if err := v.ValidateRecursive(child); err != nil {
			return fmt.Errorf("node %d: %w", i, err)
		}
	}

	return nil
}

// GetSchema returns the schema for a node type.
func (v *Validator) GetSchema(nodeType string) (Schema, bool) {
	schema, ok := v.schemas[nodeType]
	return schema, ok
}

// ListTypes returns all registered node types.
func (v *Validator) ListTypes() []string {
	types := make([]string, 0, len(v.schemas))
	for t := range v.schemas {
		types = append(types, t)
	}
	return types
}
```

**Bento Box Compliance:**
- ✅ Single responsibility: Validation orchestration
- ✅ < 150 lines
- ✅ Functions < 20 lines
- ✅ No utils package

### 2. Schema Interface

**File**: `pkg/neta/schema.go` (NEW)
**Target Size**: < 100 lines

```go
package neta

// Schema defines the validation rules for a node type.
type Schema interface {
	// Validate checks if parameters are valid
	Validate(params map[string]interface{}) error

	// Fields returns field definitions for UI forms
	Fields() []Field
}

// Field describes a parameter field.
type Field struct {
	Name        string
	Type        FieldType
	Required    bool
	Description string
	Default     interface{}
	Enum        []string // For string enums
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

// ValidationError represents a validation failure.
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
```

**Bento Box Compliance:**
- ✅ Single responsibility: Schema definition types
- ✅ < 100 lines
- ✅ Clear interfaces

### 3. HTTP Node Schema

**File**: `pkg/neta/schemas/http.go` (NEW)
**Target Size**: < 150 lines

```go
package schemas

import (
	"fmt"
	"strings"

	"bento/pkg/neta"
)

// HTTPSchema validates HTTP node parameters.
type HTTPSchema struct{}

// NewHTTPSchema creates an HTTP schema.
func NewHTTPSchema() *HTTPSchema {
	return &HTTPSchema{}
}

// Validate checks HTTP parameters.
func (s *HTTPSchema) Validate(params map[string]interface{}) error {
	var errs neta.ValidationErrors

	// Required: url
	url, ok := params["url"].(string)
	if !ok || url == "" {
		errs = append(errs, neta.ValidationError{
			Field:   "url",
			Message: "is required and must be a string",
		})
	}

	// Optional: method (must be valid HTTP method)
	if method, ok := params["method"].(string); ok {
		validMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
		if !contains(validMethods, strings.ToUpper(method)) {
			errs = append(errs, neta.ValidationError{
				Field:   "method",
				Message: fmt.Sprintf("must be one of: %s", strings.Join(validMethods, ", ")),
			})
		}
	}

	// Optional: headers (must be map)
	if headers, ok := params["headers"]; ok {
		if _, isMap := headers.(map[string]interface{}); !isMap {
			errs = append(errs, neta.ValidationError{
				Field:   "headers",
				Message: "must be a map of string to string",
			})
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

// Fields returns field definitions for forms.
func (s *HTTPSchema) Fields() []neta.Field {
	return []neta.Field{
		{
			Name:        "url",
			Type:        neta.FieldString,
			Required:    true,
			Description: "HTTP(S) URL to request",
		},
		{
			Name:        "method",
			Type:        neta.FieldString,
			Required:    false,
			Description: "HTTP method",
			Default:     "GET",
			Enum:        []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		},
		{
			Name:        "headers",
			Type:        neta.FieldMap,
			Required:    false,
			Description: "HTTP headers as key-value pairs",
		},
		{
			Name:        "body",
			Type:        neta.FieldString,
			Required:    false,
			Description: "Request body (for POST/PUT/PATCH)",
		},
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
```

**Bento Box Compliance:**
- ✅ Single responsibility: HTTP validation
- ✅ Separate file per node type
- ✅ < 150 lines

### 4. Transform (jq) Schema

**File**: `pkg/neta/schemas/transform.go` (NEW)
**Target Size**: < 100 lines

```go
package schemas

import (
	"bento/pkg/neta"
)

// JQSchema validates jq transform parameters.
type JQSchema struct{}

// NewJQSchema creates a jq schema.
func NewJQSchema() *JQSchema {
	return &JQSchema{}
}

// Validate checks jq parameters.
func (s *JQSchema) Validate(params map[string]interface{}) error {
	var errs neta.ValidationErrors

	// Required: query
	query, ok := params["query"].(string)
	if !ok || query == "" {
		errs = append(errs, neta.ValidationError{
			Field:   "query",
			Message: "is required (jq query string)",
		})
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

// Fields returns field definitions.
func (s *JQSchema) Fields() []neta.Field {
	return []neta.Field{
		{
			Name:        "query",
			Type:        neta.FieldString,
			Required:    true,
			Description: "jq query to transform data (e.g., '.data | .[] | select(.active)')",
		},
		{
			Name:        "input",
			Type:        neta.FieldString,
			Required:    false,
			Description: "Static input data (if not using previous node output)",
		},
	}
}
```

### 5. Group Schema (Sequence/Parallel)

**File**: `pkg/neta/schemas/group.go` (NEW)
**Target Size**: < 100 lines

```go
package schemas

import (
	"bento/pkg/neta"
)

// SequenceSchema validates sequence group parameters.
type SequenceSchema struct{}

// NewSequenceSchema creates a sequence schema.
func NewSequenceSchema() *SequenceSchema {
	return &SequenceSchema{}
}

// Validate checks sequence has child nodes.
func (s *SequenceSchema) Validate(params map[string]interface{}) error {
	// Sequences have no parameters - validation is structural
	return nil
}

// Fields returns empty (sequences defined by child nodes).
func (s *SequenceSchema) Fields() []neta.Field {
	return []neta.Field{
		{
			Name:        "nodes",
			Type:        neta.FieldArray,
			Required:    true,
			Description: "Child nodes to execute in sequence",
		},
	}
}

// ParallelSchema validates parallel group parameters.
type ParallelSchema struct{}

// NewParallelSchema creates a parallel schema.
func NewParallelSchema() *ParallelSchema {
	return &ParallelSchema{}
}

// Validate checks parallel has child nodes.
func (s *ParallelSchema) Validate(params map[string]interface{}) error {
	// Parallel has no parameters - validation is structural
	return nil
}

// Fields returns field definitions.
func (s *ParallelSchema) Fields() []neta.Field {
	return []neta.Field{
		{
			Name:        "nodes",
			Type:        neta.FieldArray,
			Required:    true,
			Description: "Child nodes to execute in parallel",
		},
		{
			Name:        "max_concurrent",
			Type:        neta.FieldInt,
			Required:    false,
			Description: "Maximum concurrent executions (default: unlimited)",
			Default:     0,
		},
	}
}
```

### 6. Integration with Parser

**File**: `pkg/jubako/parser.go` (modify existing)

Update `validateDefinition` to use new validation framework:

```go
// validateDefinition ensures a definition is well-formed.
func validateDefinition(def neta.Definition) error {
	// Validate version first (Phase 5.5)
	if err := neta.ValidateVersion(def.Version); err != nil {
		return fmt.Errorf("version error: %w", err)
	}

	// Validate type exists
	if def.Type == "" {
		return fmt.Errorf("type is required")
	}

	// Validate node parameters (Phase 5.7)
	validator := neta.NewValidator()
	if err := validator.Validate(def); err != nil {
		return fmt.Errorf("parameter validation: %w", err)
	}

	// Validate child nodes recursively
	if def.IsGroup() {
		for i, child := range def.Nodes {
			if err := validateDefinition(child); err != nil {
				return fmt.Errorf("node %d: %w", i, err)
			}
		}
	}

	return nil
}
```

### 7. Comprehensive Tests

**File**: `pkg/neta/validator_test.go` (NEW)
**Target Size**: < 200 lines

```go
package neta

import (
	"testing"
)

func TestValidator_HTTPNode(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name    string
		def     Definition
		wantErr bool
	}{
		{
			name: "valid http node",
			def: Definition{
				Version: "1.0",
				Type:    "http",
				Name:    "Test",
				Parameters: map[string]interface{}{
					"url":    "https://example.com",
					"method": "GET",
				},
			},
			wantErr: false,
		},
		{
			name: "missing url",
			def: Definition{
				Version: "1.0",
				Type:    "http",
				Name:    "Test",
				Parameters: map[string]interface{}{
					"method": "GET",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid method",
			def: Definition{
				Version: "1.0",
				Type:    "http",
				Name:    "Test",
				Parameters: map[string]interface{}{
					"url":    "https://example.com",
					"method": "INVALID",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.def)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidator_GetSchema(t *testing.T) {
	validator := NewValidator()

	schema, ok := validator.GetSchema("http")
	if !ok {
		t.Fatal("expected http schema to be registered")
	}

	fields := schema.Fields()
	if len(fields) == 0 {
		t.Error("expected http schema to have fields")
	}

	// Verify url field exists and is required
	var urlField *Field
	for i := range fields {
		if fields[i].Name == "url" {
			urlField = &fields[i]
			break
		}
	}

	if urlField == nil {
		t.Fatal("expected url field in http schema")
	}

	if !urlField.Required {
		t.Error("url field should be required")
	}
}
```

**File**: `pkg/neta/schemas/http_test.go` (NEW)

Add specific tests for each schema type.

## Integration with Omise Editor

The validation framework powers the editor in two ways:

### 1. Form Generation (Huh Integration)

```go
// In pkg/omise/screens/editor.go

func (m *EditorModel) createNodeForm(nodeType string) *huh.Form {
	validator := neta.NewValidator()
	schema, ok := validator.GetSchema(nodeType)
	if !ok {
		return nil // Unknown node type
	}

	var fields []huh.Field
	for _, field := range schema.Fields() {
		switch field.Type {
		case neta.FieldString:
			if len(field.Enum) > 0 {
				// Enum -> Select
				options := make([]huh.Option[string], len(field.Enum))
				for i, val := range field.Enum {
					options[i] = huh.NewOption(val, val)
				}
				fields = append(fields, huh.NewSelect[string]().
					Title(field.Name).
					Description(field.Description).
					Options(options...))
			} else {
				// String -> Input
				fields = append(fields, huh.NewInput().
					Title(field.Name).
					Description(field.Description).
					Validate(func(s string) error {
						if field.Required && s == "" {
							return fmt.Errorf("%s is required", field.Name)
						}
						return nil
					}))
			}
		case neta.FieldInt:
			// Int -> Input with validation
			fields = append(fields, huh.NewInput().
				Title(field.Name).
				Description(field.Description).
				Validate(validateInt))
		case neta.FieldBool:
			// Bool -> Confirm
			fields = append(fields, huh.NewConfirm().
				Title(field.Name).
				Description(field.Description))
		}
	}

	return huh.NewForm(huh.NewGroup(fields...))
}
```

### 2. Live Validation

```go
// In editor, validate before saving
func (m *EditorModel) validateCurrentBento() error {
	validator := neta.NewValidator()
	return validator.ValidateRecursive(m.currentBento)
}
```

## Error Messages

**Before Phase 5.7:**
```
Error: execution failed: http request failed
```

**After Phase 5.7:**
```
Error: validation failed: node "Fetch User Data":
  - url is required and must be a string
  - method must be one of: GET, POST, PUT, DELETE, PATCH
  - headers must be a map of string to string
```

## Example: Valid Bento with Validation

```yaml
version: "1.0"
type: sequence
name: API Pipeline
nodes:
  - version: "1.0"
    type: http
    name: Fetch Users
    parameters:
      url: https://api.example.com/users
      method: GET
      headers:
        Authorization: Bearer token123

  - version: "1.0"
    type: transform.jq
    name: Extract IDs
    parameters:
      query: ".[] | .id"

  - version: "1.0"
    type: http
    name: Fetch Details
    parameters:
      url: https://api.example.com/user/{{.}}
      method: GET
```

## Validation Commands

```bash
# Validate and see detailed errors
bento prepare workflow.bento.yaml

# List available node types and their schemas
bento pantry list --schemas

# Show schema for specific node type
bento pantry schema http
```

## Success Criteria

Phase 5.7 is complete when:

1. ✅ Validation framework implemented (`validator.go`, `schema.go`)
2. ✅ Schemas for all built-in node types (http, transform, group, loop, conditional)
3. ✅ Parser integration validates all nodes
4. ✅ Clear, actionable error messages
5. ✅ Schema metadata available for editor forms
6. ✅ All tests passing
7. ✅ Files < 150 lines
8. ✅ Functions < 20 lines
9. ✅ **Karen's approval granted**

## Common Pitfalls to Avoid

1. ❌ **External validation libraries** - Keep it simple, standard library only
2. ❌ **Over-abstraction** - Don't create generic "rule engines"
3. ❌ **Poor error messages** - Must tell user exactly what's wrong and how to fix
4. ❌ **Forgetting child nodes** - Validation must be recursive
5. ❌ **No schema metadata** - Must provide Fields() for editor integration

## Bento Box Compliance

- ✅ Each schema in separate file (`schemas/http.go`, etc.)
- ✅ Single responsibility per file
- ✅ No utils package
- ✅ Files < 150 lines
- ✅ Functions < 20 lines
- ✅ Clear package boundaries (neta/validation, neta/schemas)

## Future Extensions (Not in 5.7)

These are explicitly OUT OF SCOPE for Phase 5.7:

- Custom user-defined schemas
- Schema versioning (when node schemas change)
- Complex validation rules (regex, custom validators)
- Schema generation from code
- JSON Schema export

Keep it simple. Ship the basics that work.

## Execution Prompt

```
I'm ready to begin Phase 5.7: Node and Bento Validation Framework.

I have read the Bento Box Principle and will follow it.

Please implement:
- Validation framework (validator.go, schema.go)
- Schemas for all built-in node types
- Parser integration
- Comprehensive tests
- Clear error messages

Each file < 150 lines, functions < 20 lines. I will use TodoWrite to track progress and get Karen's approval before completing.
```

---

**Phase 5.7 Validation Framework**: Prevent invalid compartments, guide users to success! 🎯
