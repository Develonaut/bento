package components

import (
	"bento/pkg/omise/styles"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

// Progress wraps bubbles/progress with theme-aware styling
type Progress struct {
	progress.Model
	width int
}

// NewProgress creates a themed progress bar
func NewProgress(width int) Progress {
	p := progress.New(
		progress.WithGradient(string(styles.Primary), string(styles.Secondary)),
		progress.WithWidth(width),
	)
	return Progress{
		Model: p,
		width: width,
	}
}

// RebuildStyles updates the progress bar colors with current theme
func (p Progress) RebuildStyles() Progress {
	p.Model = progress.New(
		progress.WithGradient(string(styles.Primary), string(styles.Secondary)),
		progress.WithWidth(p.width),
	)
	return p
}

// Update handles progress messages
func (p Progress) Update(msg tea.Msg) (Progress, tea.Cmd) {
	var cmd tea.Cmd
	model, cmd := p.Model.Update(msg)
	p.Model = model.(progress.Model)
	return p, cmd
}
