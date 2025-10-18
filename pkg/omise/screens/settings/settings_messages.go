package settings

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// routeModalModeUpdate routes updates when in modal mode
func (s Settings) routeModalModeUpdate(msg tea.Msg) (Settings, tea.Cmd, bool) {
	if s.selectingDir {
		newS, cmd := s.handleDirectoryPickerMode(msg)
		return newS, cmd, true
	}

	if s.selectingTheme {
		newS, cmd := s.handleThemeFormMode(msg)
		return newS, cmd, true
	}

	if s.selectingSlowMo {
		newS, cmd := s.handleSlowMoFormMode(msg)
		return newS, cmd, true
	}

	return s, nil, false
}

// handleWindowResize processes window resize messages
func (s Settings) handleWindowResize(msg tea.WindowSizeMsg) (Settings, tea.Cmd) {
	h, v := lipgloss.NewStyle().Margin(2, 2).GetFrameSize()
	s.list.SetSize(msg.Width-h, msg.Height-v-4)
	return s, nil
}
