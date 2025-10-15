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

	// Should have 4 items: 1 create item + 3 actual bentos
	if len(items) != 4 {
		t.Errorf("Expected 4 items (1 create + 3 bentos), got %d", len(items))
	}

	// First item should be the create item
	if len(items) > 0 {
		firstItem, ok := items[0].(bentoItem)
		if !ok {
			t.Error("Expected first item to be bentoItem")
		} else if !firstItem.isNewItem {
			t.Error("Expected first item to be the create new bento item")
		}
	}

	// Verify remaining items are actual bentos with expected fields
	for i := 1; i < len(items); i++ {
		bi, ok := items[i].(bentoItem)
		if !ok {
			t.Errorf("Expected bentoItem at index %d, got %T", i, items[i])
			continue
		}
		if bi.isNewItem {
			t.Errorf("Expected regular bento at index %d, got create item", i)
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

func TestBrowser_CreateNewBentoItem(t *testing.T) {
	workDir := t.TempDir()

	browser, err := NewBrowser(workDir)
	if err != nil {
		t.Fatalf("NewBrowser() error = %v", err)
	}

	// The create item should be present
	selected := browser.getSelected()
	if selected == nil {
		t.Fatal("Expected a selected item")
	}

	if !selected.isNewItem {
		t.Error("First item should be the create new bento item")
	}

	// Pressing enter on create item should trigger CreateBentoMsg
	_, cmd := browser.handleRun(selected)
	if cmd == nil {
		t.Fatal("Expected command from handleRun on create item")
	}

	result := cmd()
	if _, ok := result.(CreateBentoMsg); !ok {
		t.Errorf("Expected CreateBentoMsg, got %T", result)
	}

	// Copy and delete should do nothing on create item
	_, copyCmd := browser.handleCopy(selected)
	if copyCmd != nil {
		t.Error("Expected nil command from handleCopy on create item")
	}

	_, deleteCmd := browser.handleDelete(selected)
	if deleteCmd != nil {
		t.Error("Expected nil command from handleDelete on create item")
	}
}
