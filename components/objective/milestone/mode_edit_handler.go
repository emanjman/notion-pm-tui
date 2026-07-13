package milestone

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// update name on ui + notion server
func (m Model) handleEditSave() (Model, tea.Cmd) {
	var cmd tea.Cmd

	if mstone, ok := m.list.SelectedItem().(DefaultItem); ok {
		isNew := strings.HasPrefix(mstone.ID, "temp")
		title := strings.TrimSpace(m.Edit.titleInput.Value())

		// stash og title (in case revert needed on failed server update)
		m.Edit.titleBackup = mstone.Name

		// optimistically update local milestone title
		mstone.Name = m.Edit.titleInput.Value()
		m.list.SetItem(m.list.Index(), mstone)
		m = m.updateMilestone(mstone)

		switch {
		case isNew && title != "":
			// create on notion under the active version; temp id gets reconciled on the response
			cmd = m.notion.AddMilestonePage(mstone.ID, title, m.versionID)
		case isNew:
			// discard an empty brand-new milestone instead of persisting it
			m = m.removeMilestoneByID(mstone.ID)
		default:
			// existing milestone: push the title update
			cmd = updateNotionMilestoneTitle(m.notion, mstone.ID, mstone.Name)
		}
	}

	m = m.switchMode(NormalMode)

	return m, cmd
}
