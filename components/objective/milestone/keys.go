package milestone

import "github.com/charmbracelet/bubbles/key"

type NeutralKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding // toggle group header
	Rename key.Binding // enter writing mode
}

var NeutralKeyMapper = NeutralKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select milestone"),
	),
	Rename: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rename milestone"),
	),
}

func (k NeutralKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select, k.Rename}
}

func (k NeutralKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select, k.Rename},
		{},
	}
}

// ---

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
