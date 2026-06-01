package milestone

import tea "github.com/charmbracelet/bubbletea"

// handle select behavior based on item type
func (m Model) onNeutralSelect() (Model, tea.Cmd) {
	selected := m.list.SelectedItem()

	if header, ok := selected.(GroupHeaderItem); ok {
		// toggle show/hide header
		group := m.groups[header.Status]
		group.Hide = !group.Hide
		m.groups[header.Status] = group
		m.list.SetItems(m.buildMilestoneList())

	} else if loadMore, ok := selected.(LoadMoreItem); ok && !loadMore.Loading {
		loadMoreMilestones(loadMore.Status)

	} else if mstone, ok := selected.(DefaultItem); ok {
		// fetch tasks for curr milestone
		switch mstone.FetchStatus {
		case FetchIdle:
			idx := m.list.Index()
			mstone.FetchStatus = FetchPending
			m.list.SetItem(idx, mstone)
			if mstone.TaskCount > 0 {
				return m, fetchInitTasks(mstone.ID, idx, m.notion)
			}
			return m, nil

		case FetchSuccess:
			return m, refreshMilestoneTasks(mstone.TaskGroups)
		}
	}

	return m, nil
}

// enter writing-mode
func (m Model) onNeutralRename() (Model, tea.Cmd) {
	if mstone, ok := m.list.SelectedItem().(DefaultItem); ok {
		// id milestone to update
		m.Focus.milestoneID = mstone.ID
		m.Focus.milestoneIdx = m.list.Index()

		// setup text input model
		m.Focus.tempTitle = initTempTitle(mstone)

		// flip to writing-mode
		m.Focus.Mode = WritingMode
		m.ActiveKeyMap = WritingKeyMapper
	}
	return m, nil
}

// nav down 1 + refresh
func (m Model) onNeutralDown() (Model, tea.Cmd) {
	m.list.CursorDown()
	return m, refreshMilestoneTasks(m.getCurrTaskGroups())
}

// nav up 1 + refresh
func (m Model) onNeutralUp() (Model, tea.Cmd) {
	m.list.CursorUp()
	return m, refreshMilestoneTasks(m.getCurrTaskGroups())
}

// nav down 5 + refresh
func (m Model) onNeutralJumpDown() (Model, tea.Cmd) {
	m.list.Select(min(len(m.list.Items())-1, m.list.Index()+5))
	return m, refreshMilestoneTasks(m.getCurrTaskGroups())
}

// nav up 5 + refresh
func (m Model) onNeutralJumpUp() (Model, tea.Cmd) {
	m.list.Select(max(0, m.list.Index()-5))
	return m, refreshMilestoneTasks(m.getCurrTaskGroups())
}
