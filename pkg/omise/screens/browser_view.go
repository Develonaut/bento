package screens

import (
	"bento/pkg/omise/styles"

	"github.com/charmbracelet/lipgloss"
)

// View renders the browser
func (b Browser) View() string {
	if b.confirmDialog != nil {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			b.list.View(),
			"",
			b.confirmDialog.View(),
		)
	}

	if b.helpView.IsFullHelpShowing() {
		return b.renderFullHelp()
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		b.list.View(),
		"",
		b.renderFooter(),
	)
}

// renderFooter shows keyboard shortcuts
func (b Browser) renderFooter() string {
	return b.helpView.RenderFooter("", b.keys)
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
