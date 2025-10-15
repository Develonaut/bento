package schemas

import (
	"strings"
	"testing"
)

func TestHTTPSchema_Validate(t *testing.T) {
	schema := NewHTTPSchema()

	tests := []struct {
		name    string
		params  map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid minimal",
			params: map[string]interface{}{
				"url": "https://example.com",
			},
			wantErr: false,
		},
		{
			name: "valid with all fields",
			params: map[string]interface{}{
				"url":    "https://example.com/api",
				"method": "POST",
				"headers": map[string]interface{}{
					"Content-Type": "application/json",
				},
				"body": `{"key": "value"}`,
			},
			wantErr: false,
		},
		{
			name: "missing url",
			params: map[string]interface{}{
				"method": "GET",
			},
			wantErr: true,
		},
		{
			name: "empty url",
			params: map[string]interface{}{
				"url": "",
			},
			wantErr: true,
		},
		{
			name: "invalid method",
			params: map[string]interface{}{
				"url":    "https://example.com",
				"method": "INVALID",
			},
			wantErr: true,
		},
		{
			name: "method case insensitive",
			params: map[string]interface{}{
				"url":    "https://example.com",
				"method": "get",
			},
			wantErr: false,
		},
		{
			name: "invalid headers not map",
			params: map[string]interface{}{
				"url":     "https://example.com",
				"headers": "not a map",
			},
			wantErr: true,
		},
		{
			name: "invalid headers value not string",
			params: map[string]interface{}{
				"url": "https://example.com",
				"headers": map[string]interface{}{
					"Content-Length": 123,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := schema.Validate(tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHTTPSchema_Fields(t *testing.T) {
	schema := NewHTTPSchema()
	fields := schema.Fields()

	if len(fields) == 0 {
		t.Fatal("expected fields to be defined")
	}

	// Check url field
	var urlField *interface{}
	for i := range fields {
		if fields[i].Name == "url" {
			if !fields[i].Required {
				t.Error("url field should be required")
			}
			urlField = new(interface{})
			break
		}
	}
	if urlField == nil {
		t.Error("expected url field to be defined")
	}

	// Check method field has enum
	for i := range fields {
		if fields[i].Name == "method" {
			if len(fields[i].Enum) == 0 {
				t.Error("method field should have enum values")
			}
			if fields[i].Required {
				t.Error("method field should be optional")
			}
			break
		}
	}
}

func TestHTTPSchema_ErrorMessages(t *testing.T) {
	schema := NewHTTPSchema()

	// Test missing url produces clear error
	err := schema.Validate(map[string]interface{}{})
	if err == nil {
		t.Fatal("expected error for missing url")
	}

	if !strings.Contains(err.Error(), "url") {
		t.Errorf("error should mention 'url', got: %v", err)
	}

	// Test invalid method produces clear error
	err = schema.Validate(map[string]interface{}{
		"url":    "https://example.com",
		"method": "INVALID",
	})
	if err == nil {
		t.Fatal("expected error for invalid method")
	}

	if !strings.Contains(err.Error(), "method") {
		t.Errorf("error should mention 'method', got: %v", err)
	}
}
