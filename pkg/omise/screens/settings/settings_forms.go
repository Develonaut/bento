package settings

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"bento/pkg/omise/config"
	"bento/pkg/omise/styles"
)

// parseSlowMoValue converts display string to milliseconds
func parseSlowMoValue(value string) int {
	switch value {
	case "Off":
		return 0
	case "250ms":
		return 250
	case "500ms":
		return 500
	case "1000ms":
		return 1000
	case "2000ms":
		return 2000
	case "4000ms":
		return 4000
	case "8000ms":
		return 8000
	default:
		return 0
	}
}

// applyThemeSelection applies the selected theme from the form
func (s Settings) applyThemeSelection() Settings {
	// Get the value from the form
	selectedValue := s.themeForm.GetValue()
	variant := styles.Variant(selectedValue)
	s.themeManager.SetVariant(variant)
	s.selectingTheme = false

	items := s.buildSettings()
	s.list.SetItems(items)

	return s
}

// applySlowMoSelection applies the selected slow-mo value from the form
func (s Settings) applySlowMoSelection() Settings {
	// Get the value from the form
	selectedValue := s.slowMoForm.GetValue()
	delayMs := parseSlowMoValue(selectedValue)
	s.config.SlowMoDelayMs = delayMs

	if err := config.Save(s.config); err != nil {
		s.statusMessage = fmt.Sprintf("Warning: Failed to save config: %v", err)
	} else {
		s.statusMessage = ""
	}

	s.selectingSlowMo = false

	items := s.buildSettings()
	s.list.SetItems(items)

	return s
}

// shouldExitThemeForm checks if user pressed exit keys
func shouldExitFormMode(msg tea.Msg) bool {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		key := keyMsg.String()
		return key == "esc" || key == "tab" || key == "shift+tab"
	}
	return false
}

// handleThemeFormMode processes input when theme form is active
func (s Settings) handleThemeFormMode(msg tea.Msg) (Settings, tea.Cmd) {
	if shouldExitFormMode(msg) {
		s.selectingTheme = false
		return s, nil
	}

	var cmd tea.Cmd
	s.themeForm, cmd = s.themeForm.Update(msg)

	if s.themeForm.IsCompleted() {
		s = s.applyThemeSelection()
		return s, func() tea.Msg {
			return styles.ThemeChangedMsg{}
		}
	}

	return s, cmd
}

// handleSlowMoFormMode processes input when slow-mo form is active
func (s Settings) handleSlowMoFormMode(msg tea.Msg) (Settings, tea.Cmd) {
	if shouldExitFormMode(msg) {
		s.selectingSlowMo = false
		return s, nil
	}

	var cmd tea.Cmd
	s.slowMoForm, cmd = s.slowMoForm.Update(msg)

	if s.slowMoForm.IsCompleted() {
		s = s.applySlowMoSelection()
		return s, nil
	}

	return s, cmd
}
