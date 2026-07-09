package milestone

import tea "github.com/charmbracelet/bubbletea"

// return to neutral-mode
func (m Model) handleDeleteCancel() (Model, tea.Cmd) {
	m = m.switchMode(NeutralMode)
	return m, nil
}

// send off req(s) to delete milestone pg + dependent task pg(s);
// should await for all req's before reacting
func (m Model) handleDeleteConfirm() (Model, tea.Cmd) {
	// optimistically remove milestone from ui
	mstone := m.Delete.milestoneBackup
	m = m.removeMilestoneByID(mstone.ID)

	err := m.notion.TrashPage(mstone.ID)
	if err != nil {
		return m, emitTrashMilestonePage(err)
	}

	m = m.switchMode(NeutralMode)
	return m, nil
}
