package components

import "bento/pkg/omise/styles"

// Footer renders the footer with keyboard shortcuts
func Footer(width int) string {
	separator := styles.Subtle.Render(" • ")
	shortcuts := []string{
		styles.HelpKey.Render("tab") + " " + styles.HelpDesc.Render("next"),
		styles.HelpKey.Render("shift+tab") + " " + styles.HelpDesc.Render("prev"),
		styles.HelpKey.Render("?") + " " + styles.HelpDesc.Render("help"),
		styles.HelpKey.Render("q") + " " + styles.HelpDesc.Render("quit"),
	}

	footerText := joinWithSeparator(shortcuts, separator)
	return styles.Footer.Width(width).Render(footerText)
}

// joinWithSeparator joins strings with a separator between them
func joinWithSeparator(items []string, sep string) string {
	if len(items) == 0 {
		return ""
	}
	result := items[0]
	for _, item := range items[1:] {
		result += sep + item
	}
	return result
}
