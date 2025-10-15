package styles

import "github.com/charmbracelet/lipgloss"

// Variant represents a sushi-themed color variant
type Variant string

const (
	VariantNasu   Variant = "Nasu"   // Purple (eggplant sushi)
	VariantWasabi Variant = "Wasabi" // Green (wasabi)
	VariantToro   Variant = "Toro"   // Pink (fatty tuna)
	VariantTamago Variant = "Tamago" // Yellow (egg sushi)
	VariantMaguro Variant = "Maguro" // Red (tuna)
	VariantSaba   Variant = "Saba"   // Cyan (mackerel)
	VariantIka    Variant = "Ika"    // White (squid)
)

// AllVariants returns all available theme variants
func AllVariants() []Variant {
	return []Variant{
		VariantNasu,
		VariantWasabi,
		VariantToro,
		VariantTamago,
		VariantMaguro,
		VariantSaba,
		VariantIka,
	}
}

// Palette defines colors for a theme variant
type Palette struct {
	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Success   lipgloss.Color
	Error     lipgloss.Color
	Warning   lipgloss.Color
	Text      lipgloss.Color
	Muted     lipgloss.Color
}

// GetPalette returns the color palette for a variant
func GetPalette(v Variant) Palette {
	switch v {
	case VariantNasu:
		return nasuPalette()
	case VariantWasabi:
		return wasabiPalette()
	case VariantToro:
		return toroPalette()
	case VariantTamago:
		return tamagoPalette()
	case VariantMaguro:
		return maguroPalette()
	case VariantSaba:
		return sabaPalette()
	case VariantIka:
		return ikaPalette()
	default:
		return maguroPalette()
	}
}

// Nasu - Purple (eggplant sushi)
func nasuPalette() Palette {
	return Palette{
		Primary:   lipgloss.Color("#BD93F9"), // Purple
		Secondary: lipgloss.Color("#FF79C6"), // Pink
		Success:   lipgloss.Color("#50FA7B"), // Green
		Error:     lipgloss.Color("#f87359"), // Red
		Warning:   lipgloss.Color("#F1FA8C"), // Yellow
		Text:      lipgloss.Color("#F8F8F2"), // White
		Muted:     lipgloss.Color("#6272A4"), // Comment
	}
}

// Wasabi - Green
func wasabiPalette() Palette {
	return Palette{
		Primary:   lipgloss.Color("#50FA7B"), // Green
		Secondary: lipgloss.Color("#8BE9FD"), // Cyan
		Success:   lipgloss.Color("#50FA7B"), // Green
		Error:     lipgloss.Color("#f87359"), // Red
		Warning:   lipgloss.Color("#F1FA8C"), // Yellow
		Text:      lipgloss.Color("#F8F8F2"), // White
		Muted:     lipgloss.Color("#6272A4"), // Comment
	}
}

// Toro - Pink (fatty tuna)
func toroPalette() Palette {
	return Palette{
		Primary:   lipgloss.Color("#FF79C6"), // Pink
		Secondary: lipgloss.Color("#BD93F9"), // Purple
		Success:   lipgloss.Color("#50FA7B"), // Green
		Error:     lipgloss.Color("#f87359"), // Red
		Warning:   lipgloss.Color("#F1FA8C"), // Yellow
		Text:      lipgloss.Color("#F8F8F2"), // White
		Muted:     lipgloss.Color("#6272A4"), // Comment
	}
}

// Tamago - Yellow (egg)
func tamagoPalette() Palette {
	return Palette{
		Primary:   lipgloss.Color("#F1FA8C"), // Yellow
		Secondary: lipgloss.Color("#FFB86C"), // Orange
		Success:   lipgloss.Color("#50FA7B"), // Green
		Error:     lipgloss.Color("#f87359"), // Red
		Warning:   lipgloss.Color("#F1FA8C"), // Yellow
		Text:      lipgloss.Color("#F8F8F2"), // White
		Muted:     lipgloss.Color("#6272A4"), // Comment
	}
}

// Maguro - Red (tuna)
func maguroPalette() Palette {
	return Palette{
		Primary:   lipgloss.Color("#f87359"), // Red
		Secondary: lipgloss.Color("#FFB86C"), // Pink
		Success:   lipgloss.Color("#50FA7B"), // Green
		Error:     lipgloss.Color("#f87359"), // Red
		Warning:   lipgloss.Color("#F1FA8C"), // Yellow
		Text:      lipgloss.Color("#F8F8F2"), // White
		Muted:     lipgloss.Color("#6272A4"), // Comment
	}
}

// Saba - Cyan (mackerel)
func sabaPalette() Palette {
	return Palette{
		Primary:   lipgloss.Color("#8BE9FD"), // Cyan
		Secondary: lipgloss.Color("#BD93F9"), // Purple
		Success:   lipgloss.Color("#50FA7B"), // Green
		Error:     lipgloss.Color("#f87359"), // Red
		Warning:   lipgloss.Color("#F1FA8C"), // Yellow
		Text:      lipgloss.Color("#F8F8F2"), // White
		Muted:     lipgloss.Color("#6272A4"), // Comment
	}
}

// Ika - White (squid)
func ikaPalette() Palette {
	return Palette{
		Primary:   lipgloss.Color("#F8F8F2"), // White
		Secondary: lipgloss.Color("#BFBFBF"), // Light Gray
		Success:   lipgloss.Color("#50FA7B"), // Green
		Error:     lipgloss.Color("#f87359"), // Red
		Warning:   lipgloss.Color("#F1FA8C"), // Yellow
		Text:      lipgloss.Color("#F8F8F2"), // White
		Muted:     lipgloss.Color("#6272A4"), // Comment
	}
}
