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
		json    string
		wantErr bool
	}{
		{
			name: "valid http node",
			json: `{
  "version": "1.0",
  "type": "http",
  "name": "Test",
  "parameters": {
    "url": "https://example.com",
    "method": "GET"
  }
}`,
			wantErr: false,
		},
		{
			name: "valid group node",
			json: `{
  "version": "1.0",
  "type": "sequence",
  "name": "Test Group",
  "nodes": [
    {
      "version": "1.0",
      "type": "http",
      "name": "Step 1",
      "parameters": {
        "url": "https://example.com"
      }
    },
    {
      "version": "1.0",
      "type": "http",
      "name": "Step 2",
      "parameters": {
        "url": "https://example.com/api"
      }
    }
  ]
}`,
			wantErr: false,
		},
		{
			name: "missing version",
			json: `{
  "type": "http",
  "name": "Test",
  "parameters": {
    "url": "https://example.com"
  }
}`,
			wantErr: true,
		},
		{
			name: "incompatible version",
			json: `{
  "version": "2.0",
  "type": "http",
  "name": "Test",
  "parameters": {
    "url": "https://example.com"
  }
}`,
			wantErr: true,
		},
		{
			name: "missing type",
			json: `{
  "version": "1.0",
  "name": "Test",
  "parameters": {
    "url": "https://example.com"
  }
}`,
			wantErr: true,
		},
		{
			name:    "invalid json",
			json:    `{"this": "is" "not" "valid"}`,
			wantErr: true,
		},
		{
			name: "group with invalid child version",
			json: `{
  "version": "1.0",
  "type": "sequence",
  "name": "Test Group",
  "nodes": [
    {
      "type": "http",
      "name": "Missing Version",
      "parameters": {
        "url": "https://example.com"
      }
    }
  ]
}`,
			wantErr: true,
		},
	}

	parser := NewParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			def, err := parser.ParseBytes([]byte(tt.json))
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

	// Valid bento file
	validJSON := `{
  "version": "1.0",
  "type": "http",
  "name": "Test Bento",
  "parameters": {
    "url": "https://example.com",
    "method": "GET"
  }
}`

	validPath := filepath.Join(tmpDir, "valid.bento.json")
	if err := os.WriteFile(validPath, []byte(validJSON), 0644); err != nil {
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
		if def.Name != "Test Bento" {
			t.Errorf("Parse() got name = %v, want Test Bento", def.Name)
		}
	})

	t.Run("non-existent file", func(t *testing.T) {
		_, err := parser.Parse(filepath.Join(tmpDir, "nonexistent.json"))
		if err == nil {
			t.Error("Parse() expected error for non-existent file")
		}
	})
}

func TestParser_Format(t *testing.T) {
	parser := NewParser()

	def := neta.Definition{
		Version: "1.0",
		Type:    "http",
		Name:    "Test",
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
			t.Errorf("Format() produced invalid JSON: %v", err)
		}
		if parsed.Type != def.Type {
			t.Errorf("Format/Parse roundtrip failed: got type %v, want %v", parsed.Type, def.Type)
		}
	})
}

func TestParser_ParseIconAndDescription(t *testing.T) {
	parser := NewParser()

	t.Run("with icon and description", func(t *testing.T) {
		json := `{
  "version": "1.0",
  "type": "http",
  "name": "Test API",
  "icon": "🌐",
  "description": "Makes an API call to fetch data",
  "parameters": {
    "url": "https://example.com",
    "method": "GET"
  }
}`

		def, err := parser.ParseBytes([]byte(json))
		if err != nil {
			t.Fatalf("ParseBytes() error = %v", err)
		}

		if def.Icon != "🌐" {
			t.Errorf("ParseBytes() got icon = %q, want %q", def.Icon, "🌐")
		}

		if def.Description != "Makes an API call to fetch data" {
			t.Errorf("ParseBytes() got description = %q, want %q", def.Description, "Makes an API call to fetch data")
		}
	})

	t.Run("without icon and description", func(t *testing.T) {
		json := `{
  "version": "1.0",
  "type": "http",
  "name": "Test API",
  "parameters": {
    "url": "https://example.com",
    "method": "GET"
  }
}`

		def, err := parser.ParseBytes([]byte(json))
		if err != nil {
			t.Fatalf("ParseBytes() error = %v", err)
		}

		if def.Icon != "" {
			t.Errorf("ParseBytes() got icon = %q, want empty string", def.Icon)
		}

		if def.Description != "" {
			t.Errorf("ParseBytes() got description = %q, want empty string", def.Description)
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
				Version: "1.0",
				Type:    "http",
				Name:    "Test",
				Parameters: map[string]interface{}{
					"url": "https://example.com",
				},
			},
			wantErr: false,
		},
		{
			name: "valid group node",
			def: neta.Definition{
				Version: "1.0",
				Type:    "sequence",
				Name:    "Test Group",
				Nodes: []neta.Definition{
					{Version: "1.0", Type: "http", Name: "Step 1", Parameters: map[string]interface{}{"url": "https://example.com"}},
					{Version: "1.0", Type: "http", Name: "Step 2", Parameters: map[string]interface{}{"url": "https://example.com"}},
				},
			},
			wantErr: false,
		},
		{
			name: "missing version",
			def: neta.Definition{
				Type: "http",
				Name: "Test",
			},
			wantErr: true,
		},
		{
			name: "incompatible version",
			def: neta.Definition{
				Version: "2.0",
				Type:    "http",
				Name:    "Test",
			},
			wantErr: true,
		},
		{
			name: "missing type",
			def: neta.Definition{
				Version: "1.0",
				Name:    "Test",
			},
			wantErr: true,
		},
		{
			name: "group with invalid child",
			def: neta.Definition{
				Version: "1.0",
				Type:    "sequence",
				Name:    "Test Group",
				Nodes: []neta.Definition{
					{Version: "1.0", Name: "Missing Type"},
				},
			},
			wantErr: true,
		},
		{
			name: "group with child missing version",
			def: neta.Definition{
				Version: "1.0",
				Type:    "sequence",
				Name:    "Test Group",
				Nodes: []neta.Definition{
					{Type: "http", Name: "Missing Version", Parameters: map[string]interface{}{"url": "https://example.com"}},
				},
			},
			wantErr: true,
		},
		{
			name: "description too long",
			def: neta.Definition{
				Version:     "1.0",
				Type:        "http",
				Name:        "Test",
				Description: string(make([]byte, 201)), // 201 characters
				Parameters: map[string]interface{}{
					"url": "https://example.com",
				},
			},
			wantErr: true,
		},
		{
			name: "icon too long",
			def: neta.Definition{
				Version: "1.0",
				Type:    "http",
				Name:    "Test",
				Icon:    "🌐🌐🌐🌐🌐🌐", // More than 10 characters
				Parameters: map[string]interface{}{
					"url": "https://example.com",
				},
			},
			wantErr: true,
		},
		{
			name: "valid with icon and description",
			def: neta.Definition{
				Version:     "1.0",
				Type:        "http",
				Name:        "Test",
				Icon:        "🌐",
				Description: "A valid description",
				Parameters: map[string]interface{}{
					"url": "https://example.com",
				},
			},
			wantErr: false,
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
