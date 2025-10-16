package components

import (
	"github.com/charmbracelet/bubbles/key"
)

type GlobalKeyMap struct {
	Quit     key.Binding
	Help     key.Binding
	Settings key.Binding
	Back     key.Binding
}

func NewGlobalKeyMap() GlobalKeyMap {
	return GlobalKeyMap{
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Settings: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "settings"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
	}
}

func (k GlobalKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Settings, k.Help, k.Quit}
}

func (k GlobalKeyMap) ShortHelpWithBack() []key.Binding {
	return []key.Binding{k.Back, k.Help, k.Quit}
}

func (k GlobalKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Settings, k.Help},
		{k.Quit},
	}
}
