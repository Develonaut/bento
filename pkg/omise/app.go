// Package omise provides the Bubble Tea TUI for Bento.
// Omise (お店) means "shop" - the customer interaction point.
package omise

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// Launch starts the TUI application
func Launch() error {
	// Check if stdout is a terminal (cross-platform)
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		// Not a terminal - exit gracefully
		return nil
	}

	p := tea.NewProgram(
		NewModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	_, err := p.Run()
	return err
}
