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
	settingsView := s.renderSettings()
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		settingsView,
		"",
		s.helpView.RenderFooterWithBack("", s.keys),
	)
}

// renderSettings renders the settings list
func (s Settings) renderSettings() string {
	var view string
	for i, item := range s.items {
		view += s.renderSetting(i, item)
	}
	return view
}

// renderSetting renders a single setting item
func (s Settings) renderSetting(index int, item settingItem) string {
	cursor := "  "
	nameStyle := styles.Normal
	if index == s.cursor {
		cursor = "> "
		nameStyle = styles.Selected
	}

	valueStyle := styles.Subtle
	if item.valueStyle.GetBold() || item.valueStyle.GetForeground() != lipgloss.Color("") {
		valueStyle = item.valueStyle
	}

	editIndicator := ""
	if item.editable {
		editIndicator = " [↵]"
	}

	return cursor + nameStyle.Render(item.name) + "\n" +
		"  " + valueStyle.Render("  "+item.value+editIndicator) + "\n" +
		"  " + styles.Subtle.Render("  "+item.desc) + "\n\n"
}
