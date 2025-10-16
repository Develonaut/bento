package omise

import (
	tea "github.com/charmbracelet/bubbletea"

	"bento/pkg/jubako"
	"bento/pkg/omise/screens"
	"bento/pkg/omise/styles"
	"bento/pkg/pantry"
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
	case screens.EditBentoMsg:
		return m.handleEditBento(msg)
	case screens.CreateBentoMsg:
		return m.handleCreateBento(msg)
	case screens.BentoOperationCompleteMsg:
		return m.handleBentoOperation(msg)
	case screens.EditorSavedMsg:
		return m.handleEditorSaved(msg)
	case screens.EditorCancelledMsg:
		return m.handleEditorCancelled(msg)
	case screens.RunBentoFromEditorMsg:
		return m.handleRunBentoFromEditor(msg)
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

	// Check if Editor screen is in modal mode (form input)
	// If so, let the screen handle tab/shift+tab instead of global handler
	if m.screen == ScreenEditor && m.editor.InModalMode() {
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
	case ScreenEditor:
		m.editor, cmd = m.editor.Update(msg)
	}

	return m, cmd
}

// handleBentoSelected switches to executor and starts bento (legacy)
func (m Model) handleBentoSelected(msg screens.BentoSelectedMsg) (tea.Model, tea.Cmd) {
	m.screen = ScreenExecutor
	m.executor = m.executor.StartBento(msg.Name, msg.Path)
	return m, m.executor.ExecuteCmd()
}

// handleWorkflowSelected switches to executor and starts bento
func (m Model) handleWorkflowSelected(msg screens.WorkflowSelectedMsg) (tea.Model, tea.Cmd) {
	m.screen = ScreenExecutor
	m.executor = m.executor.StartBento(msg.Name, msg.Path)
	return m, m.executor.ExecuteCmd()
}

// handleEditBento switches to editor for existing bento
func (m Model) handleEditBento(msg screens.EditBentoMsg) (tea.Model, tea.Cmd) {
	store, err := jubako.NewStore(m.workDir)
	if err != nil {
		return m, nil // Stay on current screen if store creation fails
	}

	registry := pantry.New()
	editor, err := screens.NewEditorEdit(store, registry, msg.Name, msg.Path)
	if err != nil {
		return m, nil // Stay on current screen if editor creation fails
	}

	m.editor = editor
	m.screen = ScreenEditor
	return m, nil
}

// handleCreateBento switches to editor for new bento
func (m Model) handleCreateBento(msg screens.CreateBentoMsg) (tea.Model, tea.Cmd) {
	store, err := jubako.NewStore(m.workDir)
	if err != nil {
		return m, nil // Stay on current screen if store creation fails
	}

	registry := pantry.New()
	m.editor = screens.NewEditorCreate(store, registry)
	m.screen = ScreenEditor
	return m, nil
}

// handleBentoOperation handles completion of copy/delete operations
func (m Model) handleBentoOperation(msg screens.BentoOperationCompleteMsg) (tea.Model, tea.Cmd) {
	// Delegate to browser screen for refresh
	return m.updateScreen(msg)
}

// handleEditorSaved returns to browser after saving
func (m Model) handleEditorSaved(msg screens.EditorSavedMsg) (tea.Model, tea.Cmd) {
	m.screen = ScreenBrowser
	// Refresh browser list
	return m.updateScreen(screens.BentoListRefreshMsg{})
}

// handleEditorCancelled returns to browser without saving
func (m Model) handleEditorCancelled(msg screens.EditorCancelledMsg) (tea.Model, tea.Cmd) {
	m.screen = ScreenBrowser
	return m, nil
}

// handleRunBentoFromEditor runs bento from editor
func (m Model) handleRunBentoFromEditor(msg screens.RunBentoFromEditorMsg) (tea.Model, tea.Cmd) {
	// Save bento first if it has a name
	if m.editor.GetBentoName() != "" {
		store, err := jubako.NewStore(m.workDir)
		if err != nil {
			// If store creation fails, still attempt to run with in-memory definition
			m.screen = ScreenExecutor
			m.executor = m.executor.StartBento(m.editor.GetBentoName(), "")
			return m, m.executor.ExecuteCmd()
		}

		if err := store.Save(m.editor.GetBentoName(), m.editor.GetDefinition()); err != nil {
			// Save failed, but still attempt to run with in-memory definition
			m.screen = ScreenExecutor
			m.executor = m.executor.StartBento(m.editor.GetBentoName(), "")
			return m, m.executor.ExecuteCmd()
		}
	}

	// Switch to executor and run
	m.screen = ScreenExecutor
	bentoName := m.editor.GetBentoName()
	if bentoName == "" {
		bentoName = "unsaved-bento"
	}
	m.executor = m.executor.StartBento(bentoName, "")
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
