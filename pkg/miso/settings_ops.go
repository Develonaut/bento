package miso

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// expandHomePath expands ~ to full path
func expandHomePath(path string) string {
	if strings.HasPrefix(path, "~") {
		if homeDir, err := os.UserHomeDir(); err == nil {
			return filepath.Join(homeDir, path[1:])
		}
	}
	return path
}

// configureBentoHome prompts for bento home directory configuration
func (m Model) configureBentoHome() (tea.Model, tea.Cmd) {
	currentHome := LoadBentoHome()
	// Resolve the path for the file picker to start in the right place
	resolvedHome, err := ResolvePath(currentHome)
	if err != nil {
		resolvedHome = expandHomePath(currentHome)
	}

	newHome := resolvedHome
	m.varHolders = map[string]*string{"BENTO_HOME": &newHome}

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewFilePicker().
				Title("Bento Home Directory").
				Description(fmt.Sprintf("Current: %s (Tip: Use {{GDRIVE}} for cross-platform paths)", CompressPath(resolvedHome))).
				CurrentDirectory(resolvedHome).
				DirAllowed(true).
				FileAllowed(false).
				ShowHidden(true).
				Value(&newHome),
		),
	).WithTheme(huh.ThemeCharm()).
		WithWidth(m.width).
		WithHeight(m.height)

	m.activeSettingsForm = bentoHomeForm
	m.currentView = formView
	return m, m.form.Init()
}

// buildThemeOptions creates the theme selection options for huh form.
func buildThemeOptions() []huh.Option[string] {
	return []huh.Option[string]{
		huh.NewOption(formatThemeOption(VariantNasu, "Purple", "eggplant sushi"), string(VariantNasu)),
		huh.NewOption(formatThemeOption(VariantWasabi, "Green", "wasabi"), string(VariantWasabi)),
		huh.NewOption(formatThemeOption(VariantToro, "Pink", "fatty tuna"), string(VariantToro)),
		huh.NewOption(formatThemeOption(VariantTamago, "Yellow", "egg sushi"), string(VariantTamago)),
		huh.NewOption(formatThemeOption(VariantTonkotsu, "Red", "pork bone broth"), string(VariantTonkotsu)),
		huh.NewOption(formatThemeOption(VariantSaba, "Cyan", "mackerel"), string(VariantSaba)),
		huh.NewOption(formatThemeOption(VariantIka, "White", "squid"), string(VariantIka)),
	}
}

// formatThemeOption creates a theme option with a color swatch
func formatThemeOption(variant Variant, color string, description string) string {
	palette := GetPalette(variant)
	// Create a colored square using the primary color
	swatch := lipgloss.NewStyle().
		Foreground(palette.Primary).
		Render("■")
	return fmt.Sprintf("%s  %s - %s (%s)", swatch, variant, color, description)
}

// configureTheme prompts for theme selection
func (m Model) configureTheme() (tea.Model, tea.Cmd) {
	currentTheme := LoadSavedTheme()
	selectedTheme := string(currentTheme)

	// Create value holder for form
	m.varHolders = map[string]*string{"THEME": &selectedTheme}

	// Build description with current theme colors
	descriptionText := buildThemeDescription(currentTheme)

	// Create form
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select Theme").
				Description(descriptionText).
				Options(buildThemeOptions()...).
				Value(&selectedTheme),
		),
	).WithTheme(huh.ThemeCharm()).
		WithWidth(m.width).
		WithHeight(m.height)

	// Set form type and switch to form view
	m.activeSettingsForm = themeForm
	m.currentView = formView
	return m, m.form.Init()
}

// buildThemeDescription creates a description showing current theme colors
func buildThemeDescription(variant Variant) string {
	palette := GetPalette(variant)
	return fmt.Sprintf("Current: %s\n\nColors:\n  Primary:   %s\n  Secondary: %s\n  Success:   %s\n  Error:     %s\n  Warning:   %s\n  Text:      %s\n  Muted:     %s",
		variant,
		palette.Primary,
		palette.Secondary,
		palette.Success,
		palette.Error,
		palette.Warning,
		palette.Text,
		palette.Muted,
	)
}

// getFormValue retrieves a value from form holders
func getFormValue(holders map[string]*string, key string) string {
	if holder, ok := holders[key]; ok && holder != nil {
		return *holder
	}
	return ""
}

// completeSettingsForm finishes a settings form and returns to settings view
func (m Model) completeSettingsForm() (tea.Model, tea.Cmd) {
	m.activeSettingsForm = noSettingsForm
	m.currentView = settingsView
	return m, nil
}

// completeBentoHomeForm handles bento home form completion
func (m Model) completeBentoHomeForm() (tea.Model, tea.Cmd) {
	currentHome := LoadBentoHome()
	newHome := getFormValue(m.varHolders, "BENTO_HOME")

	// Resolve current home for comparison
	resolvedCurrentHome, err := ResolvePath(currentHome)
	if err != nil {
		resolvedCurrentHome = expandHomePath(currentHome)
	}

	if newHome == "" || newHome == resolvedCurrentHome {
		return m.completeSettingsForm()
	}

	// Compress the path to use {{GDRIVE}} markers for portability
	compressedPath := CompressPath(newHome)

	if err := SaveBentoHome(compressedPath); err != nil {
		// Return to settings on error
		m.activeSettingsForm = noSettingsForm
		m.currentView = settingsView
		return m, nil
	}

	// Successfully changed bento home - return to settings
	m.activeSettingsForm = noSettingsForm
	m.currentView = settingsView
	return m, nil
}

// completeThemeForm handles theme form completion
func (m Model) completeThemeForm() (tea.Model, tea.Cmd) {
	currentTheme := LoadSavedTheme()
	selectedTheme := getFormValue(m.varHolders, "THEME")

	if selectedTheme == "" || selectedTheme == string(currentTheme) {
		return m.completeSettingsForm()
	}

	newVariant := Variant(selectedTheme)
	if err := SaveTheme(newVariant); err != nil {
		// Return to settings on error
		m.activeSettingsForm = noSettingsForm
		m.currentView = settingsView
		return m, nil
	}

	// Successfully changed theme - update model and return to settings
	m.theme = newVariant
	m.activeSettingsForm = noSettingsForm
	m.currentView = settingsView
	return m, nil
}
