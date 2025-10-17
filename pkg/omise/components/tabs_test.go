package components

import (
	"strings"
	"testing"
)

func TestNewTabView(t *testing.T) {
	tv := NewTabView()

	if tv.GetActiveTab() != TabBentos {
		t.Error("Default active tab should be Bentos")
	}

	tabs := tv.GetTabs()
	if len(tabs) != 4 {
		t.Errorf("Expected 4 tabs, got %d", len(tabs))
	}
}

func TestTabSwitching(t *testing.T) {
	tv := NewTabView()

	// Test NextTab
	tv = tv.NextTab()
	if tv.GetActiveTab() != TabRecipes {
		t.Error("NextTab should switch to Recipes")
	}

	// Test PrevTab
	tv = tv.PrevTab()
	if tv.GetActiveTab() != TabBentos {
		t.Error("PrevTab should switch back to Bentos")
	}

	// Test wrapping
	tv = tv.PrevTab()
	if tv.GetActiveTab() != TabSensei {
		t.Error("PrevTab should wrap to Sensei")
	}
}

func TestSetActiveTab(t *testing.T) {
	tv := NewTabView()

	tv = tv.SetActiveTab(TabMise)
	if tv.GetActiveTab() != TabMise {
		t.Error("SetActiveTab should switch to Mise")
	}
}

func TestTabFromKey(t *testing.T) {
	tv := NewTabView()

	tests := []struct {
		key      string
		expected TabID
		ok       bool
	}{
		{"1", TabBentos, true},
		{"2", TabRecipes, true},
		{"3", TabMise, true},
		{"4", TabSensei, true},
		{"5", 0, false},
		{"x", 0, false},
	}

	for _, tt := range tests {
		tabID, ok := tv.TabFromKey(tt.key)
		if ok != tt.ok {
			t.Errorf("TabFromKey(%q) ok = %v, want %v", tt.key, ok, tt.ok)
		}
		if ok && tabID != tt.expected {
			t.Errorf("TabFromKey(%q) = %v, want %v", tt.key, tabID, tt.expected)
		}
	}
}

func TestTabView(t *testing.T) {
	tv := NewTabView().SetWidth(80)
	result := tv.View()

	// Check that all tabs are present
	if !strings.Contains(result, "🍱 Bentos") {
		t.Error("View should contain Bentos tab")
	}
	if !strings.Contains(result, "🍣 Pantry") {
		t.Error("View should contain Pantry tab")
	}
	if !strings.Contains(result, "🥢 Settings") {
		t.Error("View should contain Settings tab")
	}
	if !strings.Contains(result, "🍥 Help") {
		t.Error("View should contain Help tab")
	}
}

func TestGetTab(t *testing.T) {
	tv := NewTabView()

	tab, ok := tv.GetTab(TabBentos)
	if !ok {
		t.Error("GetTab should find Bentos tab")
	}
	if tab.Name != "Bentos" {
		t.Errorf("Expected tab name Bentos, got %s", tab.Name)
	}

	_, ok = tv.GetTab(TabID(999))
	if ok {
		t.Error("GetTab should not find invalid tab")
	}
}
