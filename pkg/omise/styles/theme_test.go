package styles

import (
	"os"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestThemeInitialization_LoadsSavedTheme(t *testing.T) {
	// This test verifies the bug fix where theme wasn't loaded on init
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")

	t.Setenv("HOME", tempDir)
	defer func() {
		os.Setenv("HOME", originalHome)
	}()

	// Save a specific theme
	expectedVariant := VariantMaguro
	if err := SaveTheme(expectedVariant); err != nil {
		t.Fatalf("SaveTheme failed: %v", err)
	}

	// Simulate fresh initialization by reloading
	currentVariant = LoadSavedTheme()
	palette := GetPalette(currentVariant)

	// Apply the loaded theme
	Primary = palette.Primary
	Secondary = palette.Secondary
	Success = palette.Success
	Error = palette.Error
	Warning = palette.Warning
	Text = palette.Text
	Muted = palette.Muted
	rebuildStyles()

	// Verify colors match the saved theme
	expectedPalette := GetPalette(expectedVariant)
	if Primary != expectedPalette.Primary {
		t.Errorf("Primary color not loaded from saved theme. Expected %s, got %s", expectedPalette.Primary, Primary)
	}

	if currentVariant != expectedVariant {
		t.Errorf("currentVariant not set correctly. Expected %s, got %s", expectedVariant, currentVariant)
	}
}

func TestGlobalColors_Mutable(t *testing.T) {
	// Test that global colors can be changed (required for theme switching)
	originalPrimary := Primary

	testColor := lipgloss.Color("#123456")
	Primary = testColor

	if Primary != testColor {
		t.Error("Primary color should be mutable")
	}

	// Restore original
	Primary = originalPrimary
}

func TestStylesUpdate_AfterColorChange(t *testing.T) {
	// Test that styles use current colors after rebuild
	// This verifies the fix where styles weren't updating

	// Set a unique color
	testColor := lipgloss.Color("#ABCDEF")
	Primary = testColor

	// Rebuild styles
	rebuildStyles()

	// Check that Title style uses the new color
	if Title.GetForeground() != testColor {
		t.Error("Title style should use current Primary color after rebuild")
	}

	// Check that Selected style uses the new color
	if Selected.GetForeground() != testColor {
		t.Error("Selected style should use current Primary color after rebuild")
	}

	// Check that Goodbye style uses the new color
	if Goodbye.GetForeground() != testColor {
		t.Error("Goodbye style should use current Primary color after rebuild")
	}
}

func TestAllStyles_UseSemanticColors(t *testing.T) {
	// Verify all styles are using semantic colors, not hardcoded values
	// Set unique test colors
	Primary = lipgloss.Color("#111111")
	Secondary = lipgloss.Color("#222222")
	Success = lipgloss.Color("#333333")
	Error = lipgloss.Color("#444444")
	Warning = lipgloss.Color("#555555")
	Text = lipgloss.Color("#666666")
	Muted = lipgloss.Color("#777777")

	rebuildStyles()

	// Test Primary color usage
	if Title.GetForeground() != Primary {
		t.Error("Title should use Primary")
	}
	if Selected.GetForeground() != Primary {
		t.Error("Selected should use Primary")
	}
	if Goodbye.GetForeground() != Primary {
		t.Error("Goodbye should use Primary")
	}

	// Test Secondary color usage
	if HelpKey.GetForeground() != Secondary {
		t.Error("HelpKey should use Secondary")
	}

	// Test Muted color usage
	if Subtle.GetForeground() != Muted {
		t.Error("Subtle should use Muted")
	}
	if Footer.GetForeground() != Muted {
		t.Error("Footer should use Muted")
	}

	// Test Error color usage
	if ErrorStyle.GetForeground() != Error {
		t.Error("ErrorStyle should use Error")
	}

	// Test Success color usage
	if SuccessStyle.GetForeground() != Success {
		t.Error("SuccessStyle should use Success")
	}

	// Test Warning color usage
	if WarningStyle.GetForeground() != Warning {
		t.Error("WarningStyle should use Warning")
	}

	// Test Text color usage
	if Normal.GetForeground() != Text {
		t.Error("Normal should use Text")
	}
}

func TestCurrentVariant_SyncWithManager(t *testing.T) {
	// Test that currentVariant stays in sync when manager changes theme
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")

	t.Setenv("HOME", tempDir)
	defer func() {
		os.Setenv("HOME", originalHome)
	}()

	// Create manager and change theme
	manager := NewManager()
	newVariant := VariantWasabi

	manager.SetVariant(newVariant)

	// Verify currentVariant was updated
	if currentVariant != newVariant {
		t.Errorf("currentVariant not synced with manager. Expected %s, got %s", newVariant, currentVariant)
	}

	// Verify a new manager would use this variant
	manager2 := NewManager()
	if manager2.GetVariant() != newVariant {
		t.Errorf("New manager should use currentVariant. Expected %s, got %s", newVariant, manager2.GetVariant())
	}
}

func TestThemePersistence_AcrossSessions(t *testing.T) {
	// This test simulates the bug where theme wasn't persisting across app restarts
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")

	t.Setenv("HOME", tempDir)
	defer func() {
		os.Setenv("HOME", originalHome)
	}()

	// Session 1: Set and save theme
	variant1 := VariantSaba
	manager1 := NewManager()
	manager1.SetVariant(variant1)

	// Session 2: Fresh start (simulating app restart)
	// Reload theme from disk
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

	manager2 := NewManager()

	// Verify theme persisted
	if manager2.GetVariant() != variant1 {
		t.Errorf("Theme did not persist across sessions. Expected %s, got %s", variant1, manager2.GetVariant())
	}

	// Verify colors are correct
	expectedPalette := GetPalette(variant1)
	if Primary != expectedPalette.Primary {
		t.Error("Primary color not persisted correctly across sessions")
	}
}
