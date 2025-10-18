package bentos

import (
	"fmt"

	"bento/pkg/omise/screens/guided_creation"
	"bento/pkg/omise/screens/shared"

	tea "github.com/charmbracelet/bubbletea"
)

// handleRun runs the selected bento
func (b Browser) handleRun(item *bentoItem) (Browser, tea.Cmd) {
	return b, func() tea.Msg {
		return shared.WorkflowSelectedMsg{
			Name: item.name,
			Path: item.path,
		}
	}
}

// handleEdit opens the guided editor for the selected bento
func (b Browser) handleEdit(item *bentoItem) (Browser, tea.Cmd) {
	// TODO: Implement guided editing flow
	return b, nil
}

// handleCopy initiates bento copy
func (b Browser) handleCopy(item *bentoItem) (Browser, tea.Cmd) {
	return b, b.copyBento(item)
}

// handleDelete shows delete confirmation
func (b Browser) handleDelete(item *bentoItem) (Browser, tea.Cmd) {
	b.confirmDialog = shared.NewConfirmDialog(
		"Delete Bento",
		fmt.Sprintf("Are you sure you want to delete '%s'?", item.name),
		item.path,
	)
	return b, nil
}

// handleNew creates a new bento with guided prompts
func (b Browser) handleNew() (Browser, tea.Cmd) {
	// Create the guided modal
	modal := guided_creation.NewGuidedModal(b.store, b.store.WorkDir(), b.width, b.height)
	b.guidedModal = modal

	// Initialize the modal
	return b, modal.Init()
}
