package components

import (
	"bento/pkg/omise/styles"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// StyledList wraps bubbles/list with theme-aware styling
type StyledList struct {
	list.Model
}

// NewStyledList creates a new themed list
func NewStyledList(items []list.Item, title string) StyledList {
	delegate := styledDelegate()
	l := list.New(items, delegate, 0, 0)
	l.Title = title
	l.Styles = listStyles()
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.SetStatusBarItemName("bento", "bentos")

	return StyledList{Model: l}
}

// listStyles returns list styles with theme colors
func listStyles() list.Styles {
	s := list.DefaultStyles()
	s.Title = lipgloss.NewStyle().Foreground(styles.Primary).Bold(true)
	s.Spinner = lipgloss.NewStyle().Foreground(styles.Primary)
	s.FilterPrompt = lipgloss.NewStyle().Foreground(styles.Primary)
	s.FilterCursor = lipgloss.NewStyle().Foreground(styles.Primary)
	s.DefaultFilterCharacterMatch = lipgloss.NewStyle().Underline(true)
	s.StatusBar = lipgloss.NewStyle().Foreground(styles.Muted)
	s.StatusBarActiveFilter = lipgloss.NewStyle().Foreground(styles.Text)
	s.NoItems = lipgloss.NewStyle().Foreground(styles.Muted)
	s.PaginationStyle = lipgloss.NewStyle().Foreground(styles.Muted)
	s.HelpStyle = lipgloss.NewStyle().Foreground(styles.Muted)
	s.ActivePaginationDot = lipgloss.NewStyle().Foreground(styles.Primary)
	s.InactivePaginationDot = lipgloss.NewStyle().Foreground(styles.Secondary)
	return s
}

// styledDelegate returns a delegate with theme styling
func styledDelegate() list.DefaultDelegate {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.
		BorderLeftForeground(styles.Secondary)
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(styles.Primary).
		BorderLeftForeground(styles.Primary)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(styles.Muted).
		BorderLeftForeground(styles.Secondary)
	return delegate
}

// RebuildStyles updates the list styles with current theme colors
func (sl StyledList) RebuildStyles() StyledList {
	sl.Model.Styles = listStyles()
	delegate := styledDelegate()
	sl.Model.SetDelegate(delegate)
	return sl
}
