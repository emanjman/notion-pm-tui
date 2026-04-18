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
	return tea.Batch(
		m.notion.QueryMilestones(m.projID, notion.MilestoneUnderDevelopment, ""),
		m.notion.QueryMilestones(m.projID, notion.MilestoneIdle, ""),
		m.notion.QueryMilestones(m.projID, notion.MilestoneComplete, ""),
	)
}

func (m Model) buildMilestoneList() []list.Item {
	var items []list.Item
	for _, status := range notion.MilestoneStatusOrder() {
		group, ok := m.groups[status]
		if !ok || len(group.Milestones) == 0 {
			continue
		}
		items = append(items, NewGroupHeaderItem(status, group))
		if !group.Hide {
			for _, pg := range group.Milestones {
				items = append(items, NewDefaultItem(pg))
			}
			if group.NextCursor != nil {
				items = append(items, NewLoadMoreItem(status, group))
			}
		}
	}
	return items
}

func (m Model) getCurrTaskGroups() notion.TaskGroups {
	item := m.list.SelectedItem()

	switch item := item.(type) {
	case GroupHeaderItem:
		group := m.groups[item.Status]
		if len(group.Milestones) > 0 {
			return NewDefaultItem(group.Milestones[0]).TaskGroups
		}
	case DefaultItem:
		return item.TaskGroups
	}
	return notion.TaskGroups{}
}

func (m *Model) SetItemDelegate(d list.ItemDelegate) {
	m.list.SetDelegate(d)
}

func (m Model) updateMilestoneInGroups(item DefaultItem) Model {
	group := m.groups[item.MilestoneStatus]

	for i, pg := range group.Milestones {
		if pg.ID == item.ID {
			// sync the name back onto the page (only field editable locally)
			group.Milestones[i].Properties.Title.Title[0].PlainText = item.Name
			break
		}
	}

	m.groups[item.MilestoneStatus] = group
	m.list.SetItems(m.buildMilestoneList())
	return m
}
