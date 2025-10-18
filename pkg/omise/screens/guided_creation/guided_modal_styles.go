package guided_creation

import (
	"github.com/charmbracelet/lipgloss"
)

const guidedMaxWidth = 120

var (
	guidedIndigo = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	guidedGreen  = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
	guidedRed    = lipgloss.AdaptiveColor{Light: "#FE5F86", Dark: "#FE5F86"}
)

type GuidedStyles struct {
	Base,
	HeaderText,
	Status,
	StatusHeader,
	Highlight,
	ErrorHeaderText,
	Help,
	Breadcrumb lipgloss.Style
}

func NewGuidedStyles(lg *lipgloss.Renderer) *GuidedStyles {
	s := GuidedStyles{}
	s.Base = lg.NewStyle().
		Padding(1, 4, 0, 1)
	s.HeaderText = lg.NewStyle().
		Foreground(guidedIndigo).
		Bold(true).
		Padding(0, 1, 0, 2)
	s.Status = lg.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(guidedIndigo).
		PaddingLeft(1).
		MarginTop(1)
	s.StatusHeader = lg.NewStyle().
		Foreground(guidedGreen).
		Bold(true)
	s.Highlight = lg.NewStyle().
		Foreground(lipgloss.Color("212"))
	s.ErrorHeaderText = s.HeaderText.
		Foreground(guidedRed)
	s.Help = lg.NewStyle().
		Foreground(lipgloss.Color("240"))
	s.Breadcrumb = lg.NewStyle().
		Foreground(lipgloss.Color("246")).
		Italic(true).
		Padding(0, 1, 0, 2)
	return &s
}
