package components

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"

	"bento/pkg/omise/styles"
)

// FooterModel wraps help.Model for footer rendering
type FooterModel struct {
	help   help.Model
	width  int
	global GlobalKeyMap
}

// NewFooter creates a new footer with help model
func NewFooter() FooterModel {
	h := help.New()
	h.Styles.ShortKey = styles.HelpKey
	h.Styles.ShortDesc = styles.HelpDesc
	h.Styles.ShortSeparator = styles.Subtle
	h.Styles.Ellipsis = styles.Subtle

	return FooterModel{
		help:   h,
		global: NewGlobalKeyMap(),
	}
}

// SetWidth sets the footer width
func (f FooterModel) SetWidth(width int) FooterModel {
	f.width = width
	// Account for left and right padding (2 each = 4 total)
	f.help.Width = width - 4
	return f
}

// View renders footer with contextual keys and global keys.
// If useBackKey is true, shows back key instead of settings key.
func (f FooterModel) View(contextualKeys []key.Binding, useBackKey bool) string {
	var globalKeys []key.Binding
	if useBackKey {
		globalKeys = f.global.ShortHelpWithBack()
	} else {
		globalKeys = f.global.ShortHelp()
	}

	// Render contextual keys separately from global keys with a separator
	var helpText string
	if len(contextualKeys) > 0 {
		contextualHelp := f.help.ShortHelpView(contextualKeys)
		globalHelp := f.help.ShortHelpView(globalKeys)
		separator := styles.Subtle.Render(" │ ")
		helpText = contextualHelp + separator + globalHelp
	} else {
		// No contextual keys, just show global
		helpText = f.help.ShortHelpView(globalKeys)
	}

	return styles.Footer.Render(helpText)
}
