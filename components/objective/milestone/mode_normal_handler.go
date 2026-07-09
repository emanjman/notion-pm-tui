package milestone

import (
	"fmt"
	"log"
	"notion-project-tui/notion"

	tea "github.com/charmbracelet/bubbletea"
)

// nav down 1 + refresh
func (m Model) handleNormalDown() (Model, tea.Cmd) {
	m.list.CursorDown()
	return m, refreshMilestoneTasks(m.getCurrMilestoneID(), m.getCurrTaskGroups())
}

// nav up 1 + refresh
func (m Model) handleNormalUp() (Model, tea.Cmd) {
	m.list.CursorUp()
	return m, refreshMilestoneTasks(m.getCurrMilestoneID(), m.getCurrTaskGroups())
}

// nav down 5 + refresh
func (m Model) handleNormalJumpDown() (Model, tea.Cmd) {
	m.list.Select(min(len(m.list.Items())-1, m.list.Index()+5))
	return m, refreshMilestoneTasks(m.getCurrMilestoneID(), m.getCurrTaskGroups())
}

// nav up 5 + refresh
func (m Model) handleNormalJumpUp() (Model, tea.Cmd) {
	m.list.Select(max(0, m.list.Index()-5))
	return m, refreshMilestoneTasks(m.getCurrMilestoneID(), m.getCurrTaskGroups())
}

// handle select behavior based on item type
func (m Model) handleNormalSelect() (Model, tea.Cmd) {
	selected := m.list.SelectedItem()

	if header, ok := selected.(GroupHeaderItem); ok {
		// toggle show/hide group
		g := m.groups[header.Status]
		g.Hide = !g.Hide
		m = m.updateGroup(header.Status, g)

	} else if loadMore, ok := selected.(LoadMoreItem); ok && !loadMore.Loading {
		emitQueryMoreMilestonePages(loadMore.Status)

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
			return m, refreshMilestoneTasks(mstone.ID, mstone.TaskGroups)
		}
	}

	return m, nil
}

// enter edit-mode
func (m Model) handleNormalRename() (Model, tea.Cmd) {
	if mstone, ok := m.list.SelectedItem().(DefaultItem); ok {
		// id milestone to update
		m.Edit.milestoneID = mstone.ID
		m.Edit.milestoneIdx = m.list.Index()

		// setup text input model
		m.Edit.titleInput = m.Edit.newTitleInput(mstone)

		// flip to edit-mode
		m = m.switchMode(EditMode)
	}

	return m, nil
}

func (m Model) handleNormalAdd() (Model, tea.Cmd) {
	m.Edit.tempIDs++
	tempID := fmt.Sprintf("temp-%d", m.Edit.tempIDs)

	newPage := notion.MilestonePage{
		ID: tempID,
		// Icon: notion.Icon{Type: notion.IconEmoji, Emoji: ""},
		Properties: notion.MilestoneProperties{
			Title: notion.TitleProperty{
				Title: []notion.RichText{{Text: notion.TextContent{Content: ""}}},
			},
		},
	}

	g := m.groups[notion.MilestoneIdle]
	g.Milestones = append(g.Milestones, newPage)
	m = m.updateGroup(notion.MilestoneIdle, g)

	// find new milestone idx in list; enter edit-mode
	for i, item := range m.list.Items() {
		if mstone, ok := item.(DefaultItem); ok && mstone.ID == tempID {
			// focus on mstone
			m.list.Select(i)

			m.Edit.milestoneID = tempID
			m.Edit.milestoneIdx = i
			m.Edit.titleInput = m.Edit.newTitleInput(mstone)

			// flip to edit-mode
			m = m.switchMode(EditMode)
		}
	}

	return m, nil
}

// store mstone + switch to delete-mode; await user decision
func (m Model) handleNormalDelete() (Model, tea.Cmd) {
	log.Printf("on neutral delete") // !debug
	selected := m.list.SelectedItem()
	if cur, ok := selected.(DefaultItem); ok && cur.TaskCount == 0 {
		for _, mstone := range m.groups[cur.MilestoneStatus].Milestones {
			if mstone.ID == cur.ID {
				m.Delete.milestoneBackup = mstone
				break
			}
		}
		m = m.switchMode(DeleteMode)
	}
	return m, nil
}
