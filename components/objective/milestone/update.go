package milestone

import (
	"notion-project-tui/notion"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// dispatch messages to handlers
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case notion.MilestonePagesMsg:
		return m.onMilestonePages(msg)
	case notion.FetchMoreMilestonesMsg:
		return m.onFetchMoreMilestones(msg)
	case notion.FetchMoreTasksMsg:
		return m.onFetchMoreTasksByStatus(msg)
	case notion.ToggleTaskGroupMsg:
		return m.onToggleTaskGroup(msg)
	case notion.TaskQueryMsg:
		return m.onTaskQuery(msg)
	case UpdateNotionTitleMsg:
		return m.onUpdateTitle(msg)
	case tea.KeyMsg:
		return m.handleKey(msg)
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
	}

	// otherwise, handle from children
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
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
