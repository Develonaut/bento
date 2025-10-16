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

	// Set viewport content and render it
	content := m.renderContent()
	m.viewport.SetContent(content)
	viewportView := m.viewport.View()

	contextualKeys := m.getContextualKeys()
	footer := components.Footer(m.width, contextualKeys)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		viewportView,
		"",
		footer,
	)
}

// getContextualKeys returns contextual keys from the active screen
func (m Model) getContextualKeys() []components.KeyHelp {
	switch m.screen {
	case ScreenBrowser:
		return m.browser.ContextualKeys()
	case ScreenExecutor:
		return m.executor.ContextualKeys()
	case ScreenPantry:
		return m.pantry.ContextualKeys()
	case ScreenSettings:
		return m.settings.ContextualKeys()
	case ScreenHelp:
		return m.help.ContextualKeys()
	case ScreenEditor:
		return m.editor.ContextualKeys()
	default:
		return []components.KeyHelp{}
	}
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
