package components

import (
	"github.com/charmbracelet/lipgloss"

	"bento/pkg/omise/styles"
)

// Footer renders the footer with keyboard shortcuts
func Footer(width int) string {
	shortcuts := []string{
		styles.HelpKey.Render("tab") + " " + styles.HelpDesc.Render("next"),
		styles.HelpKey.Render("shift+tab") + " " + styles.HelpDesc.Render("prev"),
		styles.HelpKey.Render("?") + " " + styles.HelpDesc.Render("help"),
		styles.HelpKey.Render("q") + " " + styles.HelpDesc.Render("quit"),
	}

	footerText := lipgloss.JoinHorizontal(lipgloss.Left, shortcuts...)
	return styles.Footer.Width(width).Render(footerText)
}
