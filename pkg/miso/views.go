package miso

import (
	"github.com/charmbracelet/lipgloss"
)

// viewList renders the list view
func (m Model) viewList() string {
	helpStr := helpText(m.listKeys.Enter, m.listKeys.Settings, m.listKeys.Quit)
	helpStyled := lipgloss.NewStyle().
		Faint(true).
		Padding(1, 0).
		Render(helpStr)

	content := m.list.View() + "\n" + helpStyled
	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(content)
}

// viewSettings renders the settings view
func (m Model) viewSettings() string {
	helpStr := helpText(m.settingsKeys.Enter, m.settingsKeys.Back, m.settingsKeys.Quit)
	helpStyled := lipgloss.NewStyle().
		Faint(true).
		Padding(1, 0).
		Render(helpStr)

	content := m.settingsList.View() + "\n" + helpStyled
	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(content)
}

// viewSecrets renders the secrets view
func (m Model) viewSecrets() string {
	helpStr := helpText(m.secretsKeys.Add, m.secretsKeys.Delete, m.secretsKeys.Back, m.secretsKeys.Quit)
	helpStyled := lipgloss.NewStyle().
		Faint(true).
		Padding(1, 0).
		Render(helpStr)

	content := m.secretsList.View() + "\n" + helpStyled
	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(content)
}

// viewVariables renders the variables view
func (m Model) viewVariables() string {
	helpStr := helpText(m.variablesKeys.Add, m.variablesKeys.Delete, m.variablesKeys.Back, m.variablesKeys.Quit)
	helpStyled := lipgloss.NewStyle().
		Faint(true).
		Padding(1, 0).
		Render(helpStr)

	content := m.variablesList.View() + "\n" + helpStyled
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
