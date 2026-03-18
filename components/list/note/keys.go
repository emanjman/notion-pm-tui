package note

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	LeftFocus  key.Binding
	RightFocus key.Binding

	// todo: this one is specific to browsing
	FetchContent key.Binding
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

	FetchContent: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "fetch page content"),
	),

	// todo: i suppose there's top/down navs we can document
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.LeftFocus, k.RightFocus, k.FetchContent}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.LeftFocus, k.RightFocus, k.FetchContent},
		{},
	}
}
