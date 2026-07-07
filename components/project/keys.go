package project

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Help key.Binding
	Next key.Binding
	Prev key.Binding
	Quit key.Binding
}

var RootKeyMap = KeyMap{
	// todo: should this be in `explore` model?
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Next: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next tab"),
	),
	Prev: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift tab", "prev tab"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Help},
		{k.Prev, k.Next},
	}
}
