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

	// Sync tab view with current screen
	m.tabView = m.tabView.SetActiveTab(m.ScreenToTab())

	header := components.Header(m.tabView, m.width)
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
	// Add padding to content
	paddedContent := lipgloss.NewStyle().
		PaddingLeft(2).
		PaddingRight(2).
		PaddingBottom(2).
		Render(content)
	m.viewport.SetContent(paddedContent)
	return m.viewport.View()
}

// renderFooter renders the footer with contextual keys
func (m Model) renderFooter() string {
	footerModel := components.NewFooter().SetWidth(m.width)
	contextualKeys := m.getKeyBindings()
	useBackKey := m.screen == ScreenSettings || m.screen == ScreenHelp
	footer := footerModel.View(contextualKeys, useBackKey)

	// Add padding to footer
	return lipgloss.NewStyle().
		PaddingLeft(2).
		PaddingBottom(1).
		Render(footer)
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
	default:
		return ""
	}
}
