package app

import (
	"github.com/charmbracelet/bubbles/key"
	// tea "github.com/charmbracelet/bubbletea"
)

type KeyMap struct {
	Quit key.Binding
	Help key.Binding
	Next key.Binding
	Prev key.Binding
}

var RootKeyMap = KeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
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
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Help}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.Help},
		{k.Prev, k.Next},
	}
}
