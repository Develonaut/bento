package miso

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// updateList handles list view updates
func (m Model) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Select bento
			if selected, ok := m.list.SelectedItem().(BentoItem); ok {
				m.selectedBento = selected.FilePath
				return m.runBento()
			}
		case "s":
			// Go to settings
			m.currentView = settingsView
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// updateSettings handles settings view updates
func (m Model) updateSettings(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Select setting
			if selected, ok := m.settingsList.SelectedItem().(SettingsItem); ok {
				switch selected.Action {
				case "secrets":
					return m.loadSecretsView()
				case "variables":
					return m.loadVariablesView()
				case "bentohome":
					return m.configureBentoHome()
				case "theme":
					// TODO: Show theme selection
					return m, nil
				}
			}
		case "esc":
			// Return to list
			m.currentView = listView
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.settingsList, cmd = m.settingsList.Update(msg)
	return m, cmd
}

// updateSecrets handles secrets view updates
func (m Model) updateSecrets(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			// Add new secret
			return m.addSecret()
		case "d", "x":
			// Delete selected secret
			if selected, ok := m.secretsList.SelectedItem().(SecretItem); ok {
				return m.deleteSecret(selected.Key)
			}
		case "esc":
			// Return to settings
			m.currentView = settingsView
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.secretsList, cmd = m.secretsList.Update(msg)
	return m, cmd
}

// updateVariables handles variables view updates
func (m Model) updateVariables(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			// Add new variable
			return m.addVariable()
		case "d", "x":
			// Delete selected variable
			if selected, ok := m.variablesList.SelectedItem().(VariableItem); ok {
				return m.deleteVariable(selected.Key)
			}
		case "esc":
			// Return to settings
			m.currentView = settingsView
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.variablesList, cmd = m.variablesList.Update(msg)
	return m, cmd
}

// updateForm handles form view updates
func (m Model) updateForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Update the form
	form, formCmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	// Check if form is complete
	if m.form.State == huh.StateCompleted {
		// Extract values and move to execution
		return m.startExecution()
	}

	// Check for ESC to cancel
	if msg, ok := msg.(tea.KeyMsg); ok {
		if msg.String() == "esc" {
			m.currentView = listView
			return m, nil
		}
	}

	return m, tea.Batch(cmd, formCmd)
}

// updateExecution handles execution view updates
func (m Model) updateExecution(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Return to list
			m.currentView = listView
			return m, nil
		case "up", "k":
			m.logViewport.ScrollUp(1)
			return m, nil
		case "down", "j":
			m.logViewport.ScrollDown(1)
			return m, nil
		case "pgup", "b":
			m.logViewport.HalfPageUp()
			return m, nil
		case "pgdown", "f", " ":
			m.logViewport.HalfPageDown()
			return m, nil
		}
	}

	// Update viewport for other events
	m.logViewport, cmd = m.logViewport.Update(msg)
	return m, cmd
}
