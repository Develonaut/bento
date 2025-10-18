package bentos

import (
	"bento/pkg/omise/styles"

	"github.com/charmbracelet/lipgloss"
)

// View renders the browser
func (b Browser) View() string {
	// PRIORITY 1: If guided modal is active, show ONLY the modal
	if b.guidedModal != nil {
		return b.guidedModal.View()
	}

	// PRIORITY 2: If confirmation dialog is active, show it
	if b.confirmDialog != nil {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			b.list.View(),
			"",
			b.confirmDialog.View(),
		)
	}

	// PRIORITY 3: If full help is showing, render that
	if b.helpView.IsFullHelpShowing() {
		return b.renderFullHelp()
	}

	// Default: Show the list
	return b.list.View()
}

// renderFullHelp renders full help view
func (b Browser) renderFullHelp() string {
	title := styles.Title.Render("Browser - Keyboard Shortcuts")
	help := b.helpView.FullHelp(b.keys)
	footer := styles.Subtle.Render("\nPress ? to toggle help")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		help,
		footer,
	)
}
