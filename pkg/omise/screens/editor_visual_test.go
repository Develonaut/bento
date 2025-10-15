package screens

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"bento/pkg/jubako"
	"bento/pkg/neta"
	"bento/pkg/pantry"
)

func createTestEditorWithNodes() Editor {
	workDir := "testdata"
	store, _ := jubako.NewStore(workDir)
	registry := pantry.New()

	editor := NewEditorCreate(store, registry)
	editor.bentoName = "test-bento"
	editor.state = StateReview

	// Add test nodes
	editor.def.Nodes = []neta.Definition{
		{Type: "http", Name: "Node 1", Version: neta.CurrentVersion},
		{Type: "transform.jq", Name: "Node 2", Version: neta.CurrentVersion},
		{Type: "conditional.if", Name: "Node 3", Version: neta.CurrentVersion},
	}
	editor.def.Type = "group.sequence"

	return editor
}

func TestEditor_NavigateDown(t *testing.T) {
	editor := createTestEditorWithNodes()

	if editor.selectedNodeIndex != 0 {
		t.Errorf("expected initial selection 0, got %d", editor.selectedNodeIndex)
	}

	// Navigate down
	editor = editor.navigateDown()

	if editor.selectedNodeIndex != 1 {
		t.Errorf("expected selection 1 after down, got %d", editor.selectedNodeIndex)
	}

	// Navigate down again
	editor = editor.navigateDown()

	if editor.selectedNodeIndex != 2 {
		t.Errorf("expected selection 2 after second down, got %d", editor.selectedNodeIndex)
	}

	// Try to navigate past end
	editor = editor.navigateDown()

	if editor.selectedNodeIndex != 2 {
		t.Errorf("expected selection to stay at 2, got %d", editor.selectedNodeIndex)
	}
}

func TestEditor_NavigateUp(t *testing.T) {
	editor := createTestEditorWithNodes()
	editor.selectedNodeIndex = 2

	// Navigate up
	editor = editor.navigateUp()

	if editor.selectedNodeIndex != 1 {
		t.Errorf("expected selection 1 after up, got %d", editor.selectedNodeIndex)
	}

	// Navigate up again
	editor = editor.navigateUp()

	if editor.selectedNodeIndex != 0 {
		t.Errorf("expected selection 0 after second up, got %d", editor.selectedNodeIndex)
	}

	// Try to navigate before start
	editor = editor.navigateUp()

	if editor.selectedNodeIndex != 0 {
		t.Errorf("expected selection to stay at 0, got %d", editor.selectedNodeIndex)
	}
}

func TestEditor_MoveNode(t *testing.T) {
	editor := createTestEditorWithNodes()
	editor.selectedNodeIndex = 0

	originalFirstNode := editor.def.Nodes[0].Name
	originalSecondNode := editor.def.Nodes[1].Name

	// Move first node down
	editor, _ = editor.moveNode(0)

	if editor.def.Nodes[0].Name != originalSecondNode {
		t.Errorf("expected first node to be %s, got %s", originalSecondNode, editor.def.Nodes[0].Name)
	}

	if editor.def.Nodes[1].Name != originalFirstNode {
		t.Errorf("expected second node to be %s, got %s", originalFirstNode, editor.def.Nodes[1].Name)
	}
}

func TestEditor_MoveNodeAtEnd(t *testing.T) {
	editor := createTestEditorWithNodes()
	editor.selectedNodeIndex = 2

	// Try to move last node (should do nothing)
	editor, _ = editor.moveNode(2)

	if editor.def.Nodes[2].Name != "Node 3" {
		t.Error("last node should not have moved")
	}

	if editor.message != "Cannot move last node down" {
		t.Errorf("expected message about moving last node, got: %s", editor.message)
	}
}

func TestEditor_DeleteNode(t *testing.T) {
	editor := createTestEditorWithNodes()
	editor.selectedNodeIndex = 1

	// Delete middle node
	editor, _ = editor.deleteNode(1)

	if len(editor.def.Nodes) != 2 {
		t.Errorf("expected 2 nodes after delete, got %d", len(editor.def.Nodes))
	}

	if editor.def.Nodes[0].Name != "Node 1" {
		t.Error("first node should still be Node 1")
	}

	if editor.def.Nodes[1].Name != "Node 3" {
		t.Error("second node should now be Node 3")
	}
}

func TestEditor_DeleteNodeAdjustsSelection(t *testing.T) {
	editor := createTestEditorWithNodes()
	editor.selectedNodeIndex = 2

	// Delete last node
	editor, _ = editor.deleteNode(2)

	if editor.selectedNodeIndex != 1 {
		t.Errorf("expected selection adjusted to 1, got %d", editor.selectedNodeIndex)
	}
}

func TestEditor_ToggleViewMode(t *testing.T) {
	editor := createTestEditorWithNodes()

	if editor.viewMode != ViewModeList {
		t.Error("expected initial view mode to be List")
	}

	// Toggle to visual
	editor = editor.toggleViewMode()

	if editor.viewMode != ViewModeVisual {
		t.Error("expected view mode to toggle to Visual")
	}

	// Toggle back to list
	editor = editor.toggleViewMode()

	if editor.viewMode != ViewModeList {
		t.Error("expected view mode to toggle back to List")
	}
}

func TestEditor_GetNode(t *testing.T) {
	editor := createTestEditorWithNodes()

	node := editor.getNode(1)

	if node == nil {
		t.Fatal("expected node, got nil")
	}

	if node.Name != "Node 2" {
		t.Errorf("expected Node 2, got %s", node.Name)
	}

	// Test out of bounds
	node = editor.getNode(10)
	if node != nil {
		t.Error("expected nil for out of bounds index")
	}
}

func TestEditor_NavigationKeys(t *testing.T) {
	editor := createTestEditorWithNodes()

	// Test down arrow
	editor, _ = editor.handleReviewKey(tea.KeyMsg{Type: tea.KeyDown})
	if editor.selectedNodeIndex != 1 {
		t.Error("down arrow should navigate down")
	}

	// Test up arrow
	editor, _ = editor.handleReviewKey(tea.KeyMsg{Type: tea.KeyUp})
	if editor.selectedNodeIndex != 0 {
		t.Error("up arrow should navigate up")
	}

	// Test 'j' key
	editor, _ = editor.handleReviewKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if editor.selectedNodeIndex != 1 {
		t.Error("'j' key should navigate down")
	}

	// Test 'k' key
	editor, _ = editor.handleReviewKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if editor.selectedNodeIndex != 0 {
		t.Error("'k' key should navigate up")
	}
}

func TestEditor_GetNodes(t *testing.T) {
	editor := createTestEditorWithNodes()

	nodes := editor.getNodes()
	if len(nodes) != 3 {
		t.Errorf("expected 3 nodes, got %d", len(nodes))
	}

	// Test single node bento
	editor.def = neta.Definition{
		Type:    "http",
		Name:    "Single",
		Version: neta.CurrentVersion,
	}

	nodes = editor.getNodes()
	if len(nodes) != 1 {
		t.Errorf("expected 1 node for single-node bento, got %d", len(nodes))
	}

	if nodes[0].Type != "http" {
		t.Errorf("expected node type 'http', got %s", nodes[0].Type)
	}
}
