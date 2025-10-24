package miso

import (
	"github.com/charmbracelet/lipgloss"
)

// viewList renders the list view
func (m Model) viewList() string {
	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(m.list.View())
}

// viewSettings renders the settings view
func (m Model) viewSettings() string {
	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(m.settingsList.View())
}

// viewSecrets renders the secrets view
func (m Model) viewSecrets() string {
	help := "\nPress 'a' to add secret • 'd' to delete • 'esc' to go back"
	content := m.secretsList.View() + help
	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(content)
}

// viewVariables renders the variables view
func (m Model) viewVariables() string {
	help := "\nPress 'a' to add variable • 'd' to delete • 'esc' to go back"
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

	// Help text
	help := lipgloss.NewStyle().
		Padding(0, 2).
		Faint(true).
		Render("↑/↓ scroll • PgUp/PgDn page • ESC back")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		viewportContent,
		help,
	)
}
