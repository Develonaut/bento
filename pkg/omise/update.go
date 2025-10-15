package omise

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleResize(msg)
	case tea.KeyMsg:
		return m.handleKey(msg)
	default:
		return m.updateScreen(msg)
	}
}

// handleResize updates dimensions
func (m Model) handleResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height
	return m, nil
}

// handleKey processes keyboard input
func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		m.quitting = true
		return m, tea.Quit

	case "tab":
		m.screen = m.NextScreen()
		return m, nil

	case "shift+tab":
		m.screen = m.PrevScreen()
		return m, nil

	case "?", "h":
		m.screen = ScreenHelp
		return m, nil

	default:
		return m.updateScreen(msg)
	}
}

// updateScreen delegates to the current screen
func (m Model) updateScreen(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.screen {
	case ScreenBrowser:
		m.browser, cmd = m.browser.Update(msg)
	case ScreenExecutor:
		m.executor, cmd = m.executor.Update(msg)
	case ScreenPantry:
		m.pantry, cmd = m.pantry.Update(msg)
	case ScreenSettings:
		m.settings, cmd = m.settings.Update(msg)
	case ScreenHelp:
		m.help, cmd = m.help.Update(msg)
	}

	return m, cmd
}
