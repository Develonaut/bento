package components

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"

	"bento/pkg/omise/styles"
)

// HelpView wraps Bubbles help for consistent rendering
type HelpView struct {
	help   help.Model
	global GlobalKeyMap
}

// NewHelpView creates a new help view
func NewHelpView() HelpView {
	h := help.New()
	h.Styles.ShortKey = styles.HelpKey
	h.Styles.ShortDesc = styles.HelpDesc
	h.Styles.FullKey = styles.HelpKey
	h.Styles.FullDesc = styles.HelpDesc
	h.Styles.ShortSeparator = styles.Subtle
	h.Styles.Ellipsis = styles.Subtle

	return HelpView{
		help:   h,
		global: NewGlobalKeyMap(),
	}
}

// SetWidth sets the help view width
func (h HelpView) SetWidth(width int) HelpView {
	h.help.Width = width
	return h
}

// SetShowAll sets whether to show full help
func (h HelpView) SetShowAll(show bool) HelpView {
	h.help.ShowAll = show
	return h
}

// Toggle toggles between short and full help
func (h HelpView) Toggle() HelpView {
	h.help.ShowAll = !h.help.ShowAll
	return h
}

// IsFullHelpShowing returns true if full help is showing
func (h HelpView) IsFullHelpShowing() bool {
	return h.help.ShowAll
}

// ShortHelp renders short help with screen-specific keys
func (h HelpView) ShortHelp(screenKeys interface{ ShortHelp() []key.Binding }) string {
	keys := append(screenKeys.ShortHelp(), h.global.ShortHelp()...)
	return h.help.ShortHelpView(keys)
}

// FullHelp renders full help with screen-specific keys
func (h HelpView) FullHelp(screenKeys interface{ FullHelp() [][]key.Binding }) string {
	return h.help.FullHelpView(screenKeys.FullHelp())
}

// View renders help (short or full based on state)
func (h HelpView) View(screenKeys interface {
	ShortHelp() []key.Binding
	FullHelp() [][]key.Binding
}) string {
	if h.help.ShowAll {
		return h.FullHelp(screenKeys)
	}
	return h.ShortHelp(screenKeys)
}

// RenderFooter renders a footer with optional message and help
func (h HelpView) RenderFooter(message string, screenKeys interface{ ShortHelp() []key.Binding }) string {
	helpText := h.ShortHelp(screenKeys)

	if message != "" {
		helpText = styles.SuccessStyle.Render(message) + " • " + helpText
	}

	return styles.Subtle.Render(helpText)
}

// RenderFooterWithBack renders a footer with back key instead of settings
func (h HelpView) RenderFooterWithBack(message string, screenKeys interface{ ShortHelp() []key.Binding }) string {
	keys := append(screenKeys.ShortHelp(), h.global.ShortHelpWithBack()...)
	helpText := h.help.ShortHelpView(keys)

	if message != "" {
		helpText = styles.SuccessStyle.Render(message) + " • " + helpText
	}

	return styles.Subtle.Render(helpText)
}

// RenderKeys renders a list of key bindings using help component
func (h HelpView) RenderKeys(keys []key.Binding) string {
	return h.help.ShortHelpView(keys)
}
