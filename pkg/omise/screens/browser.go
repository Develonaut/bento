// Package screens provides individual TUI screens.
package screens

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Browser shows available workflows
type Browser struct {
	list list.Model
}

// NewBrowser creates a browser screen
func NewBrowser() Browser {
	delegate := browserDelegate()
	l := list.New(browserItems(), delegate, 0, 0)
	l.Title = "Available Workflows"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	return Browser{list: l}
}

// browserItems returns workflow list items
func browserItems() []list.Item {
	return []list.Item{
		workflowItem{
			name: "example-workflow",
			path: "./workflows/example.bento.yaml",
			desc: "Example workflow demonstrating HTTP and transforms",
		},
		workflowItem{
			name: "data-pipeline",
			path: "./workflows/data-pipeline.bento.yaml",
			desc: "Data processing pipeline with JQ transformations",
		},
		workflowItem{
			name: "api-integration",
			path: "./workflows/api-integration.bento.yaml",
			desc: "API integration with conditional logic",
		},
	}
}

// browserDelegate returns styled list delegate
func browserDelegate() list.DefaultDelegate {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("205")).
		BorderLeftForeground(lipgloss.Color("205"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("241"))
	return delegate
}

// Init initializes the browser
func (b Browser) Init() tea.Cmd {
	return nil
}

// Update handles browser messages
func (b Browser) Update(msg tea.Msg) (Browser, tea.Cmd) {
	var cmd tea.Cmd
	b.list, cmd = b.list.Update(msg)
	return b, cmd
}

// View renders the browser
func (b Browser) View() string {
	return b.list.View()
}

// workflowItem represents a .bento.yaml file
type workflowItem struct {
	name string
	path string
	desc string
}

// Title returns the item title
func (i workflowItem) Title() string {
	return i.name
}

// Description returns the item description
func (i workflowItem) Description() string {
	return i.desc
}

// FilterValue returns the value to filter by
func (i workflowItem) FilterValue() string {
	return i.name
}
