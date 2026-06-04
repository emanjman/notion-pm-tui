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
	// handle queries
	case notion.QueryMilestonePagesMsg:
		return m.onQueryMilestonePages(msg)
	case notion.QueryMoreMilestonePagesMsg:
		return m.onQueryMoreMilestonePages(msg)
	case notion.QueryMoreTaskPagesMsg:
		return m.onQueryMoreTaskPages(msg)
	case notion.QueryTaskPagesMsg:
		return m.onQueryTaskPages(msg)

	// handle writes
	case notion.AddMilestonePageMsg:
		return m.onAddMilestonePage(msg)
	case UpdateNotionTitleMsg:
		return m.onUpdateNotionTitle(msg)
	case notion.ToggleTaskGroupMsg:
		return m.onToggleTaskGroup(msg)

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

// builds list of milestones, grouped by status; depends on `m.groups`
func (m Model) buildMilestoneList() []list.Item {
	var items []list.Item

	for _, status := range notion.MilestoneStatusOrder() {
		group, ok := m.groups[status]

		// skip groups w/o milestones
		if !ok || len(group.Milestones) == 0 {
			continue
		}

		// add header
		items = append(items, NewGroupHeaderItem(status, group))

		if !group.Hide {
			// add milestones
			for _, mstone := range group.Milestones {
				items = append(items, NewDefaultItem(mstone))
			}

			if group.NextCursor != nil {
				// add load-more button
				items = append(items, NewLoadMoreItem(status, group))
			}
		}
	}

	return items
}

// updates single milestone; triggers list rebuild
func (m Model) updateMilestone(item DefaultItem) Model {
	for status, group := range m.groups {
		for i, pg := range group.Milestones {
			if pg.ID == item.ID {
				// sync the name back onto the page (only field editable locally)
				group.Milestones[i].Properties.Title.Title[0].PlainText = item.Name
				return m.updateGroup(status, group)
			}
		}
	}
	return m
}

// updates single group in `m.groups`; triggers list rebuild
func (m Model) updateGroup(status notion.MilestoneStatus, group notion.MilestoneGroup) Model {
	m.groups[status] = group
	m.list.SetItems(m.buildMilestoneList())
	return m
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

// remove a milestone (by id) from its group + rebuild the list
func (m Model) removeMilestoneByID(id string) Model {
	for status, group := range m.groups {
		for i, pg := range group.Milestones {
			if pg.ID == id {
				group.Milestones = append(group.Milestones[:i], group.Milestones[i+1:]...)
				return m.updateGroup(status, group)
			}
		}
	}
	return m
}

// dispatcher, where `onX()` logic still sits in `handlers.go`
func (m Model) handleKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch m.Focus.Mode {
	case EditMode:
		switch {
		case key.Matches(msg, m.editKeyMap.Save):
			return m.onEditSave()
		default:
			var cmd tea.Cmd
			m.Focus.tempTitle, cmd = m.Focus.tempTitle.Update(msg)
			return m, cmd
		}
	case NeutralMode:
		switch {
		// navigation
		case key.Matches(msg, m.neutralKeyMap.Down):
			return m.onNeutralDown()
		case key.Matches(msg, m.neutralKeyMap.Up):
			return m.onNeutralUp()
		case key.Matches(msg, m.neutralKeyMap.JumpDown):
			return m.onNeutralJumpDown()
		case key.Matches(msg, m.neutralKeyMap.JumpUp):
			return m.onNeutralJumpUp()

		// change modes
		case key.Matches(msg, m.neutralKeyMap.Rename):
			return m.onNeutralRename()
		case key.Matches(msg, m.neutralKeyMap.Add):
			return m.onNeutralAdd()

		// dynamic: change mode, launch fetches
		case key.Matches(msg, m.neutralKeyMap.Select):
			return m.onNeutralSelect()

		}
	}
	return m, nil
}
