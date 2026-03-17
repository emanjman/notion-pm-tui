package projectnotes

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	LeftFocus  key.Binding
	RightFocus key.Binding
}

var DefaultKeyMap = KeyMap{
	LeftFocus: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("<-/h", "left focus"),
	),
	RightFocus: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("->", "enter note"),
	),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.LeftFocus, k.RightFocus}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.LeftFocus, k.RightFocus},
		{},
	}
}
