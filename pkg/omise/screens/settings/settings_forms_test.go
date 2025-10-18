package settings

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"bento/pkg/omise/components"
	"bento/pkg/omise/styles"
)

// TestThemeFormActivation tests that activating theme form works
func TestThemeFormActivation(t *testing.T) {
	s := NewSettings()

	// Activate theme form
	s, _ = s.activateThemeForm()

	if !s.selectingTheme {
		t.Error("Expected selectingTheme to be true after activation")
	}

	if s.themeForm.GetForm() == nil {
		t.Error("Expected theme form to be initialized")
	}

	if s.selectedTheme != string(s.themeManager.GetVariant()) {
		t.Errorf("Expected selectedTheme to match current variant, got %s, want %s",
			s.selectedTheme, s.themeManager.GetVariant())
	}
}

// TestSlowMoFormActivation tests that activating slow-mo form works
func TestSlowMoFormActivation(t *testing.T) {
	s := NewSettings()

	// Activate slow-mo form
	s, _ = s.activateSlowMoForm()

	if !s.selectingSlowMo {
		t.Error("Expected selectingSlowMo to be true after activation")
	}

	if s.slowMoForm.GetForm() == nil {
		t.Error("Expected slow-mo form to be initialized")
	}

	expectedValue := formatSlowMoValue(s.config.SlowMoDelayMs)
	if s.selectedSlowMo != expectedValue {
		t.Errorf("Expected selectedSlowMo to be %s, got %s", expectedValue, s.selectedSlowMo)
	}
}

// TestThemeFormCompletion tests theme selection and application
func TestThemeFormCompletion(t *testing.T) {
	s := NewSettings()
	originalVariant := s.themeManager.GetVariant()

	// Activate theme form
	s, _ = s.activateThemeForm()

	// Select a different theme
	newTheme := styles.VariantMaguro
	if newTheme == originalVariant {
		newTheme = styles.VariantToro // Pick a different one
	}

	// Simulate form completion by setting the value
	s.selectedTheme = string(newTheme)

	// Create a mock form that's completed
	s.themeForm = components.NewFormSelect(
		"Select Theme",
		"Choose a sushi-themed color variant",
		buildThemeOptions(s.availableThemes),
		&s.selectedTheme,
	)

	// Apply the selection
	s = s.applyThemeSelection()

	if s.selectingTheme {
		t.Error("Expected selectingTheme to be false after applying")
	}

	if s.themeManager.GetVariant() != newTheme {
		t.Errorf("Expected theme to be %s, got %s", newTheme, s.themeManager.GetVariant())
	}
}

// TestSlowMoFormCompletion tests slow-mo selection and application
func TestSlowMoFormCompletion(t *testing.T) {
	s := NewSettings()
	s.config.SlowMoDelayMs = 0

	// Activate slow-mo form
	s, _ = s.activateSlowMoForm()

	// Select a delay value
	s.selectedSlowMo = "1000ms"

	// Create a mock form
	s.slowMoForm = components.NewFormSelect(
		"Slow-Mo Execution",
		"Slow down execution to watch node progress",
		buildSlowMoOptions(),
		&s.selectedSlowMo,
	)

	// Apply the selection
	s = s.applySlowMoSelection()

	if s.selectingSlowMo {
		t.Error("Expected selectingSlowMo to be false after applying")
	}

	if s.config.SlowMoDelayMs != 1000 {
		t.Errorf("Expected delay to be 1000, got %d", s.config.SlowMoDelayMs)
	}
}

// TestParseSlowMoValue tests parsing of slow-mo display values
func TestParseSlowMoValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Off", "Off", 0},
		{"250ms", "250ms", 250},
		{"500ms", "500ms", 500},
		{"1000ms", "1000ms", 1000},
		{"2000ms", "2000ms", 2000},
		{"4000ms", "4000ms", 4000},
		{"8000ms", "8000ms", 8000},
		{"Unknown", "invalid", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseSlowMoValue(tt.input)
			if result != tt.expected {
				t.Errorf("parseSlowMoValue(%s) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

// TestFormatSlowMoValue tests formatting of slow-mo values for display
func TestFormatSlowMoValue(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected string
	}{
		{"Zero is Off", 0, "Off"},
		{"250", 250, "250ms"},
		{"500", 500, "500ms"},
		{"1000", 1000, "1000ms"},
		{"2000", 2000, "2000ms"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatSlowMoValue(tt.input)
			if result != tt.expected {
				t.Errorf("formatSlowMoValue(%d) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

// TestShouldExitFormMode tests form exit detection
func TestShouldExitFormMode(t *testing.T) {
	tests := []struct {
		name     string
		msg      tea.Msg
		expected bool
	}{
		{"Esc exits", tea.KeyMsg{Type: tea.KeyEsc}, true},
		{"Tab exits", tea.KeyMsg{Type: tea.KeyTab}, true},
		{"Shift+Tab exits", tea.KeyMsg{Type: tea.KeyShiftTab}, true},
		{"Enter does not exit", tea.KeyMsg{Type: tea.KeyEnter}, false},
		{"Space does not exit", tea.KeyMsg{Type: tea.KeySpace}, false},
		{"Letter does not exit", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldExitFormMode(tt.msg)
			if result != tt.expected {
				t.Errorf("shouldExitFormMode(%v) = %v, want %v", tt.msg, result, tt.expected)
			}
		})
	}
}

// TestHandleThemeFormMode tests theme form message handling
func TestHandleThemeFormMode(t *testing.T) {
	s := NewSettings()
	s, _ = s.activateThemeForm()

	// Test escape exits form mode
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	s, _ = s.handleThemeFormMode(escMsg)

	if s.selectingTheme {
		t.Error("Expected Esc to exit theme form mode")
	}
}

// TestHandleSlowMoFormMode tests slow-mo form message handling
func TestHandleSlowMoFormMode(t *testing.T) {
	s := NewSettings()
	s, _ = s.activateSlowMoForm()

	// Test escape exits form mode
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	s, _ = s.handleSlowMoFormMode(escMsg)

	if s.selectingSlowMo {
		t.Error("Expected Esc to exit slow-mo form mode")
	}
}

// TestThemeFormEscapeExits tests that escape key exits theme form
func TestThemeFormEscapeExits(t *testing.T) {
	s := NewSettings()

	// Activate theme form
	s, _ = s.activateThemeForm()
	if !s.selectingTheme {
		t.Fatal("Failed to activate theme form")
	}

	// Press escape
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	s, _ = s.handleThemeFormMode(escMsg)

	if s.selectingTheme {
		t.Error("Expected escape to exit theme form")
	}
}

// TestSlowMoFormEscapeExits tests that escape key exits slow-mo form
func TestSlowMoFormEscapeExits(t *testing.T) {
	s := NewSettings()

	// Activate slow-mo form
	s, _ = s.activateSlowMoForm()
	if !s.selectingSlowMo {
		t.Fatal("Failed to activate slow-mo form")
	}

	// Press escape
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	s, _ = s.handleSlowMoFormMode(escMsg)

	if s.selectingSlowMo {
		t.Error("Expected escape to exit slow-mo form")
	}
}

// TestThemeFormTabExits tests that tab key exits theme form
func TestThemeFormTabExits(t *testing.T) {
	s := NewSettings()

	// Activate theme form
	s, _ = s.activateThemeForm()
	if !s.selectingTheme {
		t.Fatal("Failed to activate theme form")
	}

	// Press tab
	tabMsg := tea.KeyMsg{Type: tea.KeyTab}
	s, _ = s.handleThemeFormMode(tabMsg)

	if s.selectingTheme {
		t.Error("Expected tab to exit theme form")
	}
}

// TestSlowMoFormTabExits tests that tab key exits slow-mo form
func TestSlowMoFormTabExits(t *testing.T) {
	s := NewSettings()

	// Activate slow-mo form
	s, _ = s.activateSlowMoForm()
	if !s.selectingSlowMo {
		t.Fatal("Failed to activate slow-mo form")
	}

	// Press tab
	tabMsg := tea.KeyMsg{Type: tea.KeyTab}
	s, _ = s.handleSlowMoFormMode(tabMsg)

	if s.selectingSlowMo {
		t.Error("Expected tab to exit slow-mo form")
	}
}

// TestBuildSlowMoOptions tests that slow-mo options are built correctly
func TestBuildSlowMoOptions(t *testing.T) {
	options := buildSlowMoOptions()

	expectedCount := 7 // Off, 250ms, 500ms, 1000ms, 2000ms, 4000ms, 8000ms
	if len(options) != expectedCount {
		t.Errorf("Expected %d options, got %d", expectedCount, len(options))
	}

	// Check first option is "Off"
	if options[0].Label != "Off" || options[0].Value != "Off" {
		t.Errorf("Expected first option to be Off, got %s", options[0].Label)
	}

	// Check last option is "8000ms"
	lastIdx := len(options) - 1
	if options[lastIdx].Label != "8000ms" || options[lastIdx].Value != "8000ms" {
		t.Errorf("Expected last option to be 8000ms, got %s", options[lastIdx].Label)
	}
}

// TestBuildThemeOptions tests that theme options are built correctly
func TestBuildThemeOptions(t *testing.T) {
	themes := styles.AllVariants()
	options := buildThemeOptions(themes)

	if len(options) != len(themes) {
		t.Errorf("Expected %d options, got %d", len(themes), len(options))
	}

	// Check each theme is represented
	for i, theme := range themes {
		if options[i].Label != string(theme) {
			t.Errorf("Expected option %d label to be %s, got %s", i, theme, options[i].Label)
		}
		if options[i].Value != string(theme) {
			t.Errorf("Expected option %d value to be %s, got %s", i, theme, options[i].Value)
		}
	}
}

// TestApplyThemeSelectionUpdatesConfig tests that theme selection updates config
func TestApplyThemeSelectionUpdatesConfig(t *testing.T) {
	s := NewSettings()
	originalTheme := s.themeManager.GetVariant()

	// Activate theme form
	s, _ = s.activateThemeForm()

	// Select a different theme by updating the value that the form's pointer references
	newTheme := styles.VariantToro
	if newTheme == originalTheme {
		newTheme = styles.VariantMaguro
	}
	s.selectedTheme = string(newTheme)

	// Recreate form with the updated value pointer to simulate form completion
	s.themeForm = components.NewFormSelect(
		"Select Theme",
		"Choose a sushi-themed color variant",
		buildThemeOptions(s.availableThemes),
		&s.selectedTheme,
	)

	// Apply selection
	s = s.applyThemeSelection()

	// Verify theme changed
	if s.themeManager.GetVariant() != newTheme {
		t.Errorf("Expected theme to be %s, got %s", newTheme, s.themeManager.GetVariant())
	}

	// Verify settings list is updated
	items := s.buildSettings()
	themeItem := items[0].(settingItem)
	if themeItem.value != string(newTheme) {
		t.Errorf("Expected theme item value to be %s, got %s", newTheme, themeItem.value)
	}
}

// TestApplySlowMoSelectionUpdatesConfig tests that slow-mo selection updates config
func TestApplySlowMoSelectionUpdatesConfig(t *testing.T) {
	s := NewSettings()
	s.config.SlowMoDelayMs = 0

	// Activate slow-mo form
	s, _ = s.activateSlowMoForm()

	// Set a slow-mo value by updating the value that the form's pointer references
	s.selectedSlowMo = "2000ms"

	// Recreate form with the updated value pointer to simulate form completion
	s.slowMoForm = components.NewFormSelect(
		"Slow-Mo Execution",
		"Slow down execution to watch node progress",
		buildSlowMoOptions(),
		&s.selectedSlowMo,
	)

	// Apply selection
	s = s.applySlowMoSelection()

	// Verify delay changed
	if s.config.SlowMoDelayMs != 2000 {
		t.Errorf("Expected delay to be 2000, got %d", s.config.SlowMoDelayMs)
	}

	// Verify settings list is updated
	items := s.buildSettings()
	slowMoItem := items[1].(settingItem)
	if slowMoItem.value != "2000ms" {
		t.Errorf("Expected slow-mo item value to be 2000ms, got %s", slowMoItem.value)
	}
}

// TestSlowMoConfigPersistence tests that slow-mo setting is saved to config
func TestSlowMoConfigPersistence(t *testing.T) {
	// Use NewSettings to get properly initialized settings
	s := NewSettings()
	s.config.SlowMoDelayMs = 0

	// Apply a slow-mo value
	s, _ = s.activateSlowMoForm()
	s.selectedSlowMo = "500ms"

	// Recreate form with the updated value pointer
	s.slowMoForm = components.NewFormSelect(
		"Slow-Mo Execution",
		"Slow down execution to watch node progress",
		buildSlowMoOptions(),
		&s.selectedSlowMo,
	)

	s = s.applySlowMoSelection()

	// Verify the config was updated
	if s.config.SlowMoDelayMs != 500 {
		t.Errorf("Expected config delay to be 500, got %d", s.config.SlowMoDelayMs)
	}
}
