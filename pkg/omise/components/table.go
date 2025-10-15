package components

import (
	"bento/pkg/omise/styles"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// StyledTable wraps bubbles/table with theme-aware styling
type StyledTable struct {
	table.Model
}

// NewStyledTable creates a themed table
func NewStyledTable(columns []table.Column, rows []table.Row, height int) StyledTable {
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(height),
	)
	t.SetStyles(tableStyles())
	return StyledTable{Model: t}
}

// tableStyles returns table styles with theme colors
func tableStyles() table.Styles {
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

// RebuildStyles updates the table styles with current theme colors
func (st StyledTable) RebuildStyles() StyledTable {
	st.Model.SetStyles(tableStyles())
	return st
}

// Update handles table messages
func (st StyledTable) Update(msg tea.Msg) (StyledTable, tea.Cmd) {
	var cmd tea.Cmd
	st.Model, cmd = st.Model.Update(msg)
	return st, cmd
}
