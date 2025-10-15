package omise

import (
	tea "github.com/charmbracelet/bubbletea"

	"bento/pkg/omise/screens"
	"bento/pkg/omise/styles"
)

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleResize(msg)
	case tea.KeyMsg:
		return m.handleKey(msg)
	case screens.WorkflowSelectedMsg:
		return m.handleWorkflowSelected(msg)
	case styles.ThemeChangedMsg:
		return m.handleThemeChanged(msg)
	default:
		return m.updateScreen(msg)
	}
}

// handleResize updates dimensions and propagates to screens
func (m Model) handleResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height

	// Pass resize to all screens so they can update their dimensions
	return m.updateScreen(msg)
}

// handleKey processes keyboard input
func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Check if Settings screen is in modal mode (picker open)
	// If so, let the screen handle tab/shift+tab instead of global handler
	if m.screen == ScreenSettings && m.settings.InModalMode() {
		return m.updateScreen(msg)
	}

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

// handleWorkflowSelected switches to executor and starts workflow
func (m Model) handleWorkflowSelected(msg screens.WorkflowSelectedMsg) (tea.Model, tea.Cmd) {
	m.screen = ScreenExecutor
	m.executor = m.executor.StartWorkflow(msg.Name, msg.Path)
	return m, m.executor.ExecuteCmd()
}

// handleThemeChanged propagates theme change to all screens
func (m Model) handleThemeChanged(msg styles.ThemeChangedMsg) (tea.Model, tea.Cmd) {
	// Update all screens with new theme
	m.browser, _ = m.browser.Update(msg)
	m.executor, _ = m.executor.Update(msg)
	m.pantry, _ = m.pantry.Update(msg)
	m.settings, _ = m.settings.Update(msg)
	m.help, _ = m.help.Update(msg)
	return m, nil
}
