package task

import "github.com/charmbracelet/bubbles/key"

type NormalKeyMap struct {
	Up         key.Binding
	Down       key.Binding
	JumpUp     key.Binding // jump up 5
	JumpDown   key.Binding // jump down 5
	Select     key.Binding // enter focus (select) mode
	StatusPrev key.Binding // cycle status backward
	StatusNext key.Binding // cycle status forward
	AddTask    key.Binding // add new task to idle group
	Delete     key.Binding // delete task (requires confirmation)
}

var NormalKeyMapper = NormalKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	JumpUp: key.NewBinding(
		key.WithKeys("ctrl+k"),
		key.WithHelp("ctrl+k", "jump up 5"),
	),
	JumpDown: key.NewBinding(
		key.WithKeys("ctrl+j"),
		key.WithHelp("ctrl+j", "jump down 5"),
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

func (k NormalKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select}
}

func (k NormalKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select},
		{k.StatusPrev, k.StatusNext, k.AddTask, k.Delete},
	}
}
