package omise

import (
	"testing"
)

func TestNewModel(t *testing.T) {
	m := NewModel()

	if m.screen != ScreenBrowser {
		t.Errorf("Expected initial screen to be Browser, got %v", m.screen)
	}

	if m.quitting {
		t.Error("Expected quitting to be false initially")
	}
}

func TestScreenString(t *testing.T) {
	tests := []struct {
		screen Screen
		want   string
	}{
		{ScreenBrowser, "Browser"},
		{ScreenExecutor, "Executor"},
		{ScreenPantry, "Pantry"},
		{ScreenSettings, "Settings"},
		{ScreenHelp, "Help"},
	}

	for _, tt := range tests {
		got := tt.screen.String()
		if got != tt.want {
			t.Errorf("Screen(%d).String() = %q, want %q", tt.screen, got, tt.want)
		}
	}
}

func TestNextScreen(t *testing.T) {
	m := NewModel()

	m.screen = ScreenBrowser
	if next := m.NextScreen(); next != ScreenExecutor {
		t.Errorf("NextScreen from Browser = %v, want Executor", next)
	}

	m.screen = ScreenHelp
	if next := m.NextScreen(); next != ScreenBrowser {
		t.Errorf("NextScreen from Help = %v, want Browser (wrap around)", next)
	}
}

func TestPrevScreen(t *testing.T) {
	m := NewModel()

	m.screen = ScreenExecutor
	if prev := m.PrevScreen(); prev != ScreenBrowser {
		t.Errorf("PrevScreen from Executor = %v, want Browser", prev)
	}

	m.screen = ScreenBrowser
	if prev := m.PrevScreen(); prev != ScreenHelp {
		t.Errorf("PrevScreen from Browser = %v, want Help (wrap around)", prev)
	}
}
