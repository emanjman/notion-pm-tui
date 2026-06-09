package objective

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	LeftFocus           key.Binding
	RightFocus          key.Binding
	ToggleVersionSelect key.Binding
}

var DefaultKeyMap = KeyMap{
	LeftFocus: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("<-/h", "left focus"),
	),
	RightFocus: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("->/l", "right focus"),
	),
	ToggleVersionSelect: key.NewBinding(
		key.WithKeys("v"),
		key.WithHelp("v", "toggle version"),
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
