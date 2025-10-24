package miso

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// configureBentoHome prompts for bento home directory configuration
func (m Model) configureBentoHome() (tea.Model, tea.Cmd) {
	// Get current bento home
	currentHome := LoadBentoHome()

	// Expand current home if it has ~
	displayHome := currentHome
	if strings.HasPrefix(displayHome, "~") {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			displayHome = filepath.Join(homeDir, displayHome[1:])
		}
	}

	// Use Huh form with file picker for directory selection
	var newHome string

	form := huh.NewForm(
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
	).WithTheme(huh.ThemeCharm())

	if err := form.Run(); err != nil {
		// User cancelled
		m.currentView = settingsView
		return m, nil
	}

	// Skip if empty (user didn't change it)
	if newHome == "" || newHome == displayHome {
		m.currentView = settingsView
		return m, nil
	}

	// Save the new home
	if err := SaveBentoHome(newHome); err != nil {
		m.logs = fmt.Sprintf("Failed to save bento home: %v", err)
		m.currentView = executionView
		return m, nil
	}

	// Show success message
	m.logs = fmt.Sprintf("✅ Bento home set to: %s\n\nNote: Existing bentos in the old location will not be moved automatically.", newHome)
	m.currentView = executionView
	return m, nil
}
