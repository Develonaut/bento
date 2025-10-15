package styles

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestAllVariants_ReturnsAllThemes(t *testing.T) {
	variants := AllVariants()

	expectedCount := 7
	if len(variants) != expectedCount {
		t.Errorf("Expected %d variants, got %d", expectedCount, len(variants))
	}

	// Check that all expected variants are present
	expected := []Variant{
		VariantNasu,
		VariantWasabi,
		VariantToro,
		VariantTamago,
		VariantMaguro,
		VariantSaba,
		VariantIka,
	}

	for _, exp := range expected {
		found := false
		for _, v := range variants {
			if v == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected variant %s not found in AllVariants()", exp)
		}
	}
}

func TestGetPalette_AllVariants(t *testing.T) {
	// Test that all variants return valid palettes
	variants := AllVariants()

	for _, variant := range variants {
		palette := GetPalette(variant)

		// Verify all colors are set (non-empty)
		if palette.Primary == lipgloss.Color("") {
			t.Errorf("Variant %s has empty Primary color", variant)
		}
		if palette.Secondary == lipgloss.Color("") {
			t.Errorf("Variant %s has empty Secondary color", variant)
		}
		if palette.Success == lipgloss.Color("") {
			t.Errorf("Variant %s has empty Success color", variant)
		}
		if palette.Error == lipgloss.Color("") {
			t.Errorf("Variant %s has empty Error color", variant)
		}
		if palette.Warning == lipgloss.Color("") {
			t.Errorf("Variant %s has empty Warning color", variant)
		}
		if palette.Text == lipgloss.Color("") {
			t.Errorf("Variant %s has empty Text color", variant)
		}
		if palette.Muted == lipgloss.Color("") {
			t.Errorf("Variant %s has empty Muted color", variant)
		}
	}
}

func TestGetPalette_NasuPurple(t *testing.T) {
	palette := GetPalette(VariantNasu)

	// Nasu should be purple-themed
	expectedPrimary := lipgloss.Color("#BD93F9")
	if palette.Primary != expectedPrimary {
		t.Errorf("Expected Nasu Primary to be %s, got %s", expectedPrimary, palette.Primary)
	}
}

func TestGetPalette_WasabiGreen(t *testing.T) {
	palette := GetPalette(VariantWasabi)

	// Wasabi should be green-themed
	expectedPrimary := lipgloss.Color("#50FA7B")
	if palette.Primary != expectedPrimary {
		t.Errorf("Expected Wasabi Primary to be %s, got %s", expectedPrimary, palette.Primary)
	}
}

func TestGetPalette_ToroPink(t *testing.T) {
	palette := GetPalette(VariantToro)

	// Toro should be pink-themed
	expectedPrimary := lipgloss.Color("#FF79C6")
	if palette.Primary != expectedPrimary {
		t.Errorf("Expected Toro Primary to be %s, got %s", expectedPrimary, palette.Primary)
	}
}

func TestGetPalette_TamagoYellow(t *testing.T) {
	palette := GetPalette(VariantTamago)

	// Tamago should be yellow-themed
	expectedPrimary := lipgloss.Color("#F1FA8C")
	if palette.Primary != expectedPrimary {
		t.Errorf("Expected Tamago Primary to be %s, got %s", expectedPrimary, palette.Primary)
	}
}

func TestGetPalette_MaguroRed(t *testing.T) {
	palette := GetPalette(VariantMaguro)

	// Maguro should be red-themed
	expectedPrimary := lipgloss.Color("#f87359")
	if palette.Primary != expectedPrimary {
		t.Errorf("Expected Maguro Primary to be %s, got %s", expectedPrimary, palette.Primary)
	}
}

func TestGetPalette_SabaCyan(t *testing.T) {
	palette := GetPalette(VariantSaba)

	// Saba should be cyan-themed
	expectedPrimary := lipgloss.Color("#8BE9FD")
	if palette.Primary != expectedPrimary {
		t.Errorf("Expected Saba Primary to be %s, got %s", expectedPrimary, palette.Primary)
	}
}

func TestGetPalette_IkaWhite(t *testing.T) {
	palette := GetPalette(VariantIka)

	// Ika should be white-themed
	expectedPrimary := lipgloss.Color("#F8F8F2")
	if palette.Primary != expectedPrimary {
		t.Errorf("Expected Ika Primary to be %s, got %s", expectedPrimary, palette.Primary)
	}
}

func TestGetPalette_InvalidVariant(t *testing.T) {
	// Test that invalid variant returns default (Maguro)
	palette := GetPalette(Variant("InvalidVariant"))

	maguroPalette := GetPalette(VariantMaguro)
	if palette.Primary != maguroPalette.Primary {
		t.Errorf("Expected invalid variant to return Maguro palette")
	}
}

func TestGetPalette_Consistency(t *testing.T) {
	// Test that calling GetPalette multiple times returns consistent results
	variant := VariantToro
	palette1 := GetPalette(variant)
	palette2 := GetPalette(variant)

	if palette1.Primary != palette2.Primary {
		t.Error("GetPalette returned inconsistent Primary color")
	}
	if palette1.Secondary != palette2.Secondary {
		t.Error("GetPalette returned inconsistent Secondary color")
	}
	if palette1.Success != palette2.Success {
		t.Error("GetPalette returned inconsistent Success color")
	}
}

func TestVariantNames_SushiThemed(t *testing.T) {
	// Verify that all variant names are sushi-themed
	variants := AllVariants()

	sushiNames := map[Variant]string{
		VariantNasu:   "Nasu",
		VariantWasabi: "Wasabi",
		VariantToro:   "Toro",
		VariantTamago: "Tamago",
		VariantMaguro: "Maguro",
		VariantSaba:   "Saba",
		VariantIka:    "Ika",
	}

	for _, variant := range variants {
		name, exists := sushiNames[variant]
		if !exists {
			t.Errorf("Variant %s is not in the expected sushi-themed names", variant)
		}
		if string(variant) != name {
			t.Errorf("Expected variant name to be %s, got %s", name, string(variant))
		}
	}
}
