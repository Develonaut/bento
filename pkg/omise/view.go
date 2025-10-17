package omise

import (
	"github.com/charmbracelet/bubbles/key"
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
	viewportView := m.renderViewport()
	footer := m.renderFooter()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		viewportView,
		"",
		footer,
	)
}

// renderViewport renders the viewport with content
func (m Model) renderViewport() string {
	content := m.renderContent()
	m.viewport.SetContent(content)
	return m.viewport.View()
}

// renderFooter renders the footer with contextual keys
func (m Model) renderFooter() string {
	footerModel := components.NewFooter().SetWidth(m.width)
	contextualKeys := m.getKeyBindings()
	useBackKey := m.screen == ScreenEditor || m.screen == ScreenSettings || m.screen == ScreenHelp
	return footerModel.View(contextualKeys, useBackKey)
}

// getKeyBindings returns contextual key bindings from the active screen
func (m Model) getKeyBindings() []key.Binding {
	switch m.screen {
	case ScreenBrowser:
		return m.browser.KeyBindings()
	case ScreenExecutor:
		return m.executor.KeyBindings()
	case ScreenPantry:
		return m.pantry.KeyBindings()
	case ScreenSettings:
		return m.settings.KeyBindings()
	case ScreenHelp:
		return m.help.KeyBindings()
	case ScreenEditor:
		return m.editor.KeyBindings()
	default:
		return []key.Binding{}
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
