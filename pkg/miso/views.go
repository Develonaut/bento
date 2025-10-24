package miso

import (
	"github.com/charmbracelet/lipgloss"
)

// viewList renders the list view
func (m Model) viewList() string {
	help := "\n" + helpText(m.listKeys.Enter, m.listKeys.Settings, m.listKeys.Quit)
	content := m.list.View() + help
	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(content)
}

// viewSettings renders the settings view
func (m Model) viewSettings() string {
	help := "\n" + helpText(m.settingsKeys.Enter, m.settingsKeys.Back, m.settingsKeys.Quit)
	content := m.settingsList.View() + help
	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(content)
}

// viewSecrets renders the secrets view
func (m Model) viewSecrets() string {
	help := "\n" + helpText(m.secretsKeys.Add, m.secretsKeys.Delete, m.secretsKeys.Back, m.secretsKeys.Quit)
	content := m.secretsList.View() + help
	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(content)
}

// viewVariables renders the variables view
func (m Model) viewVariables() string {
	help := "\n" + helpText(m.variablesKeys.Add, m.variablesKeys.Delete, m.variablesKeys.Back, m.variablesKeys.Quit)
	content := m.variablesList.View() + help
	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(content)
}

// viewForm renders the form view
func (m Model) viewForm() string {
	if m.form == nil {
		return "Loading form..."
	}
	return m.form.View()
}

// viewExecution renders the execution view
func (m Model) viewExecution() string {
	palette := GetPalette(m.theme)

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(palette.Primary).
		Bold(true).
		Padding(1, 2)
	title := titleStyle.Render("🍱 Execution")

	// Bordered viewport container
	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(palette.Primary).
		Padding(1, 2)

	viewportContent := borderStyle.Render(m.logViewport.View())

	// Help text from key bindings
	helpStr := helpText(m.executionKeys.ScrollUp, m.executionKeys.ScrollDown,
		m.executionKeys.PageUp, m.executionKeys.PageDown,
		m.executionKeys.Back, m.executionKeys.Quit)
	help := lipgloss.NewStyle().
		Padding(0, 2).
		Faint(true).
		Render(helpStr)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		viewportContent,
		help,
	)
}
