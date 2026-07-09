package objective

import (
	"notion-project-tui/styles"

	lg "github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	// pin each panel to its computed dims so an empty child list doesn't collapse
	const padW = 2

	var (
		vstyle = lg.NewStyle().Border(lg.RoundedBorder(), true).Padding(0, 1).
			Width(m.versionWidth + padW)
		mstyle = lg.NewStyle().Border(lg.RoundedBorder(), true).Padding(0, 1).
			Width(m.milestoneWidth + padW).Height(m.panelHeight)
		tstyle = lg.NewStyle().Border(lg.RoundedBorder(), true).Padding(0, 1).
			Width(m.taskWidth + padW).Height(m.panelHeight)
	)

	var (
		on  = styles.MutedForeground
		off = styles.BorderForeground
	)

	switch m.panel {
	case VersionPanel:
		vstyle = vstyle.BorderForeground(on)
		mstyle = mstyle.BorderForeground(off)
		tstyle = tstyle.BorderForeground(off)
	case MilestonePanel:
		vstyle = vstyle.BorderForeground(off)
		mstyle = mstyle.BorderForeground(on)
		tstyle = tstyle.BorderForeground(off)
	case TaskPanel:
		vstyle = vstyle.BorderForeground(off)
		mstyle = mstyle.BorderForeground(off)
		tstyle = tstyle.BorderForeground(on)
	}

	version := vstyle.Render(m.version.View())
	mstone := mstyle.Render(m.milestone.View())
	task := tstyle.Render(m.task.View())

	primary := lg.JoinHorizontal(lg.Left, mstone, task)
	return lg.JoinVertical(lg.Center, version, primary)
}
