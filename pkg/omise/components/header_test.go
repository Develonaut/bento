package components

import (
	"strings"
	"testing"
)

type mockScreen int

func (m mockScreen) String() string {
	return "Browser"
}

func TestHeader(t *testing.T) {
	tabView := NewTabView()
	result := Header(tabView, 80)

	if !strings.Contains(result, "🍱 Bento") {
		t.Error("Header should contain Bento logo")
	}

	if !strings.Contains(result, "Bentos") {
		t.Error("Header should contain Bentos tab")
	}
}

func TestHeaderLegacy(t *testing.T) {
	screen := mockScreen(0)
	result := HeaderLegacy(screen, 80)

	if !strings.Contains(result, "🍱 Bento") {
		t.Error("Header should contain Bento logo")
	}

	if !strings.Contains(result, "Browser") {
		t.Error("Header should contain screen name")
	}
}
