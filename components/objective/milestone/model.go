package milestone

import (
	"notion-project-tui/notion"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// carries updated task groups for the selected milestone to the task panel
type MilestoneTasksMsg struct {
	Groups notion.TaskGroups // task pages keyed by status, each w/ pagination + visibility state
}

// result of notion milestone title rename
type UpdateNotionTitleMsg struct{ Err error }

type Model struct {
	notion         *notion.Client
	projID         string
	list           list.Model
	err            error
	pendingFetches int
	groups         notion.MilestoneGroups
	ActiveKeyMap   help.KeyMap // for help focus view
	neutralKeyMap  NeutralKeyMap
	writingKeyMap  WritingKeyMap
	Focus          *FocusState

	tempIDCounter int // for generating temp IDs for new milestones
}

func New(n *notion.Client, projID string) Model {
	f := FocusState{}

	l := list.New([]list.Item{}, NewItemDelegate(true, &f), 0, 0)
	l.Title = "Milestones"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	l.SetShowFilter(false)
	l.SetShowPagination(false)
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()

	return Model{
		notion:         n,
		projID:         projID,
		pendingFetches: 3,
		list:           l,
		err:            nil,
		groups:         notion.MilestoneGroups{},
		ActiveKeyMap:   NeutralKeyMapper, // default map view
		neutralKeyMap:  NeutralKeyMapper,
		writingKeyMap:  WritingKeyMapper,
		Focus:          &f,
	}
}

func (m Model) Init() tea.Cmd {
	return fetchMilestonesByStatus(m.projID, m.notion)
}
