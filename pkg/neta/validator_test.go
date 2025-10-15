package neta

import (
	"strings"
	"testing"

	"bento/pkg/neta/schemas"
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

func TestValidator_JQNode(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name    string
		def     Definition
		wantErr bool
	}{
		{
			name: "valid jq node",
			def: Definition{
				Version: "1.0",
				Type:    "jq",
				Name:    "Transform",
				Parameters: map[string]interface{}{
					"query": ".data",
				},
			},
			wantErr: false,
		},
		{
			name: "missing query",
			def: Definition{
				Version: "1.0",
				Type:    "jq",
				Name:    "Transform",
				Parameters: map[string]interface{}{},
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
	var urlField *schemas.Field
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

func TestValidator_ListTypes(t *testing.T) {
	validator := NewValidator()

	types := validator.ListTypes()
	if len(types) == 0 {
		t.Fatal("expected at least one registered type")
	}

	// Check types are sorted
	for i := 1; i < len(types); i++ {
		if types[i-1] >= types[i] {
			t.Errorf("types not sorted: %v", types)
			break
		}
	}

	// Verify expected types are present
	expectedTypes := []string{"http", "jq", "sequence", "parallel", "for", "if"}
	for _, expected := range expectedTypes {
		found := false
		for _, actual := range types {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected type %q not found in %v", expected, types)
		}
	}
}

func TestValidator_ValidateRecursive(t *testing.T) {
	validator := NewValidator()

	def := Definition{
		Version: "1.0",
		Type:    "sequence",
		Name:    "Root",
		Nodes: []Definition{
			{
				Version: "1.0",
				Type:    "http",
				Name:    "Child1",
				Parameters: map[string]interface{}{
					"url": "https://example.com",
				},
			},
			{
				Version: "1.0",
				Type:    "http",
				Name:    "Child2",
				Parameters: map[string]interface{}{
					// Missing url - should fail
					"method": "GET",
				},
			},
		},
	}

	err := validator.ValidateRecursive(def)
	if err == nil {
		t.Fatal("expected validation error for child node")
	}

	if !strings.Contains(err.Error(), "url") {
		t.Errorf("expected error about url, got: %v", err)
	}
}

func TestValidator_UnknownNodeType(t *testing.T) {
	validator := NewValidator()

	def := Definition{
		Version: "1.0",
		Type:    "unknown-type",
		Name:    "Test",
	}

	err := validator.Validate(def)
	if err == nil {
		t.Fatal("expected error for unknown node type")
	}

	if !strings.Contains(err.Error(), "unknown node type") {
		t.Errorf("expected 'unknown node type' error, got: %v", err)
	}
}
