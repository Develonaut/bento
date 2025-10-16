package components

import (
	"github.com/charmbracelet/bubbles/key"
)

type PantryKeyMap struct {
	Navigation NavigationKeyMap
	ViewDetail key.Binding
}

func NewPantryKeyMap() PantryKeyMap {
	return PantryKeyMap{
		Navigation: NewNavigationKeyMap(),
		ViewDetail: key.NewBinding(
			key.WithKeys("enter", " "),
			key.WithHelp("enter", "view details"),
		),
	}
}

func (k PantryKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Navigation.Up,
		k.ViewDetail,
	}
}

func (k PantryKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Navigation.Up, k.Navigation.Down},
		{k.ViewDetail},
	}
}
