package app

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type GlobalKeyMap struct {
	Quit key.Binding
	Help key.Binding
}

var RootKeyMap = GlobalKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c/q", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
}

func (k GlobalKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Help}
}

func (k GlobalKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.Help},
		{},
	}
}

// ---
// merge the root-level keys w/ the curr view's keys

type MergedKeyMap struct {
	curr   help.KeyMap
	global GlobalKeyMap
}

func (m MergedKeyMap) ShortHelp() []key.Binding {
	return append(m.curr.ShortHelp(), m.global.Help)
}

func (m MergedKeyMap) FullHelp() [][]key.Binding {
	return append(m.curr.FullHelp(), m.global.FullHelp()...)
}
