package milestone

import "github.com/charmbracelet/bubbles/key"

type NeutralKeyMap struct {
	Up       key.Binding
	Down     key.Binding
	JumpUp   key.Binding // jump up 5
	JumpDown key.Binding // jump down 5
	Select   key.Binding // toggle group header
	Rename   key.Binding
	Add      key.Binding
	Delete   key.Binding
}

var NeutralKeyMapper = NeutralKeyMap{
	Up: key.NewBinding(
		key.WithKeys("k"),
		key.WithHelp("k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j"),
		key.WithHelp("j", "down"),
	),
	JumpUp: key.NewBinding(
		key.WithKeys("K"),
		key.WithHelp("K", "jump up"),
	),
	JumpDown: key.NewBinding(
		key.WithKeys("J"),
		key.WithHelp("J", "jump down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select milestone"),
	),
	Rename: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rename milestone"),
	),
	Add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add milestone"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete milestone"),
	),
}

func (k NeutralKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select, k.Rename, k.Add, k.Delete}
}

func (k NeutralKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select, k.Rename, k.Add, k.Delete},
		{},
	}
}
