package tasklist

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
}

var DefaultKeyMap = KeyMap{
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
		key.WithHelp("enter", "select task"),
	),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select},
		{},
	}
}

// ----

// ? might have to rename the above keymap for disambiguity

type EditKeyMap struct {
	PrevField  key.Binding
	NextField  key.Binding
	EnableEdit key.Binding
	Exit       key.Binding
}

var DefaultEditKeyMap = EditKeyMap{
	PrevField: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("<-/h", "prev field"),
	),
	NextField: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("->/l", "right field"),
	),
	EnableEdit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "enter edit mode"),
	),
	Exit: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "save + exit"),
	),
}
