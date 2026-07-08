package objective

import (
	"notion-project-tui/components/objective/milestone"
	"notion-project-tui/components/objective/task"
	"notion-project-tui/components/objective/version"
	"notion-project-tui/notion"
	"notion-project-tui/util/keymap"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	projID  string
	loading bool
	err     error
	panel   Panel
	notion  *notion.Client

	version   version.Model
	milestone milestone.Model
	task      task.Model

	keys KeyMap
}

func New(n *notion.Client) Model {
	v := version.New(n)
	ms := milestone.New(n)
	t := task.New(n)

	return Model{
		projID:  "",
		loading: true,
		err:     nil,
		panel:   VersionPanel, // defaults here b/c version selection kicks off mstone-fetches
		notion:  n,

		version:   v,
		milestone: ms,
		task:      t,

		keys: DefaultKeyMap,
	}
}

func (m Model) Init(projID string) tea.Cmd {
	m.projID = projID
	return tea.Batch(m.version.Init(projID), m.milestone.Init(), m.task.Init())
}

// combine objective's native keymap w/ child keymap
func (m Model) KeyMap() help.KeyMap {
	switch m.panel {
	case VersionPanel:
		return keymap.JoinedKeyMap{Primary: m.keys, Secondary: m.version.ActiveKeyMap}
	case MilestonePanel:
		return keymap.JoinedKeyMap{Primary: m.keys, Secondary: m.milestone.ActiveKeyMap}
	case TaskPanel:
		return keymap.JoinedKeyMap{Primary: m.keys, Secondary: m.task.ActiveKeyMap}
	}
	return nil
}
