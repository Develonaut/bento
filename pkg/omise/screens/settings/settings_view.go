package settings

import (
	"github.com/charmbracelet/lipgloss"

	"bento/pkg/omise/styles"
)

// View renders the settings
func (s Settings) View() string {
	title := styles.Title.Render("Settings")

	if s.selectingDir {
		return s.renderDirectoryPickerView(title)
	}

	if s.selectingTheme {
		return s.themeForm.View()
	}

	if s.selectingSlowMo {
		return s.slowMoForm.View()
	}

	return s.renderSettingsListView(title)
}

// renderSettingsListView renders the normal settings list
func (s Settings) renderSettingsListView(title string) string {
	if s.statusMessage == "" {
		return s.list.View()
	}

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		s.list.View(),
		"",
		errorStyle.Render(s.statusMessage),
	)
}
