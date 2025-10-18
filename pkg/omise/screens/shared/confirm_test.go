package shared

import (
	"strings"
	"testing"
)

func TestConfirmDialog_NewConfirmDialog(t *testing.T) {
	title := "Delete Item"
	message := "Are you sure?"
	context := "/path/to/item"

	dialog := NewConfirmDialog(title, message, context)

	if dialog.title != title {
		t.Errorf("Expected title %q, got %q", title, dialog.title)
	}

	if dialog.message != message {
		t.Errorf("Expected message %q, got %q", message, dialog.message)
	}

	if dialog.Context != context {
		t.Errorf("Expected context %q, got %q", context, dialog.Context)
	}
}

func TestConfirmDialog_View(t *testing.T) {
	dialog := NewConfirmDialog("Test Title", "Test Message", "context")

	view := dialog.View()

	if view == "" {
		t.Error("View should not be empty")
	}

	// Check that title and message appear in view
	if !strings.Contains(view, "Test Title") {
		t.Error("View should contain title")
	}

	if !strings.Contains(view, "Test Message") {
		t.Error("View should contain message")
	}

	// Check for confirmation prompt
	if !strings.Contains(view, "Press Y") {
		t.Error("View should contain confirmation prompt")
	}
}

func TestConfirmDialog_ViewFormatting(t *testing.T) {
	dialog := NewConfirmDialog("Delete Bento", "Are you sure you want to delete 'test'?", "/path/test.bento.yaml")

	view := dialog.View()

	// Verify both title and message are present
	if !strings.Contains(view, "Delete Bento") {
		t.Error("View should contain title 'Delete Bento'")
	}

	if !strings.Contains(view, "Are you sure") {
		t.Error("View should contain confirmation message")
	}

	if !strings.Contains(view, "test") {
		t.Error("View should contain bento name")
	}
}
