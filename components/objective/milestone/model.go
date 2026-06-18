package milestone

import (
	"notion-project-tui/notion"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	notion         *notion.Client
	versionID      string
	list           list.Model
	err            error
	pendingFetches int
	groups         notion.MilestoneGroups

	Mode   *Mode        // ptr so the list delegate sees mode switches live (shared heap)
	Edit   *EditModeCtx // todo: do these NEED to be ptrs?
	Delete *DeleteModeCtx

	ActiveKeyMap  help.KeyMap // for help focus view
	neutralKeyMap NeutralKeyMap
	editKeyMap    EditKeyMap
	deleteKeyMap  DeleteKeyMap
}

func New(n *notion.Client, versionID string) Model {
	mode := NeutralMode
	edit := EditModeCtx{}
	del := DeleteModeCtx{}

	l := list.New([]list.Item{}, NewItemDelegate(true, &mode, &edit), 0, 0)
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
		versionID:      versionID,
		pendingFetches: 3,
		list:           l,
		err:            nil,
		groups:         notion.MilestoneGroups{},

		ActiveKeyMap:  NeutralKeyMapper, // default map view
		neutralKeyMap: NeutralKeyMapper,
		editKeyMap:    EditKeyMapper,
		deleteKeyMap:  DeleteKeyMapper,

		Mode:   &mode,
		Edit:   &edit,
		Delete: &del,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
	// return fetchMilestonesByStatus(m.versionID, m.notion)
}
