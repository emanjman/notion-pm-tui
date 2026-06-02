package milestone

import (
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// update name on ui + notion server
func (m Model) onWritingSave() (Model, tea.Cmd) {
	log.Printf("on writing save")
	var cmd tea.Cmd

	if mstone, ok := m.list.SelectedItem().(DefaultItem); ok {
		// stash og title (in case revert needed on failed server update)
		m.Focus.prevTitle = mstone.Name

		// optimistically update local milestone title
		mstone.Name = m.Focus.tempTitle.Value()
		m.list.SetItem(m.list.Index(), mstone)
		m = m.updateMilestoneInGroups(mstone)

		if strings.HasPrefix(mstone.ID, "temp") && strings.TrimSpace(mstone.Name) != "" {
			// todo: create milestone
		}

		// set cmd to send-off notion update
		cmd = updateNotionMilestoneTitle(m.notion, mstone.ID, mstone.Name)
	}

	m.ActiveKeyMap = NeutralKeyMapper
	m.Focus.Mode = NeutralMode

	return m, cmd
}
