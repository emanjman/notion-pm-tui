package explore

import "github.com/charmbracelet/bubbles/key"

type NeutralKeyMap struct {
	Select key.Binding
}

var NeutralKeyMapper = NeutralKeyMap{
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select project"),
	),
}

func (k NeutralKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Select}
}

func (k NeutralKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Select},
		{},
	}
}
