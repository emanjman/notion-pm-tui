package milestone

import (
	"notion-project-tui/notion"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// dispatch messages to handlers
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	// handle milestones
	case notion.MilestonePagesMsg:
		return m.onMilestonePages(msg)
	case notion.FetchMoreMilestonesMsg:
		return m.onFetchMoreMilestones(msg)

	// handle milestone-tasks
	case notion.FetchMoreTasksMsg:
		return m.onFetchMoreTasksByStatus(msg)
	case notion.ToggleTaskGroupMsg:
		return m.onToggleTaskGroup(msg)
	case notion.TaskQueryMsg:
		return m.onTaskQuery(msg)

	// notion write operations
	case UpdateNotionTitleMsg:
		return m.onUpdateNotionTitle(msg)

	// chores
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

// dispatcher, where `onX()` logic still sits in `handlers.go`
func (m Model) handleKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch m.Focus.Mode {
	case WritingMode:
		switch {
		case key.Matches(msg, m.writingKeyMap.Save):
			return m.onWritingSave()
		default:
			var cmd tea.Cmd
			m.Focus.tempTitle, cmd = m.Focus.tempTitle.Update(msg)
			return m, cmd
		}
	case NeutralMode:
		switch {
		case key.Matches(msg, m.neutralKeyMap.Down):
			return m.onNeutralDown()
		case key.Matches(msg, m.neutralKeyMap.Up):
			return m.onNeutralUp()
		case key.Matches(msg, m.neutralKeyMap.JumpDown):
			return m.onNeutralJumpDown()
		case key.Matches(msg, m.neutralKeyMap.JumpUp):
			return m.onNeutralJumpUp()

		case key.Matches(msg, m.neutralKeyMap.Select):
			return m.onNeutralSelect()
		case key.Matches(msg, m.neutralKeyMap.Rename):
			return m.onNeutralRename()
		case key.Matches(msg, m.neutralKeyMap.Add):
			return m.onNeutralAdd()
		}
	}
	return m, nil
}
