// Package styles provides Lip Gloss styles for the Omise TUI.
package styles

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// Semantic color assignments (mutable for theme switching)
var (
	Primary   lipgloss.Color // Main theme color
	Secondary lipgloss.Color // Secondary accents
	Success   lipgloss.Color // Success states
	Error     lipgloss.Color // Error states
	Warning   lipgloss.Color // Warning states
	Text      lipgloss.Color // Primary text
	Muted     lipgloss.Color // Muted/subtle text
)

// currentVariant tracks the active theme variant
var currentVariant Variant

// Initialize with saved theme or default Nasu (purple) variant
func init() {
	currentVariant = LoadSavedTheme()
	palette := GetPalette(currentVariant)
	Primary = palette.Primary
	Secondary = palette.Secondary
	Success = palette.Success
	Error = palette.Error
	Warning = palette.Warning
	Text = palette.Text
	Muted = palette.Muted
	rebuildStyles()
}

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
	Foreground(Text)

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

// FormTheme returns a Huh theme matching our color scheme
func FormTheme() *huh.Theme {
	theme := huh.ThemeCharm()
	theme.Focused.Base = theme.Focused.Base.BorderForeground(Secondary)
	theme.Focused.Title = theme.Focused.Title.Foreground(Primary)
	theme.Focused.SelectedOption = theme.Focused.SelectedOption.Foreground(Secondary)
	theme.Focused.UnselectedOption = theme.Focused.UnselectedOption.Foreground(Muted)
	return theme
}
