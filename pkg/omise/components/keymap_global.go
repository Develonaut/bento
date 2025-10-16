package components

import (
	"github.com/charmbracelet/bubbles/key"
)

type GlobalKeyMap struct {
	Quit       key.Binding
	ForceQuit  key.Binding
	Help       key.Binding
	NextScreen key.Binding
	PrevScreen key.Binding
}

func NewGlobalKeyMap() GlobalKeyMap {
	return GlobalKeyMap{
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
		ForceQuit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "force quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?", "h"),
			key.WithHelp("?", "toggle help"),
		),
		NextScreen: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next screen"),
		),
		PrevScreen: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev screen"),
		),
	}
}

func (k GlobalKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k GlobalKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.NextScreen, k.PrevScreen, k.Help},
		{k.Quit, k.ForceQuit},
	}
}
