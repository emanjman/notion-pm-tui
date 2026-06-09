package objective

import tea "github.com/charmbracelet/bubbletea"

func (m Model) handleWindow(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	var mstoneCmd, taskCmd tea.Cmd
	// each panel is wrapped in a rounded border (+1 col each side, +1 row each side)
	const borderColsPerPanel = 2
	const borderRowsPerPanel = 2
	// each panel has inner padding set in View() (+1 col each side)
	const paddingColsPerPanel = 2
	const panelCount = 2

	totalColOverhead := (borderColsPerPanel + paddingColsPerPanel) * panelCount
	totalRowOverhead := borderRowsPerPanel // same overhead applies to both panels

	availableWidth := msg.Width - totalColOverhead
	leftWidth := availableWidth * 30 / 100
	rightWidth := availableWidth - leftWidth

	m.milestone, mstoneCmd = m.milestone.Update(tea.WindowSizeMsg{
		Width:  leftWidth,
		Height: msg.Height - totalRowOverhead - 1,
	})
	m.task, taskCmd = m.task.Update(tea.WindowSizeMsg{
		Width:  rightWidth,
		Height: msg.Height - totalRowOverhead - 1,
	})

	return m, tea.Batch(mstoneCmd, taskCmd)
}
