package omise

import (
	"time"

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
	case screens.BentoSelectedMsg:
		return m.handleBentoSelected(msg)
	case screens.WorkflowSelectedMsg:
		return m.handleWorkflowSelected(msg)
	case screens.BentoOperationCompleteMsg:
		return m.handleBentoOperation(msg)
	case screens.StartExecutionMsg:
		return m.handleStartExecution(msg)
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

	// Update viewport size
	// Header: title (1 line) + tabs (2 lines with borders) = 3 lines
	// Footer: 3 lines (with padding)
	// Spacing: 2 lines
	// Total: 8 lines
	headerFooterHeight := 8
	m.viewport.Width = msg.Width
	m.viewport.Height = msg.Height - headerFooterHeight

	// Pass resize to all screens so they can update their dimensions
	return m.updateScreen(msg)
}

// handleKey processes keyboard input
func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Always handle quit keys globally, even when delegating to screen
	if msg.String() == "q" || msg.String() == "ctrl+c" {
		return m.handleQuit()
	}

	if m.shouldDelegateToScreen() {
		return m.updateScreen(msg)
	}

	switch msg.String() {
	case "?", "h":
		return m.handleHelpShortcut()
	case "s":
		return m.handleSettingsShortcut()
	case "tab", "shift+tab":
		return m.handleTabNavigation(msg.String())
	case "1", "2", "3", "4":
		return m.handleDirectTabAccess(msg)
	case "esc":
		return m.handleEscape(msg)
	default:
		return m.updateScreen(msg)
	}
}

// shouldDelegateToScreen checks if screen needs exclusive key handling
func (m Model) shouldDelegateToScreen() bool {
	// Settings only needs delegation when in modal mode
	if m.screen == ScreenSettings && m.settings.InModalMode() {
		return true
	}
	return false
}

// handleQuit processes quit commands
func (m Model) handleQuit() (tea.Model, tea.Cmd) {
	m.quitting = true
	return m, tea.Quit
}

// handleHelpShortcut switches to help screen
func (m Model) handleHelpShortcut() (tea.Model, tea.Cmd) {
	if tabID, ok := m.tabView.TabFromKey("4"); ok {
		m = m.SwitchToTab(tabID)
	}
	return m, nil
}

// handleSettingsShortcut switches to settings screen
func (m Model) handleSettingsShortcut() (tea.Model, tea.Cmd) {
	if tabID, ok := m.tabView.TabFromKey("3"); ok {
		m = m.SwitchToTab(tabID)
	}
	return m, nil
}

// handleTabNavigation handles tab and shift+tab
func (m Model) handleTabNavigation(key string) (tea.Model, tea.Cmd) {
	if key == "tab" {
		m.tabView = m.tabView.NextTab()
	} else {
		m.tabView = m.tabView.PrevTab()
	}
	m.screen = m.TabToScreen(m.tabView.GetActiveTab())
	return m, nil
}

// handleDirectTabAccess handles numeric tab shortcuts
func (m Model) handleDirectTabAccess(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if tabID, ok := m.tabView.TabFromKey(msg.String()); ok {
		m = m.SwitchToTab(tabID)
		return m, nil
	}
	return m.updateScreen(msg)
}

// handleEscape handles ESC key
func (m Model) handleEscape(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.screen != ScreenBrowser {
		if tabID, ok := m.tabView.TabFromKey("1"); ok {
			m = m.SwitchToTab(tabID)
		}
		return m, nil
	}
	return m.updateScreen(msg)
}

// updateScreen delegates to the current screen
func (m Model) updateScreen(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var vpCmd tea.Cmd

	// Update viewport for mouse wheel and viewport-specific keys
	m.viewport, vpCmd = m.viewport.Update(msg)

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

	return m, tea.Batch(cmd, vpCmd)
}

// handleBentoSelected switches to executor and starts bento (legacy)
func (m Model) handleBentoSelected(msg screens.BentoSelectedMsg) (tea.Model, tea.Cmd) {
	m.screen = ScreenExecutor
	m.executor = m.executor.StartBento(msg.Name, msg.Path, m.workDir)
	return m, m.executor.ExecuteCmd(m.program)
}

// handleWorkflowSelected switches to executor and queues delayed execution
func (m Model) handleWorkflowSelected(msg screens.WorkflowSelectedMsg) (tea.Model, tea.Cmd) {
	m.screen = ScreenExecutor
	m.executor = m.executor.StartBento(msg.Name, msg.Path, m.workDir)

	// Return command that sends StartExecutionMsg after 500ms delay
	return m, func() tea.Msg {
		time.Sleep(500 * time.Millisecond)
		return screens.StartExecutionMsg{
			Name:    msg.Name,
			Path:    msg.Path,
			WorkDir: m.workDir,
		}
	}
}

// handleStartExecution begins actual execution after UI transition delay
func (m Model) handleStartExecution(msg screens.StartExecutionMsg) (tea.Model, tea.Cmd) {
	return m, m.executor.ExecuteCmd(m.program)
}

// handleBentoOperation handles completion of copy/delete operations
func (m Model) handleBentoOperation(msg screens.BentoOperationCompleteMsg) (tea.Model, tea.Cmd) {
	// Delegate to browser screen for refresh
	return m.updateScreen(msg)
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
