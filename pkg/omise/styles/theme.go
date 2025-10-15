// Package styles provides Lip Gloss styles for the Omise TUI.
package styles

import "github.com/charmbracelet/lipgloss"

// Color palette - Custom Bento theme
var (
	Purple = lipgloss.Color("#7359f8")
	Green  = lipgloss.Color("#66f859")
	Pink   = lipgloss.Color("#f859a8")
	Yellow = lipgloss.Color("#f8f859")
	Orange = lipgloss.Color("#f8b659")
	Red    = lipgloss.Color("#f87359")
	Cyan   = lipgloss.Color("#5cf5db")
	White  = lipgloss.Color("#e3e2e9")
	Muted  = lipgloss.Color("#565f89")

	// Semantic color assignments
	Primary   = Orange // Main theme color
	Secondary = Green  // Secondary accents
	Success   = Green  // Success states
	Error     = Red    // Error states
	Warning   = Yellow // Warning states
	Text      = White  // Primary text
)

// Title renders a section title
var Title = lipgloss.NewStyle().
	Bold(true).
	Foreground(Primary).
	MarginBottom(1)

// Header renders the app header bar
var Header = lipgloss.NewStyle().
	Bold(true).
	Padding(0, 1).
	Background(Primary).
	Foreground(White)

// Footer renders the footer bar
var Footer = lipgloss.NewStyle().
	Foreground(Muted).
	MarginTop(1)

// ErrorStyle renders error messages
var ErrorStyle = lipgloss.NewStyle().
	Foreground(Error).
	Bold(true)

// SuccessStyle renders success messages
var SuccessStyle = lipgloss.NewStyle().
	Foreground(Success).
	Bold(true)

// WarningStyle renders warning messages
var WarningStyle = lipgloss.NewStyle().
	Foreground(Warning).
	Bold(true)

// Goodbye renders the goodbye message
var Goodbye = lipgloss.NewStyle().
	Bold(true).
	Foreground(Primary).
	Padding(1)

// Subtle renders subtle text
var Subtle = lipgloss.NewStyle().
	Foreground(Muted)

// Selected renders selected items
var Selected = lipgloss.NewStyle().
	Foreground(Primary).
	Bold(true)

// Normal renders normal text
var Normal = lipgloss.NewStyle().
	Foreground(Text)

// Box creates a bordered box
var Box = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(Primary).
	Padding(1, 2)

// HelpKey renders keyboard shortcut keys
var HelpKey = lipgloss.NewStyle().
	Foreground(Secondary).
	Bold(true)

// HelpDesc renders keyboard shortcut descriptions
var HelpDesc = lipgloss.NewStyle().
	Foreground(Muted)
