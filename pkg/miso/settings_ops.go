package miso

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
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
	displayHome := expandHomePath(currentHome)
	newHome := displayHome
	m.varHolders = map[string]*string{"BENTO_HOME": &newHome}

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewFilePicker().
				Title("Bento Home Directory").
				Description(fmt.Sprintf("Current: %s", currentHome)).
				CurrentDirectory(displayHome).
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
		huh.NewOption("Nasu - Purple (eggplant sushi)", string(VariantNasu)),
		huh.NewOption("Wasabi - Green (wasabi)", string(VariantWasabi)),
		huh.NewOption("Toro - Pink (fatty tuna)", string(VariantToro)),
		huh.NewOption("Tamago - Yellow (egg sushi)", string(VariantTamago)),
		huh.NewOption("Tonkotsu - Red (pork bone broth)", string(VariantTonkotsu)),
		huh.NewOption("Saba - Cyan (mackerel)", string(VariantSaba)),
		huh.NewOption("Ika - White (squid)", string(VariantIka)),
	}
}

// configureTheme prompts for theme selection
func (m Model) configureTheme() (tea.Model, tea.Cmd) {
	currentTheme := LoadSavedTheme()
	selectedTheme := string(currentTheme)

	// Create value holder for form
	m.varHolders = map[string]*string{"THEME": &selectedTheme}

	// Create form
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select Theme").
				Description(fmt.Sprintf("Current: %s", currentTheme)).
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
	displayHome := expandHomePath(currentHome)

	if newHome == "" || newHome == displayHome {
		return m.completeSettingsForm()
	}

	if err := SaveBentoHome(newHome); err != nil {
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
