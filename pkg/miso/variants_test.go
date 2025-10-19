// Package miso provides terminal output "seasoning" - themed styling and progress display.
//
// Tests for theme variants and color palettes.
package miso

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

// TestAllVariants verifies all 7 sushi-themed variants are defined.
func TestAllVariants(t *testing.T) {
	variants := AllVariants()

	if len(variants) != 7 {
		t.Errorf("Expected 7 variants, got %d", len(variants))
	}

	expected := []Variant{
		VariantNasu,
		VariantWasabi,
		VariantToro,
		VariantTamago,
		VariantMaguro,
		VariantSaba,
		VariantIka,
	}

	for i, v := range expected {
		if variants[i] != v {
			t.Errorf("Variant %d = %s, want %s", i, variants[i], v)
		}
	}
}

// TestGetPalette verifies each variant returns a complete palette.
func TestGetPalette(t *testing.T) {
	for _, variant := range AllVariants() {
		t.Run(string(variant), func(t *testing.T) {
			palette := GetPalette(variant)

			// Verify all palette colors are defined
			if palette.Primary == "" {
				t.Error("Primary color not defined")
			}
			if palette.Secondary == "" {
				t.Error("Secondary color not defined")
			}
			if palette.Success == "" {
				t.Error("Success color not defined")
			}
			if palette.Error == "" {
				t.Error("Error color not defined")
			}
			if palette.Warning == "" {
				t.Error("Warning color not defined")
			}
			if palette.Text == "" {
				t.Error("Text color not defined")
			}
			if palette.Muted == "" {
				t.Error("Muted color not defined")
			}
		})
	}
}

// TestGetPalette_DefaultVariant verifies invalid variant returns default (Maguro).
func TestGetPalette_DefaultVariant(t *testing.T) {
	invalid := Variant("invalid")
	palette := GetPalette(invalid)
	maguroPalette := GetPalette(VariantMaguro)

	if palette.Primary != maguroPalette.Primary {
		t.Error("Invalid variant should return Maguro palette")
	}
}

// TestVariantColors verifies specific color values for each variant.
func TestVariantColors(t *testing.T) {
	tests := []struct {
		variant Variant
		primary string // Expected hex color
		name    string // Friendly name
	}{
		{VariantNasu, "#BD93F9", "Purple (eggplant)"},
		{VariantWasabi, "#50FA7B", "Green (wasabi)"},
		{VariantToro, "#FF79C6", "Pink (fatty tuna)"},
		{VariantTamago, "#F1FA8C", "Yellow (egg)"},
		{VariantMaguro, "#f87359", "Red (tuna)"},
		{VariantSaba, "#8BE9FD", "Cyan (mackerel)"},
		{VariantIka, "#F8F8F2", "White (squid)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			palette := GetPalette(tt.variant)
			if string(palette.Primary) != tt.primary {
				t.Errorf("Primary color = %s, want %s", palette.Primary, tt.primary)
			}
		})
	}
}

// TestPaletteSemanticColors verifies semantic colors are consistent across variants.
func TestPaletteSemanticColors(t *testing.T) {
	// All variants should use same green for success
	expectedSuccess := lipgloss.Color("#50FA7B")

	for _, variant := range AllVariants() {
		palette := GetPalette(variant)
		if palette.Success != expectedSuccess {
			t.Errorf("%s: Success color = %s, want %s", variant, palette.Success, expectedSuccess)
		}
	}
}
