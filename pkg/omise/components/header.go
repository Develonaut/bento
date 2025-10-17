// Package components provides reusable TUI components.
package components

import (
	"github.com/charmbracelet/lipgloss"

	"bento/pkg/omise/styles"
)

// ScreenStringer defines the String method for screen types
type ScreenStringer interface {
	String() string
}

// Header renders the app header bar with tabs
func Header(tabView TabView, width int) string {
	title := styles.Header.Width(width).Render("🍱 Bento v0.1.0")
	tabBar := tabView.SetWidth(width).View()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		tabBar,
	)
}

// HeaderLegacy renders the old-style header (for compatibility)
func HeaderLegacy(screen ScreenStringer, width int) string {
	title := "🍱 Bento"
	screenName := screen.String()

	headerText := lipgloss.JoinHorizontal(
		lipgloss.Left,
		title,
		" | ",
		screenName,
	)

	return styles.Header.Width(width).Render(headerText)
}
