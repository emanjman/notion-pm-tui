package task

import "github.com/charmbracelet/bubbles/key"

type NeutralKeyMap struct {
	Up         key.Binding
	Down       key.Binding
	Select     key.Binding // enter focus (select) mode
	StatusPrev key.Binding // cycle status backward
	StatusNext key.Binding // cycle status forward
	AddTask    key.Binding // add new task to idle group
	Delete     key.Binding // delete task (requires confirmation)
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
		key.WithHelp("enter", "select task"),
	),
	StatusPrev: key.NewBinding(
		key.WithKeys("<", "shift+,"),
		key.WithHelp("<", "prev status"),
	),
	StatusNext: key.NewBinding(
		key.WithKeys(">", "shift+."),
		key.WithHelp(">", "next status"),
	),
	AddTask: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add task"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete task"),
	),
}

func (k NeutralKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select}
}

func (k NeutralKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select},
		{k.StatusPrev, k.StatusNext, k.AddTask, k.Delete},
	}
}

// ---

type SelectingKeyMap struct {
	Left   key.Binding // prev field
	Right  key.Binding // next field
	Select key.Binding // cycle select-options or enter rewrite mode
	Exit   key.Binding // send off changes to notion (server)
}

var SelectingKeyMapper = SelectingKeyMap{
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

func (k SelectingKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Left, k.Right, k.Select}
}

func (k SelectingKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Left, k.Right, k.Select},
		{k.Exit},
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
