package components

import (
	"github.com/charmbracelet/bubbles/key"
)

type BrowserKeyMap struct {
	Navigation  NavigationKeyMap
	Execute     key.Binding
	Run         key.Binding
	New         key.Binding
	Search      key.Binding
	ClearSearch key.Binding
	Edit        key.Binding
	Copy        key.Binding
	Delete      key.Binding
}

func NewBrowserKeyMap() BrowserKeyMap {
	return BrowserKeyMap{
		Navigation: NewNavigationKeyMap(),
		Execute: key.NewBinding(
			key.WithKeys("enter", " "),
			key.WithHelp("enter", "show actions"),
		),
		Run: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "run"),
		),
		New: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new bento"),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
		ClearSearch: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "clear search"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit"),
		),
		Copy: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "copy"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete"),
		),
	}
}

func (k BrowserKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Search,
	}
}

// ActionHelp returns action keys for selected items
func (k BrowserKeyMap) ActionHelp() []key.Binding {
	return []key.Binding{
		k.Execute,
		k.Run,
		k.Edit,
		k.Copy,
		k.Delete,
	}
}

func (k BrowserKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Navigation.Up, k.Navigation.Down},
		{k.Run, k.Edit, k.Copy, k.Delete},
		{k.New, k.Search, k.ClearSearch, k.Execute},
	}
}
