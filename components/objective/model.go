package objective

import (
	"notion-project-tui/components/objective/milestone"
	"notion-project-tui/components/objective/task"
	"notion-project-tui/notion"
	"notion-project-tui/util/keymap"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

type Panel int

const (
	MilestonePanel Panel = iota
	TaskPanel
)

type Model struct {
	projID           string
	versionID        string
	milestonesPropID string
	loading          bool
	err              error
	focus            Panel
	notion           *notion.Client
	milestone        milestone.Model
	task             task.Model
	keys             KeyMap
}

func New(n *notion.Client, projID, milestonesPropID string) Model {
	ms := milestone.New(n, projID)
	t := task.New(n)

	return Model{
		projID:           projID,
		milestonesPropID: milestonesPropID,
		loading:          true,
		err:              nil,
		focus:            MilestonePanel,
		notion:           n,
		milestone:        ms,
		task:             t,
		keys:             DefaultKeyMap,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.milestone.Init(), m.task.Init())
}

func (m Model) KeyMap() help.KeyMap {
	switch m.focus {
	case MilestonePanel:
		return keymap.JoinedKeyMap{Primary: m.keys, Secondary: m.milestone.ActiveKeyMap}
	case TaskPanel:
		return keymap.JoinedKeyMap{Primary: m.keys, Secondary: m.task.ActiveKeyMap}
	}
	return nil
}

func (m Model) InFocusMode() bool {
	return *m.milestone.Mode > milestone.NeutralMode || m.task.Focus.Mode > task.NeutralMode
}
