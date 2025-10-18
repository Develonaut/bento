package screens

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// handleRun runs the selected bento
func (b Browser) handleRun(item *bentoItem) (Browser, tea.Cmd) {
	return b, func() tea.Msg {
		return WorkflowSelectedMsg{
			Name: item.name,
			Path: item.path,
		}
	}
}

// handleEdit shows "Coming Soon" message (no-op for now)
func (b Browser) handleEdit(item *bentoItem) (Browser, tea.Cmd) {
	// Editor feature coming soon - do nothing for now
	return b, nil
}

// handleCopy initiates bento copy
func (b Browser) handleCopy(item *bentoItem) (Browser, tea.Cmd) {
	return b, b.copyBento(item)
}

// handleDelete shows delete confirmation
func (b Browser) handleDelete(item *bentoItem) (Browser, tea.Cmd) {
	b.confirmDialog = NewConfirmDialog(
		"Delete Bento",
		fmt.Sprintf("Are you sure you want to delete '%s'?", item.name),
		item.path,
	)
	return b, nil
}

// handleNew shows "Coming Soon" message (no-op for now)
func (b Browser) handleNew() (Browser, tea.Cmd) {
	// Editor feature coming soon - do nothing for now
	return b, nil
}
