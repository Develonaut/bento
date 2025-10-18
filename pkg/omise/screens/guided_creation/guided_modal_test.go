package guided_creation

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"

	"bento/pkg/jubako"
)

// TestGuidedCreation_CompleteFlow tests the entire guided creation flow
func TestGuidedCreation_CompleteFlow(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := jubako.NewStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	modal := NewGuidedModal(store, tmpDir, 120, 40)

	tm := teatest.NewTestModel(
		t,
		modal,
		teatest.WithInitialTermSize(120, 40),
	)
	defer func() {
		if err := tm.Quit(); err != nil {
			t.Logf("Failed to quit test model: %v", err)
		}
	}()

	// Stage 1: Bento Metadata
	// Icon field is first (select)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter}) // Accept default icon
	time.Sleep(50 * time.Millisecond)

	// Name field
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("My Test Bento")})
	time.Sleep(50 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(50 * time.Millisecond)

	// Description field
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("A test workflow")})
	time.Sleep(50 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(100 * time.Millisecond)

	// Stage 2: Node Type Selection
	// Should now be on node type selection
	// Default is "HTTP Request" (first option)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(100 * time.Millisecond)

	// Stage 3: HTTP Node Parameters
	// Name field
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("Fetch Data")})
	time.Sleep(50 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(50 * time.Millisecond)

	// URL field
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("https://api.example.com/data")})
	time.Sleep(50 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(50 * time.Millisecond)

	// Method field (default GET)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(50 * time.Millisecond)

	// Headers field (optional, skip)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(50 * time.Millisecond)

	// Body field (optional, skip)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(100 * time.Millisecond)

	// Stage 4: Continue prompt
	// Should now show "Add another node" or "Done - Save bento"
	// Select "Done - Save bento" (second option)
	tm.Send(tea.KeyMsg{Type: tea.KeyDown})
	time.Sleep(50 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(100 * time.Millisecond)

	// Verify bento was saved
	bentos, err := store.List()
	if err != nil {
		t.Fatalf("Failed to list bentos: %v", err)
	}

	if len(bentos) != 1 {
		t.Fatalf("Expected 1 bento after creation, got %d", len(bentos))
	}

	bentoInfo := bentos[0]

	// Load full definition to verify details (use Name which is the filename)
	def, err := store.Load(bentoInfo.Name)
	if err != nil {
		t.Fatalf("Failed to load bento definition: %v", err)
	}

	if def.Description != "A test workflow" {
		t.Errorf("Expected description 'A test workflow', got '%s'", def.Description)
	}

	if len(def.Nodes) != 1 {
		t.Fatalf("Expected 1 node, got %d", len(def.Nodes))
	}

	node := def.Nodes[0]
	if node.Name != "Fetch Data" {
		t.Errorf("Expected node name 'Fetch Data', got '%s'", node.Name)
	}

	if node.Type != "http" {
		t.Errorf("Expected node type 'http', got '%s'", node.Type)
	}

	url, ok := node.Parameters["url"].(string)
	if !ok || url != "https://api.example.com/data" {
		t.Errorf("Expected URL 'https://api.example.com/data', got '%v'", node.Parameters["url"])
	}

	t.Log("Guided creation flow completed successfully")
}

// TestGuidedCreation_MultipleNodes tests adding multiple nodes
func TestGuidedCreation_MultipleNodes(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := jubako.NewStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	modal := NewGuidedModal(store, tmpDir, 120, 40)

	tm := teatest.NewTestModel(
		t,
		modal,
		teatest.WithInitialTermSize(120, 40),
	)
	defer func() {
		if err := tm.Quit(); err != nil {
			t.Logf("Failed to quit test model: %v", err)
		}
	}()

	// Stage 1: Bento Metadata
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter}) // Icon
	time.Sleep(50 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("Multi-Node Bento")})
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter}) // Name
	time.Sleep(50 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("Test with multiple nodes")})
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter}) // Description
	time.Sleep(100 * time.Millisecond)

	// First Node: HTTP
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter}) // Select HTTP (default)
	time.Sleep(100 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("Fetch API")})
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter}) // Name
	time.Sleep(50 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("https://api.example.com")})
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter}) // URL
	time.Sleep(50 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter}) // Method (default GET)
	time.Sleep(50 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter}) // Headers (skip)
	time.Sleep(50 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter}) // Body (skip)
	time.Sleep(100 * time.Millisecond)

	// Continue: Add another node
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter}) // Select "Add another node" (default/first option)
	time.Sleep(100 * time.Millisecond)

	// Second Node: Transform.jq
	// Select transform.jq (second option in type list)
	tm.Send(tea.KeyMsg{Type: tea.KeyDown})
	time.Sleep(50 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(100 * time.Millisecond)

	// Fill jq parameters
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("Transform Data")})
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter}) // Name
	time.Sleep(50 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(".data")})
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter}) // Query
	time.Sleep(50 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter}) // Input (skip)
	time.Sleep(100 * time.Millisecond)

	// Continue: Done - Save bento
	tm.Send(tea.KeyMsg{Type: tea.KeyDown}) // Move to "Done"
	time.Sleep(50 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(100 * time.Millisecond)

	// Give file system a moment to flush
	time.Sleep(100 * time.Millisecond)

	// Verify bento with 2 nodes
	bentos, err := store.List()
	if err != nil {
		t.Fatalf("Failed to list bentos: %v", err)
	}

	if len(bentos) != 1 {
		t.Fatalf("Expected 1 bento, got %d", len(bentos))
	}

	bentoInfo := bentos[0]
	def, err := store.Load(bentoInfo.Name)
	if err != nil {
		t.Fatalf("Failed to load bento definition: %v", err)
	}

	if len(def.Nodes) != 2 {
		t.Fatalf("Expected 2 nodes, got %d", len(def.Nodes))
	}

	// Verify first node (HTTP)
	if def.Nodes[0].Type != "http" {
		t.Errorf("First node should be http, got %s", def.Nodes[0].Type)
	}

	// Verify second node (transform.jq)
	if def.Nodes[1].Type != "transform.jq" {
		t.Errorf("Second node should be transform.jq, got %s", def.Nodes[1].Type)
	}

	t.Log("Multiple nodes test completed successfully")
}

// TestGuidedCreation_CancelWithEscape tests ESC cancellation
func TestGuidedCreation_CancelWithEscape(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := jubako.NewStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	modal := NewGuidedModal(store, tmpDir, 120, 40)

	tm := teatest.NewTestModel(
		t,
		modal,
		teatest.WithInitialTermSize(120, 40),
	)
	defer func() {
		if err := tm.Quit(); err != nil {
			t.Logf("Failed to quit test model: %v", err)
		}
	}()

	// Start filling metadata
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter}) // Icon
	time.Sleep(50 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("Cancelled Bento")})
	time.Sleep(50 * time.Millisecond)

	// Press ESC to cancel
	tm.Send(tea.KeyMsg{Type: tea.KeyEsc})
	time.Sleep(100 * time.Millisecond)

	// Verify no bento was created
	bentos, err := store.List()
	if err != nil {
		t.Fatalf("Failed to list bentos: %v", err)
	}

	if len(bentos) != 0 {
		t.Errorf("Expected 0 bentos after cancellation, got %d", len(bentos))
	}

	t.Log("Cancellation test completed successfully")
}

// TestGuidedCreation_HelloWorldHTTP creates a simple hello-world HTTP bento
func TestGuidedCreation_HelloWorldHTTP(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Fill metadata
	h.fillMetadata("Hello World HTTP", "Fetch hello world from httpbin")

	// Select HTTP node type (index 0)
	h.selectNodeType(0)

	// Fill HTTP node
	h.fillHTTPNode("Get Hello", "https://httpbin.org/get", "GET")

	// Done - save bento
	h.selectContinue(1)

	// Verify
	def := h.verifyBento("Hello World HTTP")
	h.assertNodeCount(def, 1)
	h.assertHTTPNode(def.Nodes[0], "Get Hello", "https://httpbin.org/get", "GET")

	t.Log("Hello World HTTP bento created successfully")
}

// TestGuidedCreation_HelloWorldFile creates a simple hello-world file write bento
func TestGuidedCreation_HelloWorldFile(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Fill metadata
	h.fillMetadata("Hello World File", "Write hello world to a file")

	// Select File Write node type (index 2)
	h.selectNodeType(2)

	// Fill file write node
	h.fillFileWriteNode("Write Greeting", "/tmp/hello.txt", "Hello, World!")

	// Done - save bento
	h.selectContinue(1)

	// Verify
	def := h.verifyBento("Hello World File")
	h.assertNodeCount(def, 1)
	h.assertFileWriteNode(def.Nodes[0], "Write Greeting", "/tmp/hello.txt", "Hello, World!")

	t.Log("Hello World File bento created successfully")
}
