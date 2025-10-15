package omise

import (
	"github.com/charmbracelet/lipgloss"

	"bento/pkg/omise/components"
	"bento/pkg/omise/styles"
)

// View renders the TUI
func (m Model) View() string {
	if m.quitting {
		return styles.Goodbye.Render("Thanks for using Bento! 🍱\n")
	}

	header := components.Header(m.screen, m.width)
	content := m.renderContent()
	footer := components.Footer(m.width)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		content,
		footer,
	)
}

// renderContent renders the active screen
func (m Model) renderContent() string {
	switch m.screen {
	case ScreenBrowser:
		return m.browser.View()
	case ScreenExecutor:
		return m.executor.View()
	case ScreenPantry:
		return m.pantry.View()
	case ScreenSettings:
		return m.settings.View()
	case ScreenHelp:
		return m.help.View()
	case ScreenEditor:
		return m.editor.View()
	default:
		return ""
	}
}
