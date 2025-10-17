package screens

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"

	"bento/pkg/omise/components"
	"bento/pkg/omise/styles"
)

// activateSetting activates the current setting (e.g., opens theme selector)
func (s Settings) activateSetting() (Settings, tea.Cmd) {
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
		return s.activateThemeForm()
	case "Slow-Mo Execution":
		return s.activateSlowMoForm()
	case "Save Directory":
		s.selectingDir = true
		return s, s.dirPicker.Init()
	}

	return s, nil
}

// activateThemeForm creates and activates the theme selection form
func (s Settings) activateThemeForm() (Settings, tea.Cmd) {
	s.selectingTheme = true
	s.selectedTheme = string(s.themeManager.GetVariant())

	options := buildThemeOptions(s.availableThemes)
	s.themeForm = components.NewFormSelect(
		"Select Theme",
		"Choose a sushi-themed color variant",
		options,
		&s.selectedTheme,
	)
	return s, s.themeForm.Init()
}

// activateSlowMoForm creates and activates the slow-mo selection form
func (s Settings) activateSlowMoForm() (Settings, tea.Cmd) {
	s.selectingSlowMo = true
	s.selectedSlowMo = formatSlowMoValue(s.config.SlowMoDelayMs)

	options := buildSlowMoOptions()
	s.slowMoForm = components.NewFormSelect(
		"Slow-Mo Execution",
		"Slow down execution to watch node progress",
		options,
		&s.selectedSlowMo,
	)
	return s, s.slowMoForm.Init()
}

// buildThemeOptions creates select options for theme variants with color styling
func buildThemeOptions(themes []styles.Variant) []components.SelectOption {
	options := make([]components.SelectOption, len(themes))
	// Create a renderer that forces color output (even in non-TTY environments like tests)
	renderer := lipgloss.NewRenderer(os.Stdout, termenv.WithProfile(termenv.TrueColor))

	for i, variant := range themes {
		// Get the palette for this variant to extract its primary color
		palette := styles.GetPalette(variant)
		// Create a styled label with the theme's primary color
		styledLabel := renderer.NewStyle().
			Foreground(palette.Primary).
			Bold(true).
			Render(string(variant))

		options[i] = components.SelectOption{
			Label: styledLabel,
			Value: string(variant),
		}
	}
	return options
}

// buildSlowMoOptions creates select options for slow-mo delays
func buildSlowMoOptions() []components.SelectOption {
	return []components.SelectOption{
		{Label: "Off", Value: "Off"},
		{Label: "250ms", Value: "250ms"},
		{Label: "500ms", Value: "500ms"},
		{Label: "1000ms", Value: "1000ms"},
		{Label: "2000ms", Value: "2000ms"},
		{Label: "4000ms", Value: "4000ms"},
		{Label: "8000ms", Value: "8000ms"},
	}
}
