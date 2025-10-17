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
	allKeys := append(contextualKeys, globalKeys...)
	helpText := f.help.ShortHelpView(allKeys)
	return styles.Footer.Render(helpText)
}
