package objective

import tea "github.com/charmbracelet/bubbletea"

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	case tea.WindowSizeMsg:
		return m.handleWindow(msg)
	}

	// non-key messages (data, cmds) go to both panels
	var milestoneCmd, taskCmd tea.Cmd
	m.milestone, milestoneCmd = m.milestone.Update(msg)
	m.task, taskCmd = m.task.Update(msg)
	return m, tea.Batch(milestoneCmd, taskCmd)
}
