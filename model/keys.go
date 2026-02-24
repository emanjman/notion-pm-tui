package model

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Back   key.Binding
	Quit   key.Binding
	Help   key.Binding
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

	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c/q", "quit"),
	),

	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Help, k.Up, k.Down}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.Help, k.Up, k.Down},
		{},
	}
}

/*
in the Update() function:

case tea.KeyMsg:
    switch {
    case key.Matches(msg, m.keys.Quit):
        return m, tea.Quit
    case key.Matches(msg, m.keys.Up):
        m.cursor--
    }
*/
