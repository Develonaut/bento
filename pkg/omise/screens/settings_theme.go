package screens

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"bento/pkg/omise/styles"
)

// handleThemeSelection handles key input in theme selection mode
func (s Settings) handleThemeSelection(msg tea.KeyMsg) (Settings, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		s = s.moveThemeCursorUp()
	case "down", "j":
		s = s.moveThemeCursorDown()
	case "enter", " ":
		return s.selectTheme()
	case "esc", "tab", "shift+tab":
		// Exit theme selection mode
		s.selectingTheme = false
	}
	return s, nil
}

// moveThemeCursorUp moves cursor up in theme list
func (s Settings) moveThemeCursorUp() Settings {
	if s.themeCursor > 0 {
		s.themeCursor--
	}
	return s
}

// moveThemeCursorDown moves cursor down in theme list
func (s Settings) moveThemeCursorDown() Settings {
	if s.themeCursor < len(s.availableThemes)-1 {
		s.themeCursor++
	}
	return s
}

// selectTheme applies the selected theme
func (s Settings) selectTheme() (Settings, tea.Cmd) {
	selectedVariant := s.availableThemes[s.themeCursor]
	s.themeManager.SetVariant(selectedVariant)
	s.selectingTheme = false

	// Emit theme changed message (will rebuild items with new colors)
	return s, func() tea.Msg {
		return styles.ThemeChangedMsg{}
	}
}

// activateThemeSetting enters theme selection mode
func (s Settings) activateThemeSetting() Settings {
	s.selectingTheme = true
	// Set cursor to current theme
	currentVariant := s.themeManager.GetVariant()
	for i, variant := range s.availableThemes {
		if variant == currentVariant {
			s.themeCursor = i
			break
		}
	}
	return s
}

// renderThemeSelector renders the theme selection list
func (s Settings) renderThemeSelector(title string) string {
	var themeList string
	for i, variant := range s.availableThemes {
		cursor := "  "
		// Get the palette for this variant to show its primary color
		palette := styles.GetPalette(variant)
		itemStyle := lipgloss.NewStyle().Foreground(palette.Primary)

		// Highlight current selection
		if i == s.themeCursor {
			cursor = "> "
			itemStyle = itemStyle.Bold(true)
		}

		themeList += cursor + itemStyle.Render(string(variant)) + "\n"
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		styles.Subtle.Render("Choose a theme:"),
		"",
		themeList,
	)
}
