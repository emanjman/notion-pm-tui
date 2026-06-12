package objective

import (
	"notion-project-tui/styles"

	lg "github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	var (
		vstyle = lg.NewStyle().Border(lg.RoundedBorder(), true).Padding(0, 1)
		mstyle = lg.NewStyle().Border(lg.RoundedBorder(), true).Padding(0, 1)
		tstyle = lg.NewStyle().Border(lg.RoundedBorder(), true).Padding(0, 1)
	)

	var (
		on  = styles.MutedForeground
		off = styles.BorderForeground
	)

	switch m.panel {
	case VersionPanel:
		mstyle = mstyle.BorderForeground(off)
		tstyle = tstyle.BorderForeground(off)
	case MilestonePanel:
		mstyle = mstyle.BorderForeground(on)
		tstyle = tstyle.BorderForeground(off)
	case TaskPanel:
		mstyle = mstyle.BorderForeground(off)
		tstyle = tstyle.BorderForeground(on)
	}

	version := vstyle.Render(m.version.View())
	mstone := mstyle.Render(m.milestone.View())
	task := tstyle.Render(m.task.View())

	primary := lg.JoinHorizontal(lg.Left, mstone, task)
	return lg.JoinVertical(lg.Center, version, primary)
}
