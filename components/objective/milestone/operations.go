package milestone

import (
	"github.com/charmbracelet/bubbles/list"
	"notion-project-tui/notion"
)

// this file contains general operations on the model outside
// the required ones

// public list operations
func (m Model) ClearMilestones() Model {
	m.list.SetItems([]list.Item{})
	m.groups = notion.MilestoneGroups{}
	m.pendingFetches = 3
	return m
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

// id of the milestone backing the current selection; tasks created in the
// task panel hang off this via the @milestone relation
func (m Model) getCurrMilestoneID() string {
	item := m.list.SelectedItem()

	switch item := item.(type) {
	case GroupHeaderItem:
		group := m.groups[item.Status]
		if len(group.Milestones) > 0 {
			return group.Milestones[0].ID
		}
	case DefaultItem:
		return item.ID
	}
	return ""
}

// todo: this should be a value receiver, return model
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
