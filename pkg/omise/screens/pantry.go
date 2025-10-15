package screens

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"bento/pkg/omise/styles"
)

// Pantry shows available neta types
type Pantry struct {
	table table.Model
}

// NewPantry creates a pantry screen
func NewPantry() Pantry {
	t := table.New(
		table.WithColumns(pantryColumns()),
		table.WithRows(pantryRows()),
		table.WithFocused(true),
		table.WithHeight(10),
	)
	t.SetStyles(pantryTableStyle())
	return Pantry{table: t}
}

// pantryColumns returns table column definitions
func pantryColumns() []table.Column {
	return []table.Column{
		{Title: "Type", Width: 25},
		{Title: "Category", Width: 15},
		{Title: "Description", Width: 40},
	}
}

// pantryRows returns neta type data
func pantryRows() []table.Row {
	return []table.Row{
		{"http", "Network", "HTTP request execution"},
		{"transform.jq", "Data", "JQ transformation"},
		{"conditional.if", "Control", "If/else conditional logic"},
		{"loop.for", "Control", "For loop iteration"},
		{"group.sequence", "Group", "Sequential execution"},
		{"group.parallel", "Group", "Parallel execution"},
	}
}

// pantryTableStyle returns styled table configuration
func pantryTableStyle() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(styles.Primary).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(styles.Primary).
		Bold(true)
	return s
}

// Init initializes the pantry
func (p Pantry) Init() tea.Cmd {
	return nil
}

// Update handles pantry messages
func (p Pantry) Update(msg tea.Msg) (Pantry, tea.Cmd) {
	// Handle window resize to update table dimensions
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		p.table.SetHeight(msg.Height - 10)
	}

	var cmd tea.Cmd
	p.table, cmd = p.table.Update(msg)
	return p, cmd
}

// View renders the pantry
func (p Pantry) View() string {
	title := styles.Title.Render("Pantry - Available Neta Types")
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		p.table.View(),
		"",
		styles.Subtle.Render("Use ↑/↓ to navigate • Press Enter for details"),
	)
}
