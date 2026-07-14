package task

import "github.com/charmbracelet/bubbles/key"

type SelectKeyMap struct {
	Left   key.Binding // prev field
	Right  key.Binding // next field
	Select key.Binding // cycle select-options or enter rewrite mode
	Exit   key.Binding // send off changes to notion (server)
}

var SelectKeyMapper = SelectKeyMap{
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("<-/h", "prev field"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("->/l", "right field"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "enter edit mode"),
	),
	Exit: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "save + exit"),
	),
}

func (k SelectKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Left, k.Right, k.Select}
}

func (k SelectKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Left, k.Right, k.Select},
		{k.Exit},
	}
}
