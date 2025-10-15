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

	if len(items) != 3 {
		t.Errorf("Expected 3 bentos, got %d", len(items))
	}

	// Verify each item has expected fields
	for _, item := range items {
		bi, ok := item.(bentoItem)
		if !ok {
			t.Errorf("Expected bentoItem, got %T", item)
			continue
		}
		if bi.version != "1.0" {
			t.Errorf("Expected version 1.0, got %s", bi.version)
		}
		if bi.nodeType != "http" {
			t.Errorf("Expected type http, got %s", bi.nodeType)
		}
	}
}

func TestBrowser_HelpToggle(t *testing.T) {
	workDir := t.TempDir()

	browser, err := NewBrowser(workDir)
	if err != nil {
		t.Fatalf("NewBrowser() error = %v", err)
	}

	if browser.showingHelp {
		t.Error("Help should not be showing initially")
	}

	// Press '?' to show help
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	browser, _ = browser.Update(msg)

	if !browser.showingHelp {
		t.Error("Help should be showing after pressing '?'")
	}

	// Press '?' again to hide help
	browser, _ = browser.Update(msg)

	if browser.showingHelp {
		t.Error("Help should be hidden after pressing '?' again")
	}
}
