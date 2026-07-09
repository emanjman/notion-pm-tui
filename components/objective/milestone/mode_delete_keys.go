package milestone

import "github.com/charmbracelet/bubbles/key"

type DeleteKeyMap struct {
	Confirm key.Binding
	Cancel  key.Binding
}

var DeleteKeyMapper = DeleteKeyMap{
	Confirm: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("delete", "d"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("escape", "esc"),
	),
}

func (k DeleteKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Confirm, k.Cancel}
}

func (k DeleteKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Confirm, k.Cancel},
		{},
	}
}
