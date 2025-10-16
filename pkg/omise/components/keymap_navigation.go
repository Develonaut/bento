package components

import (
	"github.com/charmbracelet/bubbles/key"
)

type NavigationKeyMap struct {
	Up   key.Binding
	Down key.Binding
}

func NewNavigationKeyMap() NavigationKeyMap {
	return NavigationKeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
	}
}

func (k NavigationKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down}
}

func (k NavigationKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
	}
}
