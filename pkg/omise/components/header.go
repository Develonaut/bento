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

// Header renders the app header bar
func Header(screen ScreenStringer, width int) string {
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
