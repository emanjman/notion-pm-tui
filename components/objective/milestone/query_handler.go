package milestone

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"notion-project-tui/notion"
)

// handle received milestone-pages; re-render when all fetches resolved
func (m Model) onQueryMilestonePages(msg notion.QueryMilestonePagesMsg) (Model, tea.Cmd) {
	m.pendingFetches--

	// handle failed fetch
	if msg.Err != nil {
		log.Printf("Error: %s", msg.Err)
		m.err = msg.Err
		return m, nil
	}

	// push incoming milestones to respective status-group
	g := m.groups[msg.Status]
	g.Milestones = append(g.Milestones, msg.Pages...)
	g.NextCursor = msg.NextCursor
	m.groups[msg.Status] = g

	// still awaiting milestone-page fetches, guard against re-render
	if m.pendingFetches > 0 {
		return m, nil
	}

	// re-render + fetch init set of milestone-tasks
	m.list.SetItems(m.buildMilestoneList())
	return m, fetchInitMilestoneTasks(&m.list, m.notion)
}

// fetch more milestones for group queried
func (m Model) onQueryMoreMilestonePages(msg notion.QueryMoreMilestonePagesMsg) (Model, tea.Cmd) {
	g := m.groups[msg.Status]

	// not state-ready
	if g.NextCursor == nil || g.Loading {
		return m, nil
	}

	g.Loading = true
	m = m.updateGroup(msg.Status, g)

	return m, m.notion.QueryMilestonePages(m.versionID, msg.Status, *g.NextCursor)
}

func (m Model) onQueryTaskPages(msg notion.QueryTaskPagesMsg) (Model, tea.Cmd) {
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

// fetch more tasks if exists
func (m Model) onQueryMoreTaskPages(msg notion.QueryMoreTaskPagesMsg) (Model, tea.Cmd) {
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
