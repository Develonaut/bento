package guided_creation

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *GuidedModal) saveBento() tea.Cmd {
	return func() tea.Msg {
		// Generate filename from name
		filename := strings.ReplaceAll(strings.ToLower(m.definition.Name), " ", "-")

		// Save to store
		if err := m.store.Save(filename, *m.definition); err != nil {
			return GuidedCompleteMsg{
				Success:   false,
				Err:       err,
				Cancelled: false,
			}
		}

		return GuidedCompleteMsg{
			Success:    true,
			Definition: m.definition,
			Cancelled:  false,
		}
	}
}
