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

func TestEditor_HandleNameEntered(t *testing.T) {
	workDir := t.TempDir()
	store, err := jubako.NewStore(workDir)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	registry := pantry.New()
	editor := NewEditorCreate(store, registry)

	msg := BentoNameEnteredMsg{Name: "test-bento"}
	editor, _ = editor.Update(msg)

	if editor.bentoName != "test-bento" {
		t.Errorf("expected bentoName='test-bento', got %q", editor.bentoName)
	}

	if editor.state != StateSelectingType {
		t.Error("expected StateSelectingType after name entry")
	}
}

func TestEditor_HandleTypeSelected(t *testing.T) {
	workDir := t.TempDir()
	store, err := jubako.NewStore(workDir)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	registry := pantry.New()
	editor := NewEditorCreate(store, registry)
	editor.state = StateSelectingType

	msg := NodeTypeSelectedMsg{Type: "http"}
	editor, _ = editor.Update(msg)

	if editor.currentNodeType != "http" {
		t.Errorf("expected currentNodeType='http', got %q", editor.currentNodeType)
	}

	if editor.state != StateConfiguringNode {
		t.Error("expected StateConfiguringNode after type selection")
	}
}

func TestEditor_HandleNodeConfigured(t *testing.T) {
	workDir := t.TempDir()
	store, err := jubako.NewStore(workDir)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	registry := pantry.New()
	editor := NewEditorCreate(store, registry)
	editor.bentoName = "test-bento"
	editor.def.Name = "test-bento"

	msg := NodeConfiguredMsg{
		Type: "http",
		Name: "Test Request",
		Parameters: map[string]interface{}{
			"url": "https://example.com",
		},
	}
	editor, _ = editor.Update(msg)

	if editor.def.Type == "" {
		t.Error("definition type not set")
	}

	if editor.state != StateReview {
		t.Error("expected StateReview after node configured")
	}
}

func TestEditor_BuildNode(t *testing.T) {
	msg := NodeConfiguredMsg{
		Type: "http",
		Name: "Test Node",
		Parameters: map[string]interface{}{
			"url": "https://example.com",
		},
	}

	node := buildNode(msg)

	if node.Type != "http" {
		t.Errorf("expected type='http', got %q", node.Type)
	}

	if node.Name != "Test Node" {
		t.Errorf("expected name='Test Node', got %q", node.Name)
	}

	if node.Parameters["url"] != "https://example.com" {
		t.Error("parameters not set correctly")
	}
}

func TestEditor_AppendNode(t *testing.T) {
	def := neta.Definition{
		Version: "1.0",
		Name:    "test-bento",
	}

	node := neta.Definition{
		Type: "http",
		Name: "Node 1",
	}

	def = appendNode(def, node)

	if def.Type != "group.sequence" {
		t.Errorf("expected type='group.sequence', got %q", def.Type)
	}

	if len(def.Nodes) != 1 {
		t.Errorf("expected 1 node, got %d", len(def.Nodes))
	}
}
