package milestone

import (
	"notion-project-tui/notion"

	tea "github.com/charmbracelet/bubbletea"
)

// update name on ui + notion server
func (m Model) onWritingModeSave() (Model, tea.Cmd) {
	var cmd tea.Cmd

	if mstone, ok := m.list.SelectedItem().(DefaultItem); ok {
		// stash og title (in case revert needed on failed server update)
		m.Focus.prevTitle = mstone.Name

		// optimistically update local milestone title
		mstone.Name = m.Focus.tempTitle.Value()
		m.list.SetItem(m.list.Index(), mstone)
		m = m.updateMilestoneInGroups(mstone)

		// set cmd to send-off notion update
		cmd = updateNotionMilestoneTitle(m.notion, mstone.ID, mstone.Name)
	}

	m.ActiveKeyMap = NeutralKeyMapper
	m.Focus.Mode = NeutralMode

	return m, cmd
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

// update milestones; re-render after all req's are received
func (m Model) onMilestonePages(msg notion.MilestonePagesMsg) (Model, tea.Cmd) {
	if msg.Err != nil {
		m.err = msg.Err
		m.pendingFetches--
		return m, nil
	}

	// push incoming milestones to their group
	g := m.groups[msg.Status]
	g.Milestones = append(g.Milestones, msg.Pages...)
	g.NextCursor = msg.NextCursor
	m.groups[msg.Status] = g

	// only render + kick off task-fetches after last batch
	m.pendingFetches--
	if m.pendingFetches > 0 {
		return m, nil
	}
	m.list.SetItems(m.buildMilestoneList())
	return m, fetchInitMilestoneTasks(&m.list, m.notion)
}

// if update failed, revert optimistic ui update to og stashed title
func (m Model) onUpdateTitle(msg UpdateTitleMsg) (Model, tea.Cmd) {
	if msg.Err != nil {
		for i, item := range m.list.Items() {
			if mstone, ok := item.(DefaultItem); ok && mstone.ID == m.Focus.milestoneID {
				mstone.Name = m.Focus.prevTitle
				m.list.SetItem(i, mstone)
				m = m.updateMilestoneInGroups(mstone)
				break
			}
		}
	}
	return m, nil
}

// todo: comprehend
func (m Model) onFetchMoreMilestones(msg notion.FetchMoreMilestonesMsg) (Model, tea.Cmd) {
	g := m.groups[msg.Status]
	if g.NextCursor != nil && !g.Loading {
		cursor := *g.NextCursor
		g.Loading = true
		m.groups[msg.Status] = g
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
