package screens

import (
	"os"
	"path/filepath"
	"testing"

	"bento/pkg/jubako"
	"bento/pkg/neta"

	tea "github.com/charmbracelet/bubbletea"
)

func TestBrowser_NewBrowser(t *testing.T) {
	workDir := t.TempDir()

	browser, err := NewBrowser(workDir)
	if err != nil {
		t.Fatalf("NewBrowser() error = %v", err)
	}

	if browser.store == nil {
		t.Error("Browser store should not be nil")
	}

	if browser.discovery == nil {
		t.Error("Browser discovery should not be nil")
	}
}

func TestBrowser_CreateNewKeyboardShortcut(t *testing.T) {
	workDir := t.TempDir()

	browser, err := NewBrowser(workDir)
	if err != nil {
		t.Fatalf("NewBrowser() error = %v", err)
	}

	// Test 'n' creates bento - this works without item selected
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	_, cmd := browser.Update(msg)

	if cmd == nil {
		t.Error("Expected command for create bento, got nil")
		return
	}

	result := cmd()
	if result == nil {
		t.Error("Expected message from command")
		return
	}

	if _, ok := result.(CreateBentoMsg); !ok {
		t.Errorf("Expected CreateBentoMsg, got %T", result)
	}
}

func TestBrowser_ConfirmationDialogCancel(t *testing.T) {
	workDir := t.TempDir()

	browser, err := NewBrowser(workDir)
	if err != nil {
		t.Fatalf("NewBrowser() error = %v", err)
	}

	// Manually set up confirmation dialog
	browser.confirmDialog = NewConfirmDialog("Test", "Test message", "/test/path")

	if browser.confirmDialog == nil {
		t.Error("Confirmation dialog should be set")
	}

	// Press 'n' should cancel
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	browser, _ = browser.Update(msg)

	if browser.confirmDialog != nil {
		t.Error("Expected confirmation dialog to be hidden after cancel")
	}
}

func TestBrowser_CopyBentoOperation(t *testing.T) {
	workDir := t.TempDir()
	store, _ := jubako.NewStore(workDir)

	// Create test bento with valid parameters
	def := neta.Definition{
		Version: "1.0",
		Type:    "http",
		Name:    "test",
		Parameters: map[string]interface{}{
			"url": "https://httpbin.org/get",
		},
	}
	_ = store.Save("test", def)

	browser, err := NewBrowser(workDir)
	if err != nil {
		t.Fatalf("NewBrowser() error = %v", err)
	}

	// Test copy operation directly
	item := &bentoItem{
		name: "test",
		path: filepath.Join(workDir, "test.bento.yaml"),
	}

	cmd := browser.copyBento(item)
	if cmd == nil {
		t.Error("Expected command for copy operation")
		return
	}

	// Execute copy command
	result := cmd()
	if result == nil {
		t.Error("Expected message from copy command")
		return
	}

	opMsg, ok := result.(BentoOperationCompleteMsg)
	if !ok {
		t.Errorf("Expected BentoOperationCompleteMsg, got %T", result)
		return
	}

	if !opMsg.Success {
		t.Errorf("Copy operation failed: %v", opMsg.Error)
	}

	// Verify copy exists
	copyPath := filepath.Join(workDir, "test-copy.bento.yaml")
	if _, err := os.Stat(copyPath); os.IsNotExist(err) {
		t.Error("Expected copy file to exist")
	}
}

func TestBrowser_LoadBentos(t *testing.T) {
	workDir := t.TempDir()
	store, _ := jubako.NewStore(workDir)

	// Create test bentos with valid parameters
	bentos := []string{"test1", "test2", "test3"}
	for _, name := range bentos {
		def := neta.Definition{
			Version: "1.0",
			Type:    "http",
			Name:    name,
			Parameters: map[string]interface{}{
				"url": "https://httpbin.org/get",
			},
		}
		if err := store.Save(name, def); err != nil {
			t.Fatalf("Failed to save bento: %v", err)
		}
	}

	items, err := loadBentos(store)
	if err != nil {
		t.Fatalf("loadBentos() error = %v", err)
	}

	// Should have 3 bentos
	if len(items) != 3 {
		t.Errorf("Expected 3 bentos, got %d", len(items))
	}

	// Verify all items are actual bentos with expected fields
	for i := 0; i < len(items); i++ {
		bi, ok := items[i].(bentoItem)
		if !ok {
			t.Errorf("Expected bentoItem at index %d, got %T", i, items[i])
			continue
		}
		if bi.version != "1.0" {
			t.Errorf("Expected version 1.0 at index %d, got %s", i, bi.version)
		}
		if bi.nodeType != "http" {
			t.Errorf("Expected type http at index %d, got %s", i, bi.nodeType)
		}
	}
}

func TestBrowser_HelpToggle(t *testing.T) {
	workDir := t.TempDir()

	browser, err := NewBrowser(workDir)
	if err != nil {
		t.Fatalf("NewBrowser() error = %v", err)
	}

	if browser.helpView.IsFullHelpShowing() {
		t.Error("Help should not be showing initially")
	}

	// Press '?' to show help
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	browser, _ = browser.Update(msg)

	if !browser.helpView.IsFullHelpShowing() {
		t.Error("Help should be showing after pressing '?'")
	}

	// Press '?' again to hide help
	browser, _ = browser.Update(msg)

	if browser.helpView.IsFullHelpShowing() {
		t.Error("Help should be hidden after pressing '?' again")
	}
}

func TestBrowser_NewBentoKeyboardShortcut(t *testing.T) {
	workDir := t.TempDir()

	browser, err := NewBrowser(workDir)
	if err != nil {
		t.Fatalf("NewBrowser() error = %v", err)
	}

	// Test 'n' key creates new bento
	_, cmd := browser.handleNew()
	if cmd == nil {
		t.Fatal("Expected command from handleNew")
	}

	result := cmd()
	if _, ok := result.(CreateBentoMsg); !ok {
		t.Errorf("Expected CreateBentoMsg, got %T", result)
	}
}

func TestBrowser_LoadsBentosWithDifferentNodeTypes(t *testing.T) {
	workDir := t.TempDir()
	store, _ := jubako.NewStore(workDir)

	// Create bentos with different node types
	testCases := []struct {
		name     string
		nodeType string
		params   map[string]interface{}
	}{
		{
			name:     "test-http",
			nodeType: "http",
			params: map[string]interface{}{
				"url":    "https://httpbin.org/get",
				"method": "GET",
			},
		},
		{
			name:     "test-jq",
			nodeType: "transform.jq",
			params: map[string]interface{}{
				"query": ".data",
			},
		},
		{
			name:     "test-file-write",
			nodeType: "file.write",
			params: map[string]interface{}{
				"path":    "/tmp/test.txt",
				"content": "test content",
			},
		},
		{
			name:     "test-sequence",
			nodeType: "group.sequence",
			params:   map[string]interface{}{},
		},
	}

	for _, tc := range testCases {
		def := neta.Definition{
			Version:    "1.0",
			Type:       tc.nodeType,
			Name:       tc.name,
			Parameters: tc.params,
		}
		if err := store.Save(tc.name, def); err != nil {
			t.Fatalf("Failed to save bento %s: %v", tc.name, err)
		}
	}

	// Load bentos
	items, err := loadBentos(store)
	if err != nil {
		t.Fatalf("loadBentos() error = %v", err)
	}

	// Should have 4 bentos
	expectedCount := len(testCases)
	if len(items) != expectedCount {
		t.Errorf("Expected %d bentos, got %d", expectedCount, len(items))
	}

	// Verify all bentos were loaded
	loadedNames := make(map[string]bool)
	for _, item := range items {
		if bi, ok := item.(bentoItem); ok {
			loadedNames[bi.name] = true
		}
	}

	for _, tc := range testCases {
		if !loadedNames[tc.name] {
			t.Errorf("Expected bento %s to be loaded", tc.name)
		}
	}

	// Verify file.write bento was properly loaded
	var fileWriteBento *bentoItem
	for _, item := range items {
		if bi, ok := item.(bentoItem); ok && bi.name == "test-file-write" {
			fileWriteBento = &bi
			break
		}
	}

	if fileWriteBento == nil {
		t.Fatal("file.write bento should be present in loaded items")
	}

	if fileWriteBento.nodeType != "file.write" {
		t.Errorf("Expected nodeType file.write, got %s", fileWriteBento.nodeType)
	}
}
