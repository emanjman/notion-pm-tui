package version

import "github.com/charmbracelet/bubbles/key"

type NeutralKeyMap struct {
	Prev key.Binding
	Next key.Binding
}

var NeutralKeyMapper = NeutralKeyMap{
	Prev: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "prev"),
	),
	Next: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "next"),
	),
}

func (k NeutralKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Prev, k.Next}
}

func (k NeutralKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Prev, k.Next},
		{},
	}
}
