package components

import (
	"github.com/charmbracelet/bubbles/key"
)

type SettingsKeyMap struct {
	Navigation NavigationKeyMap
	Select     key.Binding
	Reset      key.Binding
	Back       key.Binding
}

func NewSettingsKeyMap() SettingsKeyMap {
	return SettingsKeyMap{
		Navigation: NewNavigationKeyMap(),
		Select: key.NewBinding(
			key.WithKeys("enter", " "),
			key.WithHelp("enter", "select"),
		),
		Reset: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "reset to default"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
	}
}

func (k SettingsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Navigation.Up,
		k.Select,
		k.Reset,
	}
}

func (k SettingsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Navigation.Up, k.Navigation.Down},
		{k.Select, k.Reset, k.Back},
	}
}
