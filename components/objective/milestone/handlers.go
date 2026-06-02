package milestone

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"notion-project-tui/notion"
)

// handle received milestone-pages; re-render when all received
func (m Model) onMilestonePages(msg notion.MilestonePagesMsg) (Model, tea.Cmd) {
	m.pendingFetches--

	if msg.Err != nil {
		m.err = msg.Err
		return m, nil
	}

	// push incoming milestones to their group
	g := m.groups[msg.Status]
	g.Milestones = append(g.Milestones, msg.Pages...)
	g.NextCursor = msg.NextCursor
	m.groups[msg.Status] = g

	// only render milestones + kick off task-fetches after last batch
	if m.pendingFetches > 0 {
		return m, nil
	}
	m.list.SetItems(m.buildMilestoneList())
	return m, fetchInitMilestoneTasks(&m.list, m.notion)
}

// fetch more milestones for group queried
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

// ----------------------------------------------------------------------------

// fetch more tasks if exists
func (m Model) onFetchMoreTasksByStatus(msg notion.FetchMoreTasksMsg) (Model, tea.Cmd) {
	selected := m.list.SelectedItem()

	if mstone, ok := selected.(DefaultItem); ok {
		group := mstone.TaskGroups[msg.Status]

		if group.NextCursor != nil && !group.Loading {
			idx := m.list.Index()
			cursor := *group.NextCursor
			group.Loading = true
			mstone.TaskGroups[msg.Status] = group
			m.list.SetItem(idx, mstone)

			// get tasks + refresh
			return m, tea.Batch(
				m.notion.QueryTasks(mstone.ID, msg.Status, cursor, idx),
				refreshMilestoneTasks(mstone.TaskGroups),
			)
		}
	}

	return m, nil
}

// show/hide task group
func (m Model) onToggleTaskGroup(msg notion.ToggleTaskGroupMsg) (Model, tea.Cmd) {
	selected := m.list.SelectedItem()
	if mstone, ok := selected.(DefaultItem); ok {
		group := mstone.TaskGroups[msg.Status]
		group.Hide = !group.Hide
		mstone.TaskGroups[msg.Status] = group
		m = m.updateMilestoneInGroups(mstone)
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

// ----------------------------------------------------------------------------

// reconcile the optimistic milestone with the page notion actually created:
// on success swap the temp id for the real one; on failure drop the temp item.
func (m Model) onAddMilestonePage(msg notion.AddMilestonePageMsg) (Model, tea.Cmd) {
	if msg.Err != nil {
		log.Printf("add milestone failed: %v", msg.Err)
		return m.removeMilestoneByID(msg.TempID), nil
	}

	for status, group := range m.groups {
		for i, pg := range group.Milestones {
			if pg.ID == msg.TempID {
				group.Milestones[i].ID = msg.Page.ID
				m.groups[status] = group
				m.list.SetItems(m.buildMilestoneList())

				// keep focus tracking pointed at the real id
				if m.Focus.milestoneID == msg.TempID {
					m.Focus.milestoneID = msg.Page.ID
				}
				return m, nil
			}
		}
	}
	return m, nil
}

// if update failed, revert optimistic ui update to og stashed title
func (m Model) onUpdateNotionTitle(msg UpdateNotionTitleMsg) (Model, tea.Cmd) {
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
