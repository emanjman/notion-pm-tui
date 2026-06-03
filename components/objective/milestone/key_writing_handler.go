package milestone

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// update name on ui + notion server
func (m Model) onWritingSave() (Model, tea.Cmd) {
	var cmd tea.Cmd

	if mstone, ok := m.list.SelectedItem().(DefaultItem); ok {
		isNew := strings.HasPrefix(mstone.ID, "temp")
		title := strings.TrimSpace(m.Focus.tempTitle.Value())

		// stash og title (in case revert needed on failed server update)
		m.Focus.prevTitle = mstone.Name

		// optimistically update local milestone title
		mstone.Name = m.Focus.tempTitle.Value()
		m.list.SetItem(m.list.Index(), mstone)
		m = m.updateMilestone(mstone)

		switch {
		case isNew && title != "":
			// create on notion; temp id gets reconciled on the response
			cmd = m.notion.AddMilestone(mstone.ID, title)
		case isNew:
			// discard an empty brand-new milestone instead of persisting it
			m = m.removeMilestoneByID(mstone.ID)
		default:
			// existing milestone: push the title update
			cmd = updateNotionMilestoneTitle(m.notion, mstone.ID, mstone.Name)
		}
	}

	m.ActiveKeyMap = NeutralKeyMapper
	m.Focus.Mode = NeutralMode

	return m, cmd
}
