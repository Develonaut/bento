// Package omise provides the Bubble Tea TUI for Bento.
// Omise (お店) means "shop" - the customer interaction point.
package omise

import (
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

// Launch starts the TUI application
func Launch() error {
	// Check if stdout is a terminal (cross-platform)
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		// Not a terminal - exit gracefully
		return nil
	}

	// Get or create work directory
	workDir, err := getWorkDir()
	if err != nil {
		return err
	}

	// Initialize model with work directory
	m, err := NewModelWithWorkDir(workDir)
	if err != nil {
		return err
	}

	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	_, err = p.Run()
	return err
}

// getWorkDir returns the bento work directory
func getWorkDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	workDir := filepath.Join(home, ".bento", "bentos")
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return "", err
	}

	return workDir, nil
}
