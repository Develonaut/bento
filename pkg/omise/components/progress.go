package components

import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

// Progress wraps bubbles/progress with theme-aware styling
type Progress struct {
	progress.Model
}

// NewProgress creates a themed progress bar
func NewProgress(width int) Progress {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(width),
	)
	return Progress{Model: p}
}

// Update handles progress messages
func (p Progress) Update(msg tea.Msg) (Progress, tea.Cmd) {
	var cmd tea.Cmd
	model, cmd := p.Model.Update(msg)
	p.Model = model.(progress.Model)
	return p, cmd
}
