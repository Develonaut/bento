package screens

import (
	"bento/pkg/omise/styles"

	"github.com/charmbracelet/lipgloss"
)

// View renders the browser
func (b Browser) View() string {
	// Show action menu if active
	if b.actionMenu != nil {
		return b.actionMenu.form.View()
	}

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

	// Show available actions for selected item
	actionsHint := b.renderActionsHint()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		b.list.View(),
		actionsHint,
	)
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

// renderActionsHint shows available keyboard shortcuts for selected item
func (b Browser) renderActionsHint() string {
	selected := b.getSelected()
	if selected == nil || selected.isNewItem {
		return ""
	}

	// Use help component to render action keys consistently
	actionKeys := b.keys.ActionHelp()
	hint := b.helpView.RenderKeys(actionKeys)

	return styles.Subtle.Render(hint)
}
