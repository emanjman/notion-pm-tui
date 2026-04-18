package milestone

import (
	"notion-project-tui/notion"

	tea "github.com/charmbracelet/bubbletea"
)

// receiver methods `onXxx()` methods

func (m Model) onWritingModeSave() (Model, tea.Cmd) {
	if milestone, ok := m.list.SelectedItem().(DefaultItem); ok {
		milestone.Name = m.Focus.tempTitle.Value()
		m.list.SetItem(m.list.Index(), milestone)
		m.updateMilestoneInGroups(milestone)
	}

	m.ActiveKeyMap = NeutralKeyMapper
	m.Focus.Mode = NeutralMode

	// todo: send command to update milestone title in notion
	return m, nil
}

func (m Model) onNeutralSelect() (Model, tea.Cmd) {
	selected := m.list.SelectedItem()

	if header, ok := selected.(GroupHeaderItem); ok {
		// toggle header
		group := m.groups[header.Status]
		group.Hide = !group.Hide
		m.groups[header.Status] = group
		m.list.SetItems(m.buildMilestoneList())

	} else if loadMore, ok := selected.(LoadMoreItem); ok && !loadMore.Loading {
		// load more milestones
		return m, func() tea.Msg {
			return notion.FetchMoreMilestonesMsg{Status: loadMore.Status}
		}
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

		// todo: is this necessary?
		case FetchSuccess:
			return m, refreshMilestoneTasks(mstone.TaskGroups)
		}
	}
	return m, nil
}

func (m Model) onNeutralRename() (Model, tea.Cmd) {
	if milestone, ok := m.list.SelectedItem().(DefaultItem); ok {
		m.Focus.milestoneID = milestone.ID
		m.Focus.milestoneIdx = m.list.Index()
		m.Focus.tempTitle = initTempTitle(milestone)
		m.Focus.Mode = WritingMode
		m.ActiveKeyMap = WritingKeyMapper
	}
	return m, nil
}

func (m Model) onNeutralDown() (Model, tea.Cmd) {
	m.list.CursorDown()
	return m, refreshMilestoneTasks(m.getCurrTaskGroups())
}

func (m Model) onNeutralUp() (Model, tea.Cmd) {
	m.list.CursorUp()
	return m, refreshMilestoneTasks(m.getCurrTaskGroups())
}

func (m Model) onNeutralJumpDown() (Model, tea.Cmd) {
	m.list.Select(min(len(m.list.Items())-1, m.list.Index()+5))
	return m, refreshMilestoneTasks(m.getCurrTaskGroups())
}

func (m Model) onNeutralJumpUp() (Model, tea.Cmd) {
	m.list.Select(max(0, m.list.Index()-5))
	return m, refreshMilestoneTasks(m.getCurrTaskGroups())
}

// -----------

func (m Model) onMilestonePages(msg notion.MilestonePagesMsg) (Model, tea.Cmd) {
	if msg.Err != nil {
		m.err = msg.Err
		m.pendingFetches--
		return m, nil
	}

	// append incoming pages into the correct status group
	group := m.groups[msg.Status]
	group.Milestones = append(group.Milestones, msg.Pages...)
	group.NextCursor = msg.NextCursor
	m.groups[msg.Status] = group

	m.pendingFetches--

	// only render + kick off task fetches once all 3 status batches have arrived
	if m.pendingFetches > 0 {
		return m, nil
	}

	m.list.SetItems(m.buildMilestoneList())
	return m, fetchInitMilestoneTasks(&m.list, m.notion)
}

func (m Model) onFetchMoreMilestones(msg notion.FetchMoreMilestonesMsg) (Model, tea.Cmd) {
	group := m.groups[msg.Status]
	if group.NextCursor != nil && !group.Loading {
		cursor := *group.NextCursor
		group.Loading = true
		m.groups[msg.Status] = group
		m.list.SetItems(m.buildMilestoneList())
		return m, m.notion.QueryMilestones(m.projID, msg.Status, cursor)
	}
	return m, nil
}

// todo: clean, super bloated
func (m Model) onFetchMoreTasksByStatus(msg notion.FetchMoreTasksMsg) (Model, tea.Cmd) {
	item := m.list.SelectedItem()
	if mstone, ok := item.(DefaultItem); ok {
		group := mstone.TaskGroups[msg.Status]

		if group.NextCursor != nil && !group.Loading {
			idx := m.list.Index()
			cursor := *group.NextCursor
			group.Loading = true
			mstone.TaskGroups[msg.Status] = group
			m.list.SetItem(idx, mstone)

			return m, tea.Batch(
				m.notion.QueryTasks(mstone.ID, msg.Status, cursor, idx),
				refreshMilestoneTasks(mstone.TaskGroups),
			)
		}
	}

	return m, nil
}

func (m Model) onToggleTaskGroup(msg notion.ToggleTaskGroupMsg) (Model, tea.Cmd) {
	item := m.list.SelectedItem()
	if mstone, ok := item.(DefaultItem); ok {
		group := mstone.TaskGroups[msg.Status]
		group.Hide = !group.Hide
		mstone.TaskGroups[msg.Status] = group
		m.updateMilestoneInGroups(mstone)
		return m, refreshMilestoneTasks(mstone.TaskGroups)
	}
	return m, nil
}

func (m Model) onTaskQuery(msg notion.TaskQueryMsg) (Model, tea.Cmd) {
	if msg.Err != nil {
		m.err = msg.Err
		return m, nil
	}
	item := m.list.Items()[msg.MilestoneIdx]
	if mstone, ok := item.(DefaultItem); ok {
		group := mstone.TaskGroups[msg.Status]
		group.Tasks = append(group.Tasks, msg.Pages...)
		group.NextCursor = msg.NextCursor
		group.Loading = false
		mstone.TaskGroups[msg.Status] = group
		mstone.FetchStatus = FetchSuccess
		m.list.SetItem(msg.MilestoneIdx, mstone)
		return m, refreshMilestoneTasks(mstone.TaskGroups)
	}
	return m, nil
}
