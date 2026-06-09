package objective

import tea "github.com/charmbracelet/bubbletea"

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	case tea.WindowSizeMsg:
		return m.handleWindow(msg)
	}

	// key presses go to active panel only; data messages go to both
	var milestoneCmd, taskCmd tea.Cmd

	if _, isKey := msg.(tea.KeyMsg); isKey {
		switch m.focus {
		case MilestonePanel:
			m.milestone, milestoneCmd = m.milestone.Update(msg)
			return m, milestoneCmd
		case TaskPanel:
			m.task, taskCmd = m.task.Update(msg)
			return m, taskCmd
		}
	} else {
		m.milestone, milestoneCmd = m.milestone.Update(msg)
		m.task, taskCmd = m.task.Update(msg)
	}

	return m, tea.Batch(milestoneCmd, taskCmd)
}
