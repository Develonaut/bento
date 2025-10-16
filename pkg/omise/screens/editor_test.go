package screens

import (
	"testing"

	"bento/pkg/jubako"
	"bento/pkg/neta"
	"bento/pkg/pantry"
)

func TestEditor_CreateMode(t *testing.T) {
	workDir := t.TempDir()
	store, err := jubako.NewStore(workDir)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	registry := pantry.New()
	editor := NewEditorCreate(store, registry)

	if editor.mode != EditorModeCreate {
		t.Error("expected create mode")
	}

	if editor.state != StateNaming {
		t.Error("expected naming state")
	}
}

func TestEditor_EditMode(t *testing.T) {
	workDir := t.TempDir()
	store, err := jubako.NewStore(workDir)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	// Create test bento
	def := neta.Definition{
		Version: "1.0",
		Type:    "http",
		Name:    "existing-bento",
		Parameters: map[string]interface{}{
			"url": "https://example.com",
		},
	}
	if err := store.Save("existing-bento", def); err != nil {
		t.Fatalf("failed to save bento: %v", err)
	}

	// Load in editor
	registry := pantry.New()
	editor, err := NewEditorEdit(store, registry, "existing-bento", "")
	if err != nil {
		t.Fatalf("failed to create editor: %v", err)
	}

	if editor.mode != EditorModeEdit {
		t.Error("expected edit mode")
	}

	if editor.def.Name != "existing-bento" {
		t.Error("definition not loaded")
	}
}

func TestEditor_InModalMode(t *testing.T) {
	workDir := t.TempDir()
	store, err := jubako.NewStore(workDir)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	registry := pantry.New()
	editor := NewEditorCreate(store, registry)

	tests := []struct {
		name     string
		state    EditorState
		expected bool
	}{
		{"StateNaming", StateNaming, true},
		{"StateSelectingType", StateSelectingType, true},
		{"StateConfiguringNode", StateConfiguringNode, true},
		{"StateReview", StateReview, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			editor.state = tt.state
			result := editor.InModalMode()
			if result != tt.expected {
				t.Errorf("InModalMode() = %v, want %v for state %v",
					result, tt.expected, tt.name)
			}
		})
	}
}

func TestEditor_CreateFlow_Integration(t *testing.T) {
	workDir := t.TempDir()
	store, err := jubako.NewStore(workDir)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	registry := pantry.New()
	editor := NewEditorCreate(store, registry)

	// Step 1: Verify initial state
	if editor.state != StateNaming {
		t.Errorf("expected initial state StateNaming, got %v", editor.state)
	}
	if !editor.InModalMode() {
		t.Error("expected InModalMode() = true in StateNaming")
	}

	// Step 2: Simulate name entry
	nameMsg := BentoNameEnteredMsg{Name: "integration-test-bento"}
	editor, _ = editor.Update(nameMsg)

	if editor.bentoName != "integration-test-bento" {
		t.Errorf("expected bentoName='integration-test-bento', got %q", editor.bentoName)
	}
	if editor.state != StateSelectingType {
		t.Errorf("expected state StateSelectingType, got %v", editor.state)
	}
	if !editor.InModalMode() {
		t.Error("expected InModalMode() = true in StateSelectingType")
	}

	// Step 3: Simulate type selection
	typeMsg := NodeTypeSelectedMsg{Type: "http"}
	editor, _ = editor.Update(typeMsg)

	if editor.currentNodeType != "http" {
		t.Errorf("expected currentNodeType='http', got %q", editor.currentNodeType)
	}
	if editor.state != StateConfiguringNode {
		t.Errorf("expected state StateConfiguringNode, got %v", editor.state)
	}
	if !editor.InModalMode() {
		t.Error("expected InModalMode() = true in StateConfiguringNode")
	}

	// Step 4: Simulate node configuration
	nodeMsg := NodeConfiguredMsg{
		Type: "http",
		Name: "Test HTTP Request",
		Parameters: map[string]interface{}{
			"url":    "https://api.example.com",
			"method": "GET",
		},
	}
	editor, _ = editor.Update(nodeMsg)

	if editor.state != StateReview {
		t.Errorf("expected state StateReview, got %v", editor.state)
	}
	if editor.InModalMode() {
		t.Error("expected InModalMode() = false in StateReview")
	}
	if editor.def.Type != "http" {
		t.Errorf("expected def.Type='http', got %q", editor.def.Type)
	}

	// Step 5: Verify definition is complete
	if editor.def.Name != "integration-test-bento" {
		t.Errorf("expected def.Name='integration-test-bento', got %q", editor.def.Name)
	}
	if editor.def.Version != neta.CurrentVersion {
		t.Errorf("expected def.Version=%q, got %q", neta.CurrentVersion, editor.def.Version)
	}
}
