package milestone

import "github.com/charmbracelet/bubbles/key"

type WritingKeyMap struct {
	Save key.Binding // update list item (client)
}

var WritingKeyMapper = WritingKeyMap{
	Save: key.NewBinding(
		key.WithKeys("enter", "esc"),
		key.WithHelp("enter/esc", "save changes"),
	),
}

func (k WritingKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Save}
}

func (k WritingKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Save},
		{},
	}
}
