package screens

import (
	"bento/pkg/omise/styles"

	"github.com/charmbracelet/lipgloss"
)

// View renders the settings
func (s Settings) View() string {
	title := styles.Title.Render("Settings")

	if s.selectingDir {
		return s.renderDirectoryPickerView(title)
	}

	if s.selectingTheme {
		return s.renderThemeSelector(title)
	}

	return s.renderSettingsListView(title)
}

// renderSettingsListView renders the normal settings list
func (s Settings) renderSettingsListView(title string) string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		s.list.View(),
		"",
		s.helpView.RenderFooterWithBack("", s.keys),
	)
}
