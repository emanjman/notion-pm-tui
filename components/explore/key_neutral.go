package explore

import "github.com/charmbracelet/bubbles/key"

type NeutralKeyMap struct {
	Select key.Binding
	Quit   key.Binding
}

var NeutralKeyMapper = NeutralKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select project"),
	),
}

func (k NeutralKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Select}
}

func (k NeutralKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.Select},
		{},
	}
}
