package styles

import (
	"os"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestNewManager_LoadsSavedTheme(t *testing.T) {
	// Create a temporary config directory
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")

	t.Setenv("HOME", tempDir)
	defer func() {
		os.Setenv("HOME", originalHome)
	}()

	// Save a theme before creating manager
	expectedVariant := VariantSaba
	if err := SaveTheme(expectedVariant); err != nil {
		t.Fatalf("SaveTheme failed: %v", err)
	}

	// Reset currentVariant to simulate fresh start
	currentVariant = LoadSavedTheme()

	// Create a new manager - should use the saved theme
	manager := NewManager()

	if manager.GetVariant() != expectedVariant {
		t.Errorf("Expected manager variant to be %s, got %s", expectedVariant, manager.GetVariant())
	}

	if currentVariant != expectedVariant {
		t.Errorf("Expected currentVariant to be %s, got %s", expectedVariant, currentVariant)
	}
}

func TestSetVariant_UpdatesGlobalColors(t *testing.T) {
	// Create a temporary config directory
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")

	t.Setenv("HOME", tempDir)
	defer func() {
		os.Setenv("HOME", originalHome)
	}()

	manager := NewManager()

	// Set a different variant
	newVariant := VariantMaguro
	manager.SetVariant(newVariant)

	// Check that global colors were updated
	expectedPalette := GetPalette(newVariant)
	if Primary != expectedPalette.Primary {
		t.Errorf("Expected Primary to be %s, got %s", expectedPalette.Primary, Primary)
	}
	if Secondary != expectedPalette.Secondary {
		t.Errorf("Expected Secondary to be %s, got %s", expectedPalette.Secondary, Secondary)
	}
	if Success != expectedPalette.Success {
		t.Errorf("Expected Success to be %s, got %s", expectedPalette.Success, Success)
	}
}

func TestSetVariant_UpdatesCurrentVariant(t *testing.T) {
	// Create a temporary config directory
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")

	t.Setenv("HOME", tempDir)
	defer func() {
		os.Setenv("HOME", originalHome)
	}()

	manager := NewManager()

	// Set a different variant
	newVariant := VariantTamago
	manager.SetVariant(newVariant)

	// Check that currentVariant was updated
	if currentVariant != newVariant {
		t.Errorf("Expected currentVariant to be %s, got %s", newVariant, currentVariant)
	}

	// Check that manager's variant was updated
	if manager.GetVariant() != newVariant {
		t.Errorf("Expected manager variant to be %s, got %s", newVariant, manager.GetVariant())
	}
}

func TestSetVariant_SavesTheme(t *testing.T) {
	// Create a temporary config directory
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")

	t.Setenv("HOME", tempDir)
	defer func() {
		os.Setenv("HOME", originalHome)
	}()

	manager := NewManager()

	// Set a variant
	variant := VariantIka
	manager.SetVariant(variant)

	// Verify it was saved to disk
	loaded := LoadSavedTheme()
	if loaded != variant {
		t.Errorf("Expected saved theme to be %s, got %s", variant, loaded)
	}
}

func TestApplyTheme_UpdatesAllColors(t *testing.T) {
	// Test that applyTheme updates all semantic colors
	palette := Palette{
		Primary:   lipgloss.Color("#FF0000"),
		Secondary: lipgloss.Color("#00FF00"),
		Success:   lipgloss.Color("#0000FF"),
		Error:     lipgloss.Color("#FFFF00"),
		Warning:   lipgloss.Color("#FF00FF"),
		Text:      lipgloss.Color("#00FFFF"),
		Muted:     lipgloss.Color("#888888"),
	}

	applyTheme(palette)

	if Primary != palette.Primary {
		t.Errorf("Primary not updated correctly")
	}
	if Secondary != palette.Secondary {
		t.Errorf("Secondary not updated correctly")
	}
	if Success != palette.Success {
		t.Errorf("Success not updated correctly")
	}
	if Error != palette.Error {
		t.Errorf("Error not updated correctly")
	}
	if Warning != palette.Warning {
		t.Errorf("Warning not updated correctly")
	}
	if Text != palette.Text {
		t.Errorf("Text not updated correctly")
	}
	if Muted != palette.Muted {
		t.Errorf("Muted not updated correctly")
	}
}

func TestRebuildStyles_RecreatesStyles(t *testing.T) {
	// Set specific colors
	Primary = lipgloss.Color("#AAAAAA")
	Secondary = lipgloss.Color("#BBBBBB")
	Muted = lipgloss.Color("#CCCCCC")

	// Rebuild styles
	rebuildStyles()

	// Verify that styles use the new colors
	// We can't directly compare styles, but we can render and check
	titleStyle := Title
	if titleStyle.GetForeground() != Primary {
		t.Errorf("Title style not using Primary color after rebuild")
	}

	headerStyle := Header
	if headerStyle.GetBackground() != Primary {
		t.Errorf("Header style not using Primary background after rebuild")
	}

	footerStyle := Footer
	if footerStyle.GetForeground() != Muted {
		t.Errorf("Footer style not using Muted color after rebuild")
	}
}

func TestNextVariant_CyclesThroughThemes(t *testing.T) {
	// Create a temporary config directory
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")

	t.Setenv("HOME", tempDir)
	defer func() {
		os.Setenv("HOME", originalHome)
	}()

	manager := NewManager()
	startVariant := manager.GetVariant()

	// Cycle through all variants
	variants := AllVariants()
	for i := 0; i < len(variants); i++ {
		next := manager.NextVariant()
		if next == startVariant && i < len(variants)-1 {
			t.Errorf("NextVariant returned to start too early at iteration %d", i)
		}
	}

	// After cycling through all, should be back to start
	current := manager.GetVariant()
	if current != startVariant {
		t.Errorf("After cycling through all variants, expected to be back at %s, got %s", startVariant, current)
	}
}
