package components

import "bento/pkg/omise/styles"

// KeyHelp represents a keyboard shortcut with its description
type KeyHelp struct {
	Key  string
	Desc string
}

// Footer renders the footer with global and contextual keyboard shortcuts
func Footer(width int, contextualKeys []KeyHelp) string {
	separator := styles.Subtle.Render(" • ")

	// Build contextual keys section
	contextualShortcuts := []string{}
	for _, key := range contextualKeys {
		contextualShortcuts = append(contextualShortcuts,
			styles.HelpKey.Render(key.Key)+" "+styles.HelpDesc.Render(key.Desc))
	}

	// Build global keys section
	globalKeys := []string{
		styles.HelpKey.Render("esc") + " " + styles.HelpDesc.Render("back"),
		styles.HelpKey.Render("s") + " " + styles.HelpDesc.Render("settings"),
		styles.HelpKey.Render("?") + " " + styles.HelpDesc.Render("help"),
		styles.HelpKey.Render("q") + " " + styles.HelpDesc.Render("quit"),
	}

	// Join sections with proper separator
	var footerText string
	if len(contextualShortcuts) > 0 {
		contextualText := joinWithSeparator(contextualShortcuts, separator)
		globalText := joinWithSeparator(globalKeys, separator)
		footerText = contextualText + " " + styles.Subtle.Render("|") + " " + globalText
	} else {
		footerText = joinWithSeparator(globalKeys, separator)
	}

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
