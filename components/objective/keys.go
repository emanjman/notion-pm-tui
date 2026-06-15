package objective

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	FocusMilestones key.Binding
	FocusTasks      key.Binding
	FocusVersions   key.Binding
	UnfocusVersions key.Binding
}

var DefaultKeyMap = KeyMap{
	FocusMilestones: key.NewBinding(
		key.WithKeys("ctrl+h"),
		key.WithHelp("^h", "focus milestones"),
	),
	FocusTasks: key.NewBinding(
		key.WithKeys("ctrl+l"),
		key.WithHelp("^l", "focus tasks"),
	),
	FocusVersions: key.NewBinding(
		key.WithKeys("ctrl+k"),
		key.WithHelp("^k", "focus versions"),
	),
	UnfocusVersions: key.NewBinding(
		key.WithKeys("ctrl+j"),
		key.WithHelp("^j", "unfocus versions"),
	),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.FocusMilestones, k.FocusTasks, k.FocusVersions}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.FocusMilestones, k.FocusTasks, k.FocusVersions},
		{},
	}
}
