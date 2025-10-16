package components

import (
	"github.com/charmbracelet/bubbles/key"
)

type EditorKeyMap struct {
	Navigation NavigationKeyMap
	Edit       key.Binding
	Move       key.Binding
	Delete     key.Binding
	Add        key.Binding
	Save       key.Binding
	Run        key.Binding
	ToggleView key.Binding
	Cancel     key.Binding
}

func NewEditorKeyMap() EditorKeyMap {
	return EditorKeyMap{
		Navigation: NewNavigationKeyMap(),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit node"),
		),
		Move: key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("m", "move node"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete node"),
		),
		Add: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add node"),
		),
		Save: key.NewBinding(
			key.WithKeys("s", "ctrl+s"),
			key.WithHelp("s", "save"),
		),
		Run: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "run"),
		),
		ToggleView: key.NewBinding(
			key.WithKeys("v"),
			key.WithHelp("v", "toggle view"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
	}
}

func (k EditorKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Navigation.Up,
		k.Edit,
		k.Add,
		k.Save,
	}
}

func (k EditorKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Navigation.Up, k.Navigation.Down},
		{k.Edit, k.Move, k.Delete, k.Add},
		{k.Run, k.ToggleView, k.Save, k.Cancel},
	}
}
