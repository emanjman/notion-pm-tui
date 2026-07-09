package objective

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleWindow(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	var versionCmd, mstoneCmd, taskCmd tea.Cmd

	// each panel has a border, holding 1space on each side
	const panelBorderWidth = 2
	const panelBorderHeight = 2

	const panelPaddingWidth = 2
	const panelCount = 2

	totalWidthOverhead := (panelBorderWidth + panelPaddingWidth) * panelCount
	totalHeightOverhead := panelBorderHeight // same overhead applies to both panels

	availableWidth := msg.Width - totalWidthOverhead
	leftWidth := availableWidth * 30 / 100
	rightWidth := availableWidth - leftWidth

	versionHeight := 1
	versionWidth := msg.Width - (panelBorderHeight + panelPaddingWidth)
	panelHeight := msg.Height - totalHeightOverhead - 1 - versionHeight - panelBorderHeight

	// log.Printf("%v %v %v %v %v %v", totalWidthOverhead, totalHeightOverhead, availableWidth, leftWidth, rightWidth, versionHeight)

	m.versionWidth = versionWidth
	m.milestoneWidth = leftWidth
	m.taskWidth = rightWidth
	m.panelHeight = panelHeight

	m.version, versionCmd = m.version.Update(tea.WindowSizeMsg{
		Width:  versionWidth,
		Height: versionHeight,
	})

	m.milestone, mstoneCmd = m.milestone.Update(tea.WindowSizeMsg{
		Width:  leftWidth,
		Height: panelHeight,
	})
	m.task, taskCmd = m.task.Update(tea.WindowSizeMsg{
		Width:  rightWidth,
		Height: panelHeight,
	})

	return m, tea.Batch(versionCmd, mstoneCmd, taskCmd)
}
