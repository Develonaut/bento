package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"bento/pkg/omise/styles"
)

// TabID identifies a tab
type TabID int

const (
	TabBentos TabID = iota
	TabRecipes
	TabMise
	TabSensei
	tabCount
)

// Tab represents a single tab
type Tab struct {
	ID    TabID
	Icon  string
	Name  string
	Key   string
	Index int
}

// TabView manages tab navigation and rendering
type TabView struct {
	activeTab TabID
	tabs      []Tab
	width     int
}

// NewTabView creates a new tab view
func NewTabView() TabView {
	return TabView{
		activeTab: TabBentos,
		tabs: []Tab{
			{ID: TabBentos, Icon: "🍱", Name: "Bentos", Key: "1", Index: 0},
			{ID: TabRecipes, Icon: "🍣", Name: "Pantry", Key: "2", Index: 1},
			{ID: TabMise, Icon: "🥢", Name: "Settings", Key: "3", Index: 2},
			{ID: TabSensei, Icon: "🍥", Name: "Help", Key: "4", Index: 3},
		},
		width: 80,
	}
}

// SetWidth updates the tab view width
func (tv TabView) SetWidth(width int) TabView {
	tv.width = width
	return tv
}

// SetActiveTab sets the active tab
func (tv TabView) SetActiveTab(tab TabID) TabView {
	if tab >= 0 && tab < tabCount {
		tv.activeTab = tab
	}
	return tv
}

// NextTab cycles to the next tab
func (tv TabView) NextTab() TabView {
	tv.activeTab = (tv.activeTab + 1) % tabCount
	return tv
}

// PrevTab cycles to the previous tab
func (tv TabView) PrevTab() TabView {
	if tv.activeTab == 0 {
		tv.activeTab = tabCount - 1
	} else {
		tv.activeTab--
	}
	return tv
}

// GetActiveTab returns the current active tab ID
func (tv TabView) GetActiveTab() TabID {
	return tv.activeTab
}

// View renders the tab bar
func (tv TabView) View() string {
	var renderedTabs []string

	numTabs := len(tv.tabs)
	baseWidth := tv.width / numTabs
	remainder := tv.width % numTabs

	for i, tab := range tv.tabs {
		tabWidth := baseWidth
		// Distribute remainder evenly starting from first tab
		if i < remainder {
			tabWidth++
		}
		renderedTabs = append(renderedTabs, tv.renderTab(tab, tabWidth))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	// Check if we need to add a gap to fill remaining width
	actualWidth := lipgloss.Width(row)
	if actualWidth < tv.width {
		gap := tv.width - actualWidth
		// Create gap with bottom border only to complete the tab bar line
		gapBorder := lipgloss.Border{
			Bottom: "─",
		}
		gapStyle := lipgloss.NewStyle().
			Width(gap).
			Border(gapBorder, false, false, true, false).
			BorderForeground(styles.InactiveTabBorder)
		row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gapStyle.Render(""))
	}

	return row
}

// renderTab renders a single tab
func (tv TabView) renderTab(tab Tab, width int) string {
	label := tab.Icon + " " + tab.Name

	isActive := tab.ID == tv.activeTab

	style := tv.getTabStyle(isActive, width)

	return style.Render(label)
}

// getTabStyle returns the appropriate style for a tab
func (tv TabView) getTabStyle(isActive bool, width int) lipgloss.Style {
	baseStyle := lipgloss.NewStyle().
		Padding(0, 1).
		Align(lipgloss.Center)

	var styled lipgloss.Style

	if isActive {
		// Active tab: blank bottom border to connect with content
		activeBorder := lipgloss.Border{
			Top:         "─",
			Bottom:      " ",
			Left:        "│",
			Right:       "│",
			TopLeft:     "╭",
			TopRight:    "╮",
			BottomLeft:  "┘",
			BottomRight: "└",
		}
		styled = baseStyle.
			Border(activeBorder, true, true, true, true).
			BorderForeground(styles.ActiveTabBorder).
			Bold(true).
			Foreground(styles.ActiveTabForeground)
	} else {
		// Inactive tab: full border
		inactiveBorder := lipgloss.Border{
			Top:         "─",
			Bottom:      "─",
			Left:        "│",
			Right:       "│",
			TopLeft:     "╭",
			TopRight:    "╮",
			BottomLeft:  "┴",
			BottomRight: "┴",
		}
		styled = baseStyle.
			Border(inactiveBorder, true, true, true, true).
			BorderForeground(styles.InactiveTabBorder).
			Foreground(styles.InactiveTabForeground)
	}

	// Calculate content width by subtracting frame size
	if width > 0 {
		frameSize := styled.GetHorizontalFrameSize() - 2
		contentWidth := width - frameSize
		if contentWidth > 0 {
			styled = styled.Width(contentWidth)
		}
	}

	return styled
}

// TabBar renders a simple tab indicator line
func (tv TabView) TabBar() string {
	var parts []string
	tabWidth := tv.width / len(tv.tabs)

	for _, tab := range tv.tabs {
		isActive := tab.ID == tv.activeTab
		var char string
		if isActive {
			char = "═"
		} else {
			char = "─"
		}
		parts = append(parts, strings.Repeat(char, tabWidth))
	}

	bar := strings.Join(parts, "")
	// Trim to exact width
	if len(bar) > tv.width {
		bar = bar[:tv.width]
	}
	return bar
}

// GetTabs returns all tabs
func (tv TabView) GetTabs() []Tab {
	return tv.tabs
}

// GetTab returns a tab by ID
func (tv TabView) GetTab(id TabID) (Tab, bool) {
	for _, tab := range tv.tabs {
		if tab.ID == id {
			return tab, true
		}
	}
	return Tab{}, false
}

// TabFromKey returns the tab ID for a given key
func (tv TabView) TabFromKey(key string) (TabID, bool) {
	for _, tab := range tv.tabs {
		if tab.Key == key {
			return tab.ID, true
		}
	}
	return 0, false
}
