package notebook

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	LeftFocus  key.Binding
	RightFocus key.Binding

	Up   key.Binding
	Down key.Binding

	Enter key.Binding
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

	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),

	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "fetch content"),
	),

	// todo: i suppose there's top/down navs we can document
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.LeftFocus, k.RightFocus}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.LeftFocus, k.RightFocus},
		{},
	}
}
