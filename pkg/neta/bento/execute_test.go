package bento

import (
	"context"
	"testing"

	"bento/pkg/itamae"
	"bento/pkg/jubako"
	"bento/pkg/neta"
	"bento/pkg/neta/transform"
	"bento/pkg/pantry"
)

func TestExecute_Execute(t *testing.T) {
	// Setup: Create a temporary workspace with test bentos
	tmpDir := t.TempDir()
	store, err := jubako.NewStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// Create a simple test bento that transforms input
	testBentoDef := neta.Definition{
		Version: "1.0",
		Type:    "transform.jq",
		Name:    "Simple Transform",
		Parameters: map[string]interface{}{
			"query": ".message",
			"input": map[string]interface{}{"message": "Hello from sub-bento!"},
		},
	}

	if err := store.Save("test-bento", testBentoDef); err != nil {
		t.Fatalf("Failed to save test bento: %v", err)
	}

	// Setup registry and chef for execution
	registry := pantry.New()
	_ = registry.Register("transform.jq", transform.NewJQ())
	chef := itamae.New(registry)

	tests := []struct {
		name    string
		params  map[string]interface{}
		wantErr bool
		verify  func(t *testing.T, output interface{})
	}{
		{
			name: "execute sub-bento successfully",
			params: map[string]interface{}{
				"bento": "test-bento",
			},
			wantErr: false,
			verify: func(t *testing.T, output interface{}) {
				if output == nil {
					t.Fatal("Output should not be nil")
				}
				// The jq transform should extract "Hello from sub-bento!"
				if str, ok := output.(string); ok {
					if str != "Hello from sub-bento!" {
						t.Errorf("Output = %v, want 'Hello from sub-bento!'", str)
					}
				} else {
					t.Errorf("Output type = %T, want string", output)
				}
			},
		},
		{
			name: "missing bento parameter",
			params: map[string]interface{}{
				"inputs": map[string]interface{}{"test": "value"},
			},
			wantErr: true,
		},
		{
			name: "bento not found",
			params: map[string]interface{}{
				"bento": "nonexistent-bento",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := NewExecute(store, chef)
			result, err := exec.Execute(context.Background(), tt.params)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Execute() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Execute() unexpected error: %v", err)
				return
			}

			// Verify output
			if tt.verify != nil {
				tt.verify(t, result.Output)
			}
		})
	}
}

func TestExecute_WithParameterPassing(t *testing.T) {
	// Setup workspace
	tmpDir := t.TempDir()
	store, err := jubako.NewStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// Create a bento that uses parameters from inputs
	// This bento expects an "input" parameter to be passed
	paramBentoDef := neta.Definition{
		Version: "1.0",
		Type:    "transform.jq",
		Name:    "Parameter Transform",
		Parameters: map[string]interface{}{
			"query": ".value * 2",
			// Input will be provided by caller
		},
	}

	if err := store.Save("param-bento", paramBentoDef); err != nil {
		t.Fatalf("Failed to save param bento: %v", err)
	}

	// Setup registry and chef
	registry := pantry.New()
	_ = registry.Register("transform.jq", transform.NewJQ())
	chef := itamae.New(registry)

	// Execute with inputs
	exec := NewExecute(store, chef)
	params := map[string]interface{}{
		"bento": "param-bento",
		"inputs": map[string]interface{}{
			"input": map[string]interface{}{"value": 21},
		},
	}

	result, err := exec.Execute(context.Background(), params)
	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	// Result should be 42 (21 * 2)
	if result.Output == nil {
		t.Fatal("Output should not be nil")
	}

	// JQ can return numbers as int or float64
	switch v := result.Output.(type) {
	case int:
		if v != 42 {
			t.Errorf("Output = %v, want 42", v)
		}
	case float64:
		if v != 42.0 {
			t.Errorf("Output = %v, want 42", v)
		}
	default:
		t.Errorf("Output type = %T, want int or float64", result.Output)
	}
}
