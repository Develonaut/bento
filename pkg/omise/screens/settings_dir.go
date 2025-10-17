package screens

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"bento/pkg/omise/components"
	"bento/pkg/omise/config"
	"bento/pkg/omise/styles"
)

// handleDirectorySelected processes directory selection message
func (s Settings) handleDirectorySelected(msg components.DirSelectedMsg) (Settings, tea.Cmd) {
	s.config.SaveDirectory = msg.Path
	if err := config.Save(s.config); err != nil {
		s.statusMessage = fmt.Sprintf("Warning: Failed to save config: %v", err)
	} else {
		s.statusMessage = ""
	}
	items := s.buildSettings()
	s.list.SetItems(items)
	s.selectingDir = false
	return s, nil
}

// handleDirectoryPickerMode processes input when directory picker is active
func (s Settings) handleDirectoryPickerMode(msg tea.Msg) (Settings, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "tab", "shift+tab":
			// Exit directory picker mode
			s.selectingDir = false
			return s, nil
		case "r":
			// Reset to default directory
			s.dirPicker = s.dirPicker.ResetToDefault()
			return s, nil
		}
	case styles.ThemeChangedMsg:
		s.dirPicker = s.dirPicker.RebuildStyles()
	}

	var cmd tea.Cmd
	s.dirPicker, cmd = s.dirPicker.Update(msg)
	return s, cmd
}

// renderDirectoryPickerView renders the directory selection interface
func (s Settings) renderDirectoryPickerView(title string) string {
	pickerView := s.dirPicker.View()
	currentPath := styles.Subtle.Render("Current: " + s.dirPicker.CurrentDirectory)
	help := styles.Subtle.Render("↑/↓: navigate  enter: open dir  s: select current  r: reset  esc: cancel")
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		styles.Subtle.Render("Select save directory:"),
		currentPath,
		"",
		pickerView,
		"",
		help,
	)
}

// resetDirectorySetting resets the save directory to default
func (s Settings) resetDirectorySetting() (Settings, tea.Cmd) {
	defaultCfg := config.Default()
	s.config.SaveDirectory = defaultCfg.SaveDirectory
	if err := config.Save(s.config); err != nil {
		s.statusMessage = fmt.Sprintf("Warning: Failed to save config: %v", err)
	} else {
		s.statusMessage = ""
	}
	items := s.buildSettings()
	s.list.SetItems(items)
	s.dirPicker = components.NewDirPicker(defaultCfg.SaveDirectory, defaultCfg.SaveDirectory)
	return s, nil
}
