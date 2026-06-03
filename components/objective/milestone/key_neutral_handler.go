package milestone

import (
	"fmt"
	"notion-project-tui/notion"

	tea "github.com/charmbracelet/bubbletea"
)

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

// handle select behavior based on item type
func (m Model) onNeutralSelect() (Model, tea.Cmd) {
	selected := m.list.SelectedItem()

	if header, ok := selected.(GroupHeaderItem); ok {
		// toggle show/hide group
		g := m.groups[header.Status]
		g.Hide = !g.Hide
		m = m.updateGroup(header.Status, g)

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

func (m Model) onNeutralAdd() (Model, tea.Cmd) {
	m.tempIDCounter++
	tempID := fmt.Sprintf("temp-%d", m.tempIDCounter)

	newMilestonepage := notion.MilestonePage{
		ID: tempID,
		// Icon: notion.Icon{Type: notion.IconEmoji, Emoji: ""},
		Properties: notion.MilestoneProperties{
			Title: notion.TitleProperty{
				Title: []notion.RichText{{Text: notion.TextContent{Content: ""}}},
			},
		},
	}

	g := m.groups[notion.MilestoneIdle]
	g.Milestones = append(g.Milestones, newMilestonepage)
	m = m.updateGroup(notion.MilestoneIdle, g)

	// find new milestone idx in list; enter writing-mode
	for i, item := range m.list.Items() {
		if mstone, ok := item.(DefaultItem); ok && mstone.ID == tempID {
			// focus on mstone
			m.list.Select(i)

			m.Focus.milestoneID = tempID
			m.Focus.milestoneIdx = i
			m.Focus.tempTitle = initTempTitle(mstone)

			// flip to writing-mode
			m.Focus.Mode = WritingMode
			m.ActiveKeyMap = WritingKeyMapper
		}
	}

	return m, nil
}
