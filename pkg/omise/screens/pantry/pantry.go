package pantry

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"bento/pkg/omise/components"
	"bento/pkg/omise/styles"
)

// Pantry shows available neta types
type Pantry struct {
	table    components.StyledTable
	helpView components.HelpView
	keys     components.PantryKeyMap
}

// NewPantry creates a pantry screen
func NewPantry() Pantry {
	return Pantry{
		table: components.NewStyledTable(
			pantryColumns(),
			pantryRows(),
			10,
		),
		helpView: components.NewHelpView(),
		keys:     components.NewPantryKeyMap(),
	}
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

// Init initializes the pantry
func (p Pantry) Init() tea.Cmd {
	return nil
}

// Update handles pantry messages
func (p Pantry) Update(msg tea.Msg) (Pantry, tea.Cmd) {
	// Handle theme changes
	if _, ok := msg.(styles.ThemeChangedMsg); ok {
		p.table = p.table.RebuildStyles()
	}

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
	)
}

// KeyBindings returns the contextual key bindings for the footer
func (p Pantry) KeyBindings() []key.Binding {
	// Pantry is read-only, no contextual actions
	return []key.Binding{}
}
