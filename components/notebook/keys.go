package notebook

import (
	"github.com/charmbracelet/bubbles/key"
)

type BrowserKeyMap struct {
	Right key.Binding
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
}

var BrowserKeys = BrowserKeyMap{
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("->/l", "goto reader"),
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
		key.WithHelp("enter", "fetch/edit content"),
	),
}

func (k BrowserKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Right, k.Up, k.Down, k.Enter}
}

func (k BrowserKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Right, k.Up, k.Down, k.Enter},
		{},
	}
}

// ---

type ReaderKeyMap struct {
	Left       key.Binding
	Up         key.Binding
	JumpUp key.Binding
	Down       key.Binding
	JumpDown key.Binding
	Enter      key.Binding
	OpenEditor key.Binding
}

var ReaderKeys = ReaderKeyMap{
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("<-/h", "goto browser"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	JumpUp: key.NewBinding(
		key.WithKeys("ctrl+up", "ctrl+k"),
		key.WithHelp("ctrl+↑/ctrl+k", "up 5"),
	),
	JumpDown: key.NewBinding(
		key.WithKeys("ctrl+down", "ctrl+j"),
		key.WithHelp("ctrl+↓/ctrl+j", "down 5"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "edit mode"),
	),
	OpenEditor: key.NewBinding(
		key.WithKeys("ctrl+e"),
		key.WithHelp("ctrl+e", "open in editor"),
	),
}

func (k ReaderKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Left, k.Up, k.Down, k.Enter, k.OpenEditor}
}

func (k ReaderKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Left, k.Up, k.Down, k.Enter, k.OpenEditor},
		{},
	}
}

// ---

type EditorKeyMap struct {
	Esc key.Binding
}

var EditorKeys = EditorKeyMap{
	Esc: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "read mode"),
	),
}

func (k EditorKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Esc}
}

func (k EditorKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Esc},
		{},
	}
}
