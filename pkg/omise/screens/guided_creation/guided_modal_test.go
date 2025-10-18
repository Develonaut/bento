package guided_creation

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"

	"bento/pkg/jubako"
	"bento/pkg/neta"
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

	// Verify bento was saved (with retry logic for file system flush)
	var bentos []jubako.BentoInfo
	var def neta.Definition
	var loadErr error

	// Retry up to 5 times with increasing delays
	for i := 0; i < 5; i++ {
		time.Sleep(time.Duration(100+i*100) * time.Millisecond)

		bentos, err = store.List()
		if err != nil {
			continue
		}

		if len(bentos) != 1 {
			continue
		}

		def, loadErr = store.Load(bentos[0].Name)
		if loadErr == nil {
			// Successfully loaded
			break
		}
	}

	if err != nil {
		t.Fatalf("Failed to list bentos: %v", err)
	}

	if len(bentos) != 1 {
		t.Fatalf("Expected 1 bento after creation, got %d", len(bentos))
	}

	if loadErr != nil {
		t.Fatalf("Failed to load bento definition: %v", loadErr)
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

	// Verify bento with retry logic for file system flush
	var bentos []jubako.BentoInfo
	var def neta.Definition
	var loadErr error

	// Retry up to 5 times with increasing delays
	for i := 0; i < 5; i++ {
		time.Sleep(time.Duration(100+i*100) * time.Millisecond)

		bentos, err = store.List()
		if err != nil {
			continue
		}

		if len(bentos) != 1 {
			continue
		}

		def, loadErr = store.Load(bentos[0].Name)
		if loadErr == nil {
			// Successfully loaded
			break
		}
	}

	if err != nil {
		t.Fatalf("Failed to list bentos: %v", err)
	}

	if len(bentos) != 1 {
		t.Fatalf("Expected 1 bento, got %d", len(bentos))
	}

	if loadErr != nil {
		t.Fatalf("Failed to load bento definition: %v", loadErr)
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

// TestGuidedCreation_SimpleSequenceGroup creates a sequence group with no children
func TestGuidedCreation_SimpleSequenceGroup(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Fill metadata
	h.fillMetadata("Simple Sequence", "A sequence group with no children")

	// Select Sequence Group node type (index 3)
	h.selectNodeType(3)

	// Fill sequence node
	h.fillSequenceNode("My Sequence")

	// Group context menu appears - select "Save bento" (index 3)
	h.selectGroupContext(3)

	// Verify
	def := h.verifyBento("Simple Sequence")
	h.assertNodeCount(def, 1)
	h.assertGroupNode(def.Nodes[0], "group.sequence", "My Sequence", 0)

	t.Log("Simple sequence group created successfully")
}

// TestGuidedCreation_SequenceWithChildren creates a sequence with 2 child nodes
func TestGuidedCreation_SequenceWithChildren(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Fill metadata
	h.fillMetadata("Sequence With Children", "A sequence group with child nodes")

	// Create sequence group
	h.selectNodeType(3) // Sequence Group
	h.fillSequenceNode("Data Pipeline")

	// Group context menu - select "Add child to group" (index 0)
	h.selectGroupContext(0)

	// Add first child: HTTP node
	h.selectNodeType(0) // HTTP
	h.fillHTTPNode("Fetch Users", "https://api.example.com/users", "GET")

	// Continue menu - select "Add another node" (index 0)
	h.selectContinue(0)

	// Add second child: Transform node
	h.selectNodeType(1) // Transform.jq
	h.tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("Extract Names")})
	h.tm.Send(tea.KeyMsg{Type: tea.KeyEnter}) // Name
	time.Sleep(50 * time.Millisecond)
	h.tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(".[] | .name")})
	h.tm.Send(tea.KeyMsg{Type: tea.KeyEnter}) // Query
	time.Sleep(50 * time.Millisecond)
	h.tm.Send(tea.KeyMsg{Type: tea.KeyEnter}) // Input (skip)
	time.Sleep(100 * time.Millisecond)

	// Continue menu - select "Done - Save bento" (index 1)
	h.selectContinue(1)

	// Verify
	def := h.verifyBento("Sequence With Children")
	h.assertNodeCount(def, 1)

	// Check sequence group
	sequence := def.Nodes[0]
	h.assertGroupNode(sequence, "group.sequence", "Data Pipeline", 2)

	// Check children
	h.assertChildNode(sequence, 0, "http", "Fetch Users")
	h.assertChildNode(sequence, 1, "transform.jq", "Extract Names")

	t.Log("Sequence with children created successfully")
}

// TestGuidedCreation_NestedGroups creates a sequence containing another sequence
func TestGuidedCreation_NestedGroups(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Fill metadata
	h.fillMetadata("Nested Groups", "A sequence containing another sequence")

	// Create outer sequence
	h.selectNodeType(3) // Sequence Group
	h.fillSequenceNode("Outer Sequence")

	// Add child to outer sequence
	h.selectGroupContext(0) // "Add child to group"

	// Create inner sequence as child
	h.selectNodeType(3) // Sequence Group
	h.fillSequenceNode("Inner Sequence")

	// Add child to inner sequence
	h.selectGroupContext(0) // "Add child to group"

	// Add HTTP node to inner sequence
	h.selectNodeType(0) // HTTP
	h.fillHTTPNode("Nested HTTP", "https://api.example.com/nested", "GET")

	// Done with inner sequence - go back to outer
	h.selectContinue(1) // "Done - Save bento" will save since we're nested

	// The flow should detect we're nested and just complete

	// Verify
	def := h.verifyBento("Nested Groups")
	h.assertNodeCount(def, 1)

	// Check outer sequence
	outerSeq := def.Nodes[0]
	h.assertGroupNode(outerSeq, "group.sequence", "Outer Sequence", 1)

	// Check inner sequence
	innerSeq := h.assertChildNode(outerSeq, 0, "group.sequence", "Inner Sequence")
	if len(innerSeq.Nodes) != 1 {
		t.Errorf("Expected inner sequence to have 1 child, got %d", len(innerSeq.Nodes))
	}

	// Check HTTP node inside inner sequence
	if len(innerSeq.Nodes) > 0 {
		httpNode := innerSeq.Nodes[0]
		h.assertHTTPNode(httpNode, "Nested HTTP", "https://api.example.com/nested", "GET")
	}

	t.Log("Nested groups created successfully")
}

// TestGuidedCreation_ParallelWithChildren creates a parallel group with 3 children
func TestGuidedCreation_ParallelWithChildren(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Fill metadata
	h.fillMetadata("Parallel Processing", "Process multiple requests in parallel")

	// Create parallel group
	h.selectNodeType(4) // Parallel Group
	h.fillParallelNode("Parallel APIs")

	// Add child to parallel group
	h.selectGroupContext(0) // "Add child to group"

	// Add first child: HTTP to API 1
	h.selectNodeType(0) // HTTP
	h.fillHTTPNode("API 1", "https://api1.example.com", "GET")
	h.selectContinue(0) // "Add another node"

	// Add second child: HTTP to API 2
	h.selectNodeType(0) // HTTP
	h.fillHTTPNode("API 2", "https://api2.example.com", "GET")
	h.selectContinue(0) // "Add another node"

	// Add third child: HTTP to API 3
	h.selectNodeType(0) // HTTP
	h.fillHTTPNode("API 3", "https://api3.example.com", "GET")
	h.selectContinue(1) // "Done - Save bento"

	// Verify
	def := h.verifyBento("Parallel Processing")
	h.assertNodeCount(def, 1)

	// Check parallel group
	parallel := def.Nodes[0]
	h.assertGroupNode(parallel, "group.parallel", "Parallel APIs", 3)

	// Check all three children
	h.assertChildNode(parallel, 0, "http", "API 1")
	h.assertChildNode(parallel, 1, "http", "API 2")
	h.assertChildNode(parallel, 2, "http", "API 3")

	t.Log("Parallel with children created successfully")
}

// TestGuidedCreation_MixedGroupsAndNodes creates a complex structure
func TestGuidedCreation_MixedGroupsAndNodes(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Fill metadata
	h.fillMetadata("Mixed Structure", "Groups and nodes at root level")

	// Add first node: HTTP at root
	h.selectNodeType(0) // HTTP
	h.fillHTTPNode("Root HTTP", "https://api.example.com/init", "GET")
	h.selectContinue(0) // "Add another node"

	// Add sequence group at root
	h.selectNodeType(3) // Sequence Group
	h.fillSequenceNode("Processing")

	// Add child to sequence
	h.selectGroupContext(0) // "Add child to group"

	// Add HTTP inside sequence
	h.selectNodeType(0) // HTTP
	h.fillHTTPNode("Process Step", "https://api.example.com/process", "POST")

	// Done with sequence - back to root via "Done with current level"
	h.selectContinue(1) // "Done - Save bento" completes everything

	// Verify
	def := h.verifyBento("Mixed Structure")
	h.assertNodeCount(def, 2)

	// Check root HTTP node
	h.assertHTTPNode(def.Nodes[0], "Root HTTP", "https://api.example.com/init", "GET")

	// Check sequence group
	sequence := def.Nodes[1]
	h.assertGroupNode(sequence, "group.sequence", "Processing", 1)
	h.assertChildNode(sequence, 0, "http", "Process Step")

	t.Log("Mixed structure created successfully")
}

// TestGuidedCreation_BreadcrumbRendering tests the breadcrumb display
func TestGuidedCreation_BreadcrumbRendering(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := jubako.NewStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	modal := NewGuidedModal(store, tmpDir, 120, 40)

	// Test 1: At root level
	breadcrumb := modal.renderBreadcrumb()
	expectedRoot := "Context: Root"
	if !contains(breadcrumb, expectedRoot) {
		t.Errorf("Expected breadcrumb to contain '%s' at root, got: %s", expectedRoot, breadcrumb)
	}

	// Test 2: Push one level deep
	parentNode := &neta.Definition{Name: "Parent Group", Type: "group.sequence"}
	modal.pushParentContext(parentNode)

	breadcrumb = modal.renderBreadcrumb()
	expectedLevel1 := "Context: Root > Parent Group"
	if !contains(breadcrumb, expectedLevel1) {
		t.Errorf("Expected breadcrumb '%s', got: %s", expectedLevel1, breadcrumb)
	}

	// Test 3: Push two levels deep
	childNode := &neta.Definition{Name: "Child Group", Type: "group.sequence"}
	modal.pushParentContext(childNode)

	breadcrumb = modal.renderBreadcrumb()
	expectedLevel2 := "Context: Root > Parent Group > Child Group"
	if !contains(breadcrumb, expectedLevel2) {
		t.Errorf("Expected breadcrumb '%s', got: %s", expectedLevel2, breadcrumb)
	}

	// Test 4: Pop back one level
	popped := modal.popParentContext()
	if !popped {
		t.Error("Expected pop to succeed")
	}

	breadcrumb = modal.renderBreadcrumb()
	if !contains(breadcrumb, expectedLevel1) {
		t.Errorf("Expected breadcrumb '%s' after pop, got: %s", expectedLevel1, breadcrumb)
	}

	// Test 5: Pop back to root
	popped = modal.popParentContext()
	if !popped {
		t.Error("Expected second pop to succeed")
	}

	breadcrumb = modal.renderBreadcrumb()
	if !contains(breadcrumb, expectedRoot) {
		t.Errorf("Expected breadcrumb '%s' after second pop, got: %s", expectedRoot, breadcrumb)
	}

	// Test 6: Try to pop at root (should fail)
	popped = modal.popParentContext()
	if popped {
		t.Error("Expected pop at root to fail")
	}

	t.Log("Breadcrumb rendering test completed successfully")
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestGuidedCreation_DeleteNode tests deleting a node with Ctrl+D
func TestGuidedCreation_DeleteNode(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Fill metadata
	h.fillMetadata("Test Bento", "Testing delete")
	time.Sleep(100 * time.Millisecond)

	// Add a complete HTTP node first
	h.selectNodeType(0)
	h.fillHTTPNode("fetch_data", "https://api.example.com", "GET")
	time.Sleep(100 * time.Millisecond)

	// Add another node
	h.selectContinue(0) // add
	time.Sleep(100 * time.Millisecond)

	// Select file.write
	h.selectNodeType(2)
	time.Sleep(100 * time.Millisecond)

	// Fill node name but then delete it
	h.tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("test_node")})
	time.Sleep(50 * time.Millisecond)

	// Delete the node with Ctrl+D
	h.tm.Send(tea.KeyMsg{Type: tea.KeyCtrlD})
	time.Sleep(100 * time.Millisecond)

	// Now save
	h.selectContinue(1) // done
	time.Sleep(300 * time.Millisecond)

	// Verify only 1 node was saved (the delete worked)
	def := h.verifyBento("Test Bento")
	h.assertNodeCount(def, 1)
	h.assertHTTPNode(def.Nodes[0], "fetch_data", "https://api.example.com", "GET")
}

// TestGuidedCreation_HistoryCapture tests that history is captured correctly
func TestGuidedCreation_HistoryCapture(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := jubako.NewStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	modal := NewGuidedModal(store, tmpDir, 120, 40)

	// Manually test history functions
	// Start at metadata stage
	if modal.stage != guidedStageMetadata {
		t.Fatalf("Expected to start at metadata stage")
	}

	// Capture snapshot
	modal.captureSnapshot()

	// Should have captured metadata stage
	if len(modal.history.stages) != 1 {
		t.Errorf("Expected 1 history entry after capture, got %d", len(modal.history.stages))
	}

	if len(modal.history.stages) > 0 && modal.history.stages[0] != guidedStageMetadata {
		t.Errorf("Expected first history entry to be metadata stage")
	}

	// Move to next stage and capture again
	modal.stage = guidedStageNodeTypeSelect
	modal.captureSnapshot()

	if len(modal.history.stages) != 2 {
		t.Errorf("Expected 2 history entries after second capture, got %d", len(modal.history.stages))
	}

	t.Log("History capture test completed successfully")
}

// TestGuidedCreation_EditMenu tests the edit menu functionality
func TestGuidedCreation_EditMenu(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := jubako.NewStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// Create an existing bento to edit
	originalDef := neta.Definition{
		Version:     "1.0",
		Type:        "group.sequence",
		Name:        "Original Bento",
		Description: "Original description",
		Icon:        "🍱",
		Nodes: []neta.Definition{
			{
				Version: "1.0",
				Type:    "http",
				Name:    "test_node",
				Parameters: map[string]interface{}{
					"url":    "https://example.com",
					"method": "GET",
				},
			},
		},
	}

	if err := store.Save("original-bento", originalDef); err != nil {
		t.Fatalf("Failed to save original bento: %v", err)
	}

	// Load it for editing
	modal, err := NewGuidedModalForEdit(store, tmpDir, 120, 40, "original-bento")
	if err != nil {
		t.Fatalf("Failed to create edit modal: %v", err)
	}

	// Should start at edit menu
	if modal.stage != guidedStageEditMenu {
		t.Errorf("Expected to start at edit menu stage, got %d", modal.stage)
	}

	// Should be in editing mode
	if !modal.editing {
		t.Error("Expected editing flag to be true")
	}

	// Definition should be loaded
	if modal.definition.Name != "Original Bento" {
		t.Errorf("Expected name 'Original Bento', got '%s'", modal.definition.Name)
	}

	if len(modal.definition.Nodes) != 1 {
		t.Errorf("Expected 1 node, got %d", len(modal.definition.Nodes))
	}

	t.Log("Edit menu initialization test completed successfully")
}

// TestGuidedCreation_EditMetadata tests editing bento metadata
func TestGuidedCreation_EditMetadata(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := jubako.NewStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// Create an existing bento
	originalDef := neta.Definition{
		Version:     "1.0",
		Type:        "group.sequence",
		Name:        "Test Bento",
		Description: "Original description",
		Icon:        "🍱",
		Nodes:       []neta.Definition{},
	}

	if err := store.Save("test-bento", originalDef); err != nil {
		t.Fatalf("Failed to save bento: %v", err)
	}

	// Load for editing
	modal, err := NewGuidedModalForEdit(store, tmpDir, 120, 40, "test-bento")
	if err != nil {
		t.Fatalf("Failed to create edit modal: %v", err)
	}

	tm := teatest.NewTestModel(t, modal, teatest.WithInitialTermSize(120, 40))
	defer func() { _ = tm.Quit() }()

	// Select "Edit metadata" option (first option)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(100 * time.Millisecond)

	// Should now be on metadata form
	// Change name
	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlU}) // Clear field
	time.Sleep(50 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("Updated Bento")})
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(50 * time.Millisecond)

	// Change description
	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlU}) // Clear field
	time.Sleep(50 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("Updated description")})
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(100 * time.Millisecond)

	// Should be back at edit menu
	// Now save
	tm.Send(tea.KeyMsg{Type: tea.KeyDown}) // Move to "Add a new node"
	time.Sleep(20 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyDown}) // Move to "Delete a node"
	time.Sleep(20 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyDown}) // Move to "Save and exit"
	time.Sleep(20 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyDown}) // Move past delete
	time.Sleep(20 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(500 * time.Millisecond)

	// Verify metadata was updated
	var def neta.Definition
	var loadErr error

	for i := 0; i < 5; i++ {
		time.Sleep(time.Duration(100+i*100) * time.Millisecond)
		def, loadErr = store.Load("test-bento")
		if loadErr == nil {
			break
		}
	}

	if loadErr != nil {
		t.Fatalf("Failed to load updated bento: %v", loadErr)
	}

	if def.Name != "Updated Bento" {
		t.Errorf("Expected name 'Updated Bento', got '%s'", def.Name)
	}

	if def.Description != "Updated description" {
		t.Errorf("Expected description 'Updated description', got '%s'", def.Description)
	}

	// Version should be incremented
	if def.Version != "1.1" {
		t.Errorf("Expected version 1.1, got %s", def.Version)
	}

	t.Log("Edit metadata test completed successfully")
}

// TestGuidedCreation_DeleteNode_EditMenu tests deleting a node via edit menu
func TestGuidedCreation_DeleteNode_EditMenu(t *testing.T) {

	tmpDir := t.TempDir()
	store, err := jubako.NewStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// Create a bento with 2 nodes
	originalDef := neta.Definition{
		Version: "1.0",
		Type:    "group.sequence",
		Name:    "Delete Test",
		Icon:    "🍱",
		Nodes: []neta.Definition{
			{
				Version: "1.0",
				Type:    "http",
				Name:    "node_to_keep",
				Parameters: map[string]interface{}{
					"url":    "https://keep.com",
					"method": "GET",
				},
			},
			{
				Version: "1.0",
				Type:    "http",
				Name:    "node_to_delete",
				Parameters: map[string]interface{}{
					"url":    "https://delete.com",
					"method": "POST",
				},
			},
		},
	}

	if err := store.Save("delete-test", originalDef); err != nil {
		t.Fatalf("Failed to save bento: %v", err)
	}

	// Load for editing
	modal, err := NewGuidedModalForEdit(store, tmpDir, 120, 40, "delete-test")
	if err != nil {
		t.Fatalf("Failed to create edit modal: %v", err)
	}

	tm := teatest.NewTestModel(t, modal, teatest.WithInitialTermSize(120, 40))
	defer func() { _ = tm.Quit() }()

	// Navigate to "Delete a node"
	tm.Send(tea.KeyMsg{Type: tea.KeyDown}) // Edit node
	time.Sleep(20 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyDown}) // Add node
	time.Sleep(20 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyDown}) // Delete node
	time.Sleep(20 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(100 * time.Millisecond)

	// Should show node list
	// Select first node (default selection) which is node_to_keep
	// Press Enter to delete it
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(100 * time.Millisecond)

	// Should be back at edit menu
	// Now save
	tm.Send(tea.KeyMsg{Type: tea.KeyDown}) // Edit node
	time.Sleep(20 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyDown}) // Add node
	time.Sleep(20 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyDown}) // Delete node
	time.Sleep(20 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyDown}) // Save
	time.Sleep(20 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(500 * time.Millisecond)

	// Verify node was deleted
	var def neta.Definition
	var loadErr error

	for i := 0; i < 5; i++ {
		time.Sleep(time.Duration(100+i*100) * time.Millisecond)
		def, loadErr = store.Load("delete-test")
		if loadErr == nil {
			break
		}
	}

	if loadErr != nil {
		t.Fatalf("Failed to load updated bento: %v", loadErr)
	}

	if len(def.Nodes) != 1 {
		t.Fatalf("Expected 1 node after deletion, got %d", len(def.Nodes))
	}

	if def.Nodes[0].Name != "node_to_delete" {
		t.Errorf("Expected remaining node to be 'node_to_delete', got '%s'", def.Nodes[0].Name)
	}

	t.Log("Delete node test completed successfully")
}

// TestGuidedCreation_EditNodeParameters tests editing node parameters with pre-filled values
func TestGuidedCreation_EditNodeParameters(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := jubako.NewStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// Create a bento with an HTTP node
	originalDef := neta.Definition{
		Version: "1.0",
		Type:    "group.sequence",
		Name:    "Edit Params Test",
		Icon:    "🍱",
		Nodes: []neta.Definition{
			{
				Version: "1.0",
				Type:    "http",
				Name:    "api_call",
				Parameters: map[string]interface{}{
					"url":    "https://original.com",
					"method": "GET",
				},
			},
		},
	}

	if err := store.Save("edit-params-test", originalDef); err != nil {
		t.Fatalf("Failed to save bento: %v", err)
	}

	// Load for editing
	modal, err := NewGuidedModalForEdit(store, tmpDir, 120, 40, "edit-params-test")
	if err != nil {
		t.Fatalf("Failed to create edit modal: %v", err)
	}

	tm := teatest.NewTestModel(t, modal, teatest.WithInitialTermSize(120, 40))
	defer func() { _ = tm.Quit() }()

	// Navigate to "Edit an existing node"
	tm.Send(tea.KeyMsg{Type: tea.KeyDown}) // Edit an existing node
	time.Sleep(20 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(100 * time.Millisecond)

	// Should show node list with api_call pre-selected
	// Press Enter to edit it
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(100 * time.Millisecond)

	// Should now be on HTTP node edit form with pre-filled values
	// The form should have:
	// - Name: api_call
	// - URL: https://original.com
	// - Method: GET

	// Change the URL by clearing and entering new value
	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlU}) // Clear name field
	time.Sleep(50 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("updated_api_call")})
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(50 * time.Millisecond)

	// Change URL
	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlU}) // Clear URL field
	time.Sleep(50 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("https://updated.com")})
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(50 * time.Millisecond)

	// Keep method as GET (just press Enter to skip)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(50 * time.Millisecond)

	// Skip headers (Enter)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(50 * time.Millisecond)

	// Skip body (Enter)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(100 * time.Millisecond)

	// Should be back at edit menu
	// Navigate to Save
	tm.Send(tea.KeyMsg{Type: tea.KeyDown}) // Edit node
	time.Sleep(20 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyDown}) // Add node
	time.Sleep(20 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyDown}) // Delete node
	time.Sleep(20 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyDown}) // Save
	time.Sleep(20 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(500 * time.Millisecond)

	// Verify node was updated
	var def neta.Definition
	var loadErr error

	for i := 0; i < 5; i++ {
		time.Sleep(time.Duration(100+i*100) * time.Millisecond)
		def, loadErr = store.Load("edit-params-test")
		if loadErr == nil {
			break
		}
	}

	if loadErr != nil {
		t.Fatalf("Failed to load updated bento: %v", loadErr)
	}

	if len(def.Nodes) != 1 {
		t.Fatalf("Expected 1 node, got %d", len(def.Nodes))
	}

	// Verify node was updated
	node := def.Nodes[0]
	if node.Name != "updated_api_call" {
		t.Errorf("Expected name 'updated_api_call', got '%s'", node.Name)
	}

	if url, ok := node.Parameters["url"].(string); !ok || url != "https://updated.com" {
		t.Errorf("Expected URL 'https://updated.com', got '%v'", node.Parameters["url"])
	}

	if method, ok := node.Parameters["method"].(string); !ok || method != "GET" {
		t.Errorf("Expected method 'GET', got '%v'", node.Parameters["method"])
	}

	t.Log("Edit node parameters test completed successfully")
}

// TestGuidedCreation_NoChangeNoVersionIncrement tests that version doesn't increment if nothing changed
func TestGuidedCreation_NoChangeNoVersionIncrement(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := jubako.NewStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// Create a bento
	originalDef := neta.Definition{
		Version: "1.0",
		Type:    "group.sequence",
		Name:    "No Change Test",
		Icon:    "🍱",
		Nodes: []neta.Definition{
			{
				Version: "1.0",
				Type:    "http",
				Name:    "api_call",
				Parameters: map[string]interface{}{
					"url":    "https://example.com",
					"method": "GET",
				},
			},
		},
	}

	if err := store.Save("no-change-test", originalDef); err != nil {
		t.Fatalf("Failed to save bento: %v", err)
	}

	// Load for editing
	modal, err := NewGuidedModalForEdit(store, tmpDir, 120, 40, "no-change-test")
	if err != nil {
		t.Fatalf("Failed to create edit modal: %v", err)
	}

	tm := teatest.NewTestModel(t, modal, teatest.WithInitialTermSize(120, 40))
	defer func() { _ = tm.Quit() }()

	// Just save without making any changes
	// Navigate to "Save and exit" (index 4)
	tm.Send(tea.KeyMsg{Type: tea.KeyDown}) // Edit node
	time.Sleep(20 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyDown}) // Add node
	time.Sleep(20 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyDown}) // Delete node
	time.Sleep(20 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyDown}) // Save
	time.Sleep(20 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(500 * time.Millisecond)

	// Verify version was NOT incremented
	var def neta.Definition
	var loadErr error

	for i := 0; i < 5; i++ {
		time.Sleep(time.Duration(100+i*100) * time.Millisecond)
		def, loadErr = store.Load("no-change-test")
		if loadErr == nil {
			break
		}
	}

	if loadErr != nil {
		t.Fatalf("Failed to load bento: %v", loadErr)
	}

	// Version should still be 1.0 (not incremented)
	if def.Version != "1.0" {
		t.Errorf("Expected version to stay at 1.0 when no changes made, got %s", def.Version)
	}

	t.Log("No change no version increment test completed successfully")
}
