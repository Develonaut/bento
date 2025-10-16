package screens

import (
	tea "github.com/charmbracelet/bubbletea"
)

// copyBento duplicates a bento file
func (b Browser) copyBento(item *bentoItem) tea.Cmd {
	return func() tea.Msg {
		def, err := b.store.Load(item.name)
		if err != nil {
			return newOperationError("copy", err)
		}

		newName := generateCopyName(item.name)
		def.Name = newName

		if err := b.store.Save(newName, def); err != nil {
			return newOperationError("copy", err)
		}

		return newOperationSuccess("copy")
	}
}

// newOperationError creates an error operation message
func newOperationError(operation string, err error) BentoOperationCompleteMsg {
	return BentoOperationCompleteMsg{
		Operation: operation,
		Success:   false,
		Error:     err,
	}
}

// newOperationSuccess creates a success operation message
func newOperationSuccess(operation string) BentoOperationCompleteMsg {
	return BentoOperationCompleteMsg{
		Operation: operation,
		Success:   true,
	}
}

// deleteBento removes a bento file
func (b Browser) deleteBento(path string) tea.Cmd {
	return func() tea.Msg {
		name := extractBentoName(path)
		if err := b.store.Delete(name); err != nil {
			return newOperationError("delete", err)
		}
		return newOperationSuccess("delete")
	}
}
