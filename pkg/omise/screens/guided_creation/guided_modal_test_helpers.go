package guided_creation

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"

	"bento/pkg/jubako"
	"bento/pkg/neta"
)

// testHelper provides common test utilities for guided creation tests
type testHelper struct {
	t     *testing.T
	tm    *teatest.TestModel
	store *jubako.Store
}

// newTestHelper creates a new test helper
func newTestHelper(t *testing.T) *testHelper {
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

	return &testHelper{
		t:     t,
		tm:    tm,
		store: store,
	}
}

// cleanup cleans up test resources
func (h *testHelper) cleanup() {
	if err := h.tm.Quit(); err != nil {
		h.t.Logf("Failed to quit test model: %v", err)
	}
}

// fillMetadata fills the bento metadata form
func (h *testHelper) fillMetadata(name, description string) {
	// Icon (default)
	h.tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(50 * time.Millisecond)

	// Name
	h.tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(name)})
	h.tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(50 * time.Millisecond)

	// Description
	h.tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(description)})
	h.tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(100 * time.Millisecond)
}

// selectNodeType selects a node type by index (0 = HTTP, 1 = jq, 2 = file.write, etc.)
func (h *testHelper) selectNodeType(index int) {
	for i := 0; i < index; i++ {
		h.tm.Send(tea.KeyMsg{Type: tea.KeyDown})
		time.Sleep(50 * time.Millisecond)
	}
	h.tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(100 * time.Millisecond)
}

// fillHTTPNode fills an HTTP node form
func (h *testHelper) fillHTTPNode(name, url, method string) {
	// Name
	h.tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(name)})
	h.tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(50 * time.Millisecond)

	// URL
	h.tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(url)})
	h.tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(50 * time.Millisecond)

	// Method (select or use default GET)
	if method != "" && method != "GET" {
		// Move to desired method
		methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
		for i, m := range methods {
			if m == method {
				for j := 0; j < i; j++ {
					h.tm.Send(tea.KeyMsg{Type: tea.KeyDown})
					time.Sleep(20 * time.Millisecond)
				}
				break
			}
		}
	}
	h.tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(50 * time.Millisecond)

	// Headers (skip)
	h.tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(50 * time.Millisecond)

	// Body (skip)
	h.tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(100 * time.Millisecond)
}

// fillFileWriteNode fills a file.write node form
func (h *testHelper) fillFileWriteNode(name, path, content string) {
	// Name
	h.tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(name)})
	h.tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(50 * time.Millisecond)

	// Path
	h.tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(path)})
	h.tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(50 * time.Millisecond)

	// Content
	if content != "" {
		h.tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(content)})
	}
	h.tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(100 * time.Millisecond)
}

// selectContinue selects either "add" (0) or "done" (1) on the continue prompt
func (h *testHelper) selectContinue(choice int) {
	for i := 0; i < choice; i++ {
		h.tm.Send(tea.KeyMsg{Type: tea.KeyDown})
		time.Sleep(50 * time.Millisecond)
	}
	h.tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(100 * time.Millisecond)
}

// verifyBento loads and verifies a saved bento with retry logic for file system flushes
func (h *testHelper) verifyBento(expectedName string) neta.Definition {
	var def neta.Definition
	var bentos []jubako.BentoInfo
	var err error

	// Retry up to 5 times with increasing delays
	for i := 0; i < 5; i++ {
		time.Sleep(time.Duration(100+i*100) * time.Millisecond)

		bentos, err = h.store.List()
		if err != nil {
			continue
		}

		if len(bentos) != 1 {
			continue
		}

		def, err = h.store.Load(bentos[0].Name)
		if err == nil {
			// Successfully loaded
			break
		}
	}

	if err != nil {
		h.t.Fatalf("Failed to load bento after retries: %v", err)
	}

	if len(bentos) != 1 {
		h.t.Fatalf("Expected 1 bento, got %d", len(bentos))
	}

	if def.Name != expectedName {
		h.t.Errorf("Expected name '%s', got '%s'", expectedName, def.Name)
	}

	return def
}

// assertNodeCount verifies the number of nodes in a definition
func (h *testHelper) assertNodeCount(def neta.Definition, expected int) {
	if len(def.Nodes) != expected {
		h.t.Fatalf("Expected %d nodes, got %d", expected, len(def.Nodes))
	}
}

// assertHTTPNode verifies an HTTP node's properties
func (h *testHelper) assertHTTPNode(node neta.Definition, name, url, method string) {
	if node.Type != "http" {
		h.t.Errorf("Expected http node, got '%s'", node.Type)
	}

	if node.Name != name {
		h.t.Errorf("Expected node name '%s', got '%s'", name, node.Name)
	}

	if actualURL, ok := node.Parameters["url"].(string); !ok || actualURL != url {
		h.t.Errorf("Expected URL '%s', got '%v'", url, node.Parameters["url"])
	}

	if actualMethod, ok := node.Parameters["method"].(string); !ok || actualMethod != method {
		h.t.Errorf("Expected method '%s', got '%v'", method, node.Parameters["method"])
	}
}

// assertFileWriteNode verifies a file.write node's properties
func (h *testHelper) assertFileWriteNode(node neta.Definition, name, path, content string) {
	if node.Type != "file.write" {
		h.t.Errorf("Expected file.write node, got '%s'", node.Type)
	}

	if node.Name != name {
		h.t.Errorf("Expected node name '%s', got '%s'", name, node.Name)
	}

	if actualPath, ok := node.Parameters["path"].(string); !ok || actualPath != path {
		h.t.Errorf("Expected path '%s', got '%v'", path, node.Parameters["path"])
	}

	if content != "" {
		if actualContent, ok := node.Parameters["content"].(string); !ok || actualContent != content {
			h.t.Errorf("Expected content '%s', got '%v'", content, node.Parameters["content"])
		}
	}
}
