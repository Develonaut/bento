package settings

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"bento/pkg/omise/config"
	"bento/pkg/omise/styles"
)

// resetCurrentSetting resets the current setting to its default value
func (s Settings) resetCurrentSetting() (Settings, tea.Cmd) {
	selected := s.list.SelectedItem()
	if selected == nil {
		return s, nil
	}

	item, ok := selected.(settingItem)
	if !ok || !item.editable {
		return s, nil
	}

	switch item.name {
	case "Theme":
		return s.resetThemeSetting()
	case "Slow-Mo Execution":
		return s.resetSlowMoSetting()
	case "Save Directory":
		return s.resetDirectorySetting()
	}

	return s, nil
}

// resetThemeSetting resets the theme to default (Maguro)
func (s Settings) resetThemeSetting() (Settings, tea.Cmd) {
	s.themeManager.SetVariant(styles.VariantMaguro)
	items := s.buildSettings()
	s.list.SetItems(items)
	return s, func() tea.Msg {
		return styles.ThemeChangedMsg{}
	}
}

// resetSlowMoSetting resets slow-mo to off (0ms)
func (s Settings) resetSlowMoSetting() (Settings, tea.Cmd) {
	s.config.SlowMoDelayMs = 0
	if err := config.Save(s.config); err != nil {
		s.statusMessage = fmt.Sprintf("Warning: Failed to save config: %v", err)
	} else {
		s.statusMessage = ""
	}
	items := s.buildSettings()
	s.list.SetItems(items)
	return s, nil
}
