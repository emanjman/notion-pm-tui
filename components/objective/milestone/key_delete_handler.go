package milestone

import tea "github.com/charmbracelet/bubbletea"

// return to neutral-mode
func (m Model) onDeleteCancel() (Model, tea.Cmd) {
	m = m.switchMode(NeutralMode)
	return m, nil
}

// send off req(s) to delete milestone pg + dependent task pg(s);
// should await for all req's before reacting
func (m Model) onDeleteConfirm() (Model, tea.Cmd) {
	// optimistically remove milestone from ui

	// todo: get taskIDs
	// todo: get milestoneID
	// todo: set `pending` tasks to await
	// todo: kickoff tea batch req

	return m, nil
}
