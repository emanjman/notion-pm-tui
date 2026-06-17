package objective

import (
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

func (m Model) handleWindow(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	var versionCmd, mstoneCmd, taskCmd tea.Cmd

	// each panel is wrapped in a rounded border (+1 col each side, +1 row each side)
	const panelBorderWidth = 2
	const panelBorderHeight = 2
	// each panel has inner padding set in View() (+1 col each side)
	const panelPaddingWidth = 2
	const panelCount = 2

	totalWidthOverhead := (panelBorderWidth + panelPaddingWidth) * panelCount
	totalHeightOverhead := panelBorderHeight // same overhead applies to both panels

	availableWidth := msg.Width - totalWidthOverhead
	leftWidth := availableWidth * 30 / 100
	rightWidth := availableWidth - leftWidth

	// todo: is this circular?
	versionHeight := lg.NewStyle().SetString(m.version.View()).GetHeight()

	m.version, versionCmd = m.version.Update(tea.WindowSizeMsg{
		Width:  availableWidth,
		Height: versionHeight,
	})
	m.milestone, mstoneCmd = m.milestone.Update(tea.WindowSizeMsg{
		Width:  leftWidth,
		Height: msg.Height - totalHeightOverhead - 1 - versionHeight,
	})
	m.task, taskCmd = m.task.Update(tea.WindowSizeMsg{
		Width:  rightWidth,
		Height: msg.Height - totalHeightOverhead - 1 - versionHeight,
	})

	return m, tea.Batch(versionCmd, mstoneCmd, taskCmd)
}
