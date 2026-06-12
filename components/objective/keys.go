package objective

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	LeftFocus           key.Binding
	RightFocus          key.Binding
	ToggleVersionSelect key.Binding
}

var DefaultKeyMap = KeyMap{
	LeftFocus: key.NewBinding(
		key.WithKeys("ctrl+h"),
		key.WithHelp("^h", "milestones"),
	),
	RightFocus: key.NewBinding(
		key.WithKeys("ctrl+l"),
		key.WithHelp("^l", "tasks"),
	),
	ToggleVersionSelect: key.NewBinding(
		key.WithKeys("ctrl+k"),
		key.WithHelp("^k", "versions"),
	),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.LeftFocus, k.RightFocus, k.ToggleVersionSelect}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.LeftFocus, k.RightFocus, k.ToggleVersionSelect},
		{},
	}
}
