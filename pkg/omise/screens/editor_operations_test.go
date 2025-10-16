package screens

import (
	"testing"

	"bento/pkg/jubako"
	"bento/pkg/neta"
	"bento/pkg/pantry"
)

// Handler tests

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

// Node operation tests

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

func TestEditor_AppendNode_EmptyDef(t *testing.T) {
	// Test appendNode when def.Type is empty (converts to group.sequence)
	def := neta.Definition{
		Version: "1.0",
		Name:    "test-bento",
		// Type is empty
	}

	node := neta.Definition{
		Type: "http",
		Name: "Node 1",
		Parameters: map[string]interface{}{
			"url": "https://example.com/1",
		},
	}

	def = appendNode(def, node)

	if def.Type != "group.sequence" {
		t.Errorf("expected def.Type='group.sequence', got %q", def.Type)
	}

	if len(def.Nodes) != 1 {
		t.Errorf("expected 1 node, got %d", len(def.Nodes))
	}
}

func TestEditor_AppendNode_PreservesExistingType(t *testing.T) {
	// Test appendNode when def.Type already exists
	// (current behavior: doesn't convert, may be improved in future)
	def := neta.Definition{
		Version: "1.0",
		Name:    "test-bento",
		Type:    "http", // Already has a type
	}

	node := neta.Definition{
		Type: "http",
		Name: "Node 2",
		Parameters: map[string]interface{}{
			"url": "https://example.com/2",
		},
	}

	def = appendNode(def, node)

	// Current behavior: Type is preserved as-is
	if def.Type != "http" {
		t.Errorf("expected def.Type='http', got %q", def.Type)
	}

	// Node is appended to Nodes array
	if len(def.Nodes) != 1 {
		t.Errorf("expected 1 node, got %d", len(def.Nodes))
	}
}
