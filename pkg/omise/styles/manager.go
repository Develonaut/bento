package styles

import "github.com/charmbracelet/lipgloss"

// Manager manages the active theme variant
type Manager struct {
	variant Variant
	palette Palette
}

// NewManager creates a theme manager with the current variant
func NewManager() *Manager {
	// Use the variant already loaded in init()
	palette := GetPalette(currentVariant)

	return &Manager{
		variant: currentVariant,
		palette: palette,
	}
}

// SetVariant changes the active theme variant and saves it
func (m *Manager) SetVariant(v Variant) {
	m.variant = v
	m.palette = GetPalette(v)
	currentVariant = v // Update global current variant
	applyTheme(m.palette)

	// Save theme preference (ignore errors)
	_ = SaveTheme(v)
}

// GetVariant returns the current variant
func (m *Manager) GetVariant() Variant {
	return m.variant
}

// CurrentVariant returns the global current variant
func CurrentVariant() Variant {
	return currentVariant
}

// NextVariant cycles to the next theme variant
func (m *Manager) NextVariant() Variant {
	variants := AllVariants()
	for i, v := range variants {
		if v == m.variant {
			next := variants[(i+1)%len(variants)]
			m.SetVariant(next)
			return next
		}
	}
	return m.variant
}

// applyTheme updates global styles with palette
func applyTheme(p Palette) {
	// Update semantic colors
	Primary = p.Primary
	Secondary = p.Secondary
	Success = p.Success
	Error = p.Error
	Warning = p.Warning
	Text = p.Text
	Muted = p.Muted

	// Rebuild styles with new colors
	rebuildStyles()
}

// rebuildStyles recreates style objects with new colors
func rebuildStyles() {
	rebuildTextStyles()
	rebuildLayoutStyles()
	rebuildStateStyles()
	rebuildHelpStyles()
}

// rebuildTextStyles updates text-related styles
func rebuildTextStyles() {
	Title = lipgloss.NewStyle().Bold(true).Foreground(Primary).MarginBottom(1)
	Normal = lipgloss.NewStyle().Foreground(Text)
	Subtle = lipgloss.NewStyle().Foreground(Muted)
	Selected = lipgloss.NewStyle().Foreground(Primary).Bold(true)
}

// rebuildLayoutStyles updates layout-related styles
func rebuildLayoutStyles() {
	Header = lipgloss.NewStyle().Bold(true).Padding(0, 1).Background(Primary).Foreground(Text)
	Footer = lipgloss.NewStyle().Foreground(Muted).MarginTop(1)
	Box = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(Primary).Padding(1, 2)
}

// rebuildStateStyles updates state-related styles
func rebuildStateStyles() {
	ErrorStyle = lipgloss.NewStyle().Foreground(Error).Bold(true)
	SuccessStyle = lipgloss.NewStyle().Foreground(Success).Bold(true)
	WarningStyle = lipgloss.NewStyle().Foreground(Warning).Bold(true)
	Goodbye = lipgloss.NewStyle().Bold(true).Foreground(Primary).Padding(1)
}

// rebuildHelpStyles updates help-related styles
func rebuildHelpStyles() {
	HelpKey = lipgloss.NewStyle().Foreground(Secondary).Bold(true)
	HelpDesc = lipgloss.NewStyle().Foreground(Muted)
}
