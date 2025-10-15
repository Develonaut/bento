package components

import (
	"bento/pkg/omise/styles"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Spinner wraps bubbles/spinner with theme-aware styling
type Spinner struct {
	spinner.Model
}

// NewSpinner creates a themed spinner
func NewSpinner() Spinner {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(styles.Primary)
	return Spinner{Model: s}
}

// RebuildStyles updates the spinner style with current theme colors
func (s Spinner) RebuildStyles() Spinner {
	s.Model.Style = lipgloss.NewStyle().Foreground(styles.Primary)
	return s
}

// Update handles spinner messages
func (s Spinner) Update(msg tea.Msg) (Spinner, tea.Cmd) {
	var cmd tea.Cmd
	s.Model, cmd = s.Model.Update(msg)
	return s, cmd
}
