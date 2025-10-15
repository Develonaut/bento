// Package screens provides individual TUI screens.
package screens

import (
	"bento/pkg/omise/components"
	"bento/pkg/omise/styles"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// BentoSelectedMsg signals that a bento was selected for execution
type BentoSelectedMsg struct {
	Name string
	Path string
}

// Browser shows available bentos
type Browser struct {
	list components.StyledList
}

// NewBrowser creates a browser screen
func NewBrowser() Browser {
	items := browserItems()
	l := components.NewStyledList(items, "Available Bentos")
	return Browser{list: l}
}

// browserItems returns bento list items
func browserItems() []list.Item {
	return []list.Item{
		bentoItem{
			name: "example-bento",
			path: "./bentos/example.bento.yaml",
			desc: "Example bento demonstrating HTTP and transforms",
		},
		bentoItem{
			name: "data-pipeline",
			path: "./bentos/data-pipeline.bento.yaml",
			desc: "Data processing pipeline with JQ transformations",
		},
		bentoItem{
			name: "api-integration",
			path: "./bentos/api-integration.bento.yaml",
			desc: "API integration with conditional logic",
		},
	}
}

// Init initializes the browser
func (b Browser) Init() tea.Cmd {
	return nil
}

// Update handles browser messages
func (b Browser) Update(msg tea.Msg) (Browser, tea.Cmd) {
	// Handle theme changes
	if _, ok := msg.(styles.ThemeChangedMsg); ok {
		b.list = b.list.RebuildStyles()
	}

	// Handle window resize to update list dimensions
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		h, v := lipgloss.NewStyle().Margin(2, 2).GetFrameSize()
		b.list.SetSize(msg.Width-h, msg.Height-v-4)
	}

	// Handle Enter or Space key to select bento
	if msg, ok := msg.(tea.KeyMsg); ok && (msg.String() == "enter" || msg.String() == " ") {
		if item, ok := b.list.SelectedItem().(bentoItem); ok {
			return b, func() tea.Msg {
				return BentoSelectedMsg{
					Name: item.name,
					Path: item.path,
				}
			}
		}
	}

	var cmd tea.Cmd
	b.list.Model, cmd = b.list.Model.Update(msg)
	return b, cmd
}

// View renders the browser
func (b Browser) View() string {
	return b.list.View()
}

// bentoItem represents a .bento.yaml file
type bentoItem struct {
	name string
	path string
	desc string
}

// Title returns the item title
func (i bentoItem) Title() string {
	return i.name
}

// Description returns the item description
func (i bentoItem) Description() string {
	return i.desc
}

// FilterValue returns the value to filter by
func (i bentoItem) FilterValue() string {
	return i.name
}
