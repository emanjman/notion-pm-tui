package milestone

import "github.com/charmbracelet/bubbles/key"

type EditKeyMap struct {
	Save key.Binding // update list item (client)
}

var EditKeyMapper = EditKeyMap{
	Save: key.NewBinding(
		key.WithKeys("enter", "esc"),
		key.WithHelp("enter/esc", "save changes"),
	),
}

func (k EditKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Save}
}

func (k EditKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Save},
		{},
	}
}
