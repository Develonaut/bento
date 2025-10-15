package jubako

import (
	"os"
	"path/filepath"
	"testing"

	"bento/pkg/neta"
)

func TestParser_ParseBytes(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		wantErr bool
	}{
		{
			name: "valid http node",
			yaml: `type: http
name: Test
parameters:
  url: https://example.com
  method: GET`,
			wantErr: false,
		},
		{
			name: "valid group node",
			yaml: `type: group.sequence
name: Test Group
nodes:
  - type: http
    name: Step 1
    parameters:
      url: https://example.com
  - type: http
    name: Step 2
    parameters:
      url: https://example.com/api`,
			wantErr: false,
		},
		{
			name: "missing type",
			yaml: `name: Test
parameters:
  url: https://example.com`,
			wantErr: true,
		},
		{
			name:    "invalid yaml",
			yaml:    `this: is: not: valid: yaml:`,
			wantErr: true,
		},
		{
			name: "group with invalid child",
			yaml: `type: group.sequence
name: Test Group
nodes:
  - name: Missing Type
    parameters:
      url: https://example.com`,
			wantErr: true,
		},
	}

	parser := NewParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			def, err := parser.ParseBytes([]byte(tt.yaml))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && def.Type == "" {
				t.Error("ParseBytes() returned empty definition")
			}
		})
	}
}

func TestParser_Parse(t *testing.T) {
	// Create temp dir for test files
	tmpDir := t.TempDir()

	// Valid workflow file
	validYAML := `type: http
name: Test Workflow
parameters:
  url: https://example.com
  method: GET`

	validPath := filepath.Join(tmpDir, "valid.bento.yaml")
	if err := os.WriteFile(validPath, []byte(validYAML), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	parser := NewParser()

	t.Run("valid file", func(t *testing.T) {
		def, err := parser.Parse(validPath)
		if err != nil {
			t.Errorf("Parse() error = %v", err)
			return
		}
		if def.Type != "http" {
			t.Errorf("Parse() got type = %v, want http", def.Type)
		}
		if def.Name != "Test Workflow" {
			t.Errorf("Parse() got name = %v, want Test Workflow", def.Name)
		}
	})

	t.Run("non-existent file", func(t *testing.T) {
		_, err := parser.Parse(filepath.Join(tmpDir, "nonexistent.yaml"))
		if err == nil {
			t.Error("Parse() expected error for non-existent file")
		}
	})
}

func TestParser_Format(t *testing.T) {
	parser := NewParser()

	def := neta.Definition{
		Type: "http",
		Name: "Test",
		Parameters: map[string]interface{}{
			"url":    "https://example.com",
			"method": "GET",
		},
	}

	t.Run("format valid definition", func(t *testing.T) {
		data, err := parser.Format(def)
		if err != nil {
			t.Errorf("Format() error = %v", err)
			return
		}
		if len(data) == 0 {
			t.Error("Format() returned empty data")
		}

		// Verify we can parse it back
		parsed, err := parser.ParseBytes(data)
		if err != nil {
			t.Errorf("Format() produced invalid YAML: %v", err)
		}
		if parsed.Type != def.Type {
			t.Errorf("Format/Parse roundtrip failed: got type %v, want %v", parsed.Type, def.Type)
		}
	})
}

func TestValidateDefinition(t *testing.T) {
	tests := []struct {
		name    string
		def     neta.Definition
		wantErr bool
	}{
		{
			name: "valid single node",
			def: neta.Definition{
				Type: "http",
				Name: "Test",
			},
			wantErr: false,
		},
		{
			name: "valid group node",
			def: neta.Definition{
				Type: "group.sequence",
				Name: "Test Group",
				Nodes: []neta.Definition{
					{Type: "http", Name: "Step 1"},
					{Type: "http", Name: "Step 2"},
				},
			},
			wantErr: false,
		},
		{
			name: "missing type",
			def: neta.Definition{
				Name: "Test",
			},
			wantErr: true,
		},
		{
			name: "group with invalid child",
			def: neta.Definition{
				Type: "group.sequence",
				Name: "Test Group",
				Nodes: []neta.Definition{
					{Name: "Missing Type"},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDefinition(tt.def)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateDefinition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
