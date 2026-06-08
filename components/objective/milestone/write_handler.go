package milestone

import (
	"log"
	"notion-project-tui/notion"

	tea "github.com/charmbracelet/bubbletea"
)

// reconcile the optimistic milestone-creation w/ result of actual notion-page creation
func (m Model) onAddMilestonePage(msg notion.AddMilestonePageMsg) (Model, tea.Cmd) {
	if msg.Err != nil {
		log.Printf("Add milestone failed: %v", msg.Err)
		return m.removeMilestoneByID(msg.TempID), nil
	}

	for status, group := range m.groups {
		for i, pg := range group.Milestones {
			if pg.ID == msg.TempID {
				group.Milestones[i].ID = msg.Page.ID
				m = m.updateGroup(status, group)

				// keep focus tracking pointed at the real id
				if m.Edit.milestoneID == msg.TempID {
					m.Edit.milestoneID = msg.Page.ID
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
			if mstone, ok := item.(DefaultItem); ok && mstone.ID == m.Edit.milestoneID {
				mstone.Name = m.Edit.titleBackup
				m.list.SetItem(i, mstone)
				m = m.updateMilestone(mstone)
				break
			}
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
		m = m.updateMilestone(mstone)
		return m, refreshMilestoneTasks(mstone.TaskGroups)
	}
	return m, nil
}

// if trash-page failed, restore optimistic ui deletion
func (m Model) onTrashMilestonePage(msg TrashMilestonePageMsg) (Model, tea.Cmd) {
	if msg.Err != nil {
		m.err = msg.Err
		status, err := notion.MilestoneStatusFromString(m.Delete.milestoneBackup.Properties.Status)
		g := m.groups[status]
		g.Milestones = append(g.Milestones, m.Delete.milestoneBackup)
		m = m.updateGroup(status, g)
	}
	return m, nil
}
