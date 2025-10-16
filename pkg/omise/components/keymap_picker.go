package components

import (
	"github.com/charmbracelet/bubbles/key"
)

type PickerKeyMap struct {
	Navigation NavigationKeyMap
	Select     key.Binding
	Reset      key.Binding
	Cancel     key.Binding
}

func NewPickerKeyMap() PickerKeyMap {
	return PickerKeyMap{
		Navigation: NewNavigationKeyMap(),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Reset: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "reset to default"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
	}
}

func (k PickerKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Navigation.Up,
		k.Select,
		k.Reset,
	}
}

func (k PickerKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Navigation.Up, k.Navigation.Down},
		{k.Select, k.Reset, k.Cancel},
	}
}
