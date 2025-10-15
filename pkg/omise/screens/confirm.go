package screens

import (
	"bento/pkg/omise/styles"

	"github.com/charmbracelet/lipgloss"
)

// ConfirmDialog is a simple yes/no confirmation dialog
type ConfirmDialog struct {
	title   string
	message string
	context string // Context data (e.g., path to delete)
}

// NewConfirmDialog creates a confirmation dialog
func NewConfirmDialog(title, message, context string) *ConfirmDialog {
	return &ConfirmDialog{
		title:   title,
		message: message,
		context: context,
	}
}

// View renders the confirmation dialog
func (c *ConfirmDialog) View() string {
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Warning).
		Padding(1, 2).
		Width(50)

	title := styles.WarningStyle.Render(c.title)
	message := styles.Subtle.Render(c.message)
	prompt := styles.Subtle.Render("\nPress Y to confirm, N/Esc to cancel")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		message,
		prompt,
	)

	return box.Render(content)
}
