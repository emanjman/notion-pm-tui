package milestone

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type NeutralKeyMap struct {
	Up       key.Binding
	Down     key.Binding
	JumpUp   key.Binding // jump up 5
	JumpDown key.Binding // jump down 5
	Select   key.Binding // toggle group header
	Rename   key.Binding // enter writing mode
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
		key.WithHelp("enter", "select milestone"),
	),
	Rename: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rename milestone"),
	),
}

func (k NeutralKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select, k.Rename}
}

func (k NeutralKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select, k.Rename},
		{},
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

// ---

// dispatcher, where `onX()` logic still sits in `handlers.go`
func (m Model) handleKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch m.Focus.Mode {
	case WritingMode:
		switch {
		case key.Matches(msg, m.writingKeyMap.Save):
			return m.onWritingModeSave()
		default:
			var cmd tea.Cmd
			m.Focus.tempTitle, cmd = m.Focus.tempTitle.Update(msg)
			return m, cmd
		}
	case NeutralMode:
		switch {
		case key.Matches(msg, m.neutralKeyMap.Select):
			return m.onNeutralSelect()
		case key.Matches(msg, m.neutralKeyMap.Rename):
			return m.onNeutralRename()
		case key.Matches(msg, m.neutralKeyMap.Down):
			return m.onNeutralDown()
		case key.Matches(msg, m.neutralKeyMap.Up):
			return m.onNeutralUp()
		case key.Matches(msg, m.neutralKeyMap.JumpDown):
			return m.onNeutralJumpDown()
		case key.Matches(msg, m.neutralKeyMap.JumpUp):
			return m.onNeutralJumpUp()
		}
	}
	return m, nil
}
