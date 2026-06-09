package objective

import (
	"notion-project-tui/styles"

	lg "github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	leftStyle := lg.NewStyle().Border(lg.RoundedBorder(), true).Padding(0, 1)
	rightStyle := lg.NewStyle().Border(lg.RoundedBorder(), true).Padding(0, 1)

	if m.focus == MilestonePanel {
		leftStyle = leftStyle.BorderForeground(styles.MutedForeground)
		rightStyle = rightStyle.BorderForeground(styles.BorderForeground)
	} else {
		rightStyle = rightStyle.BorderForeground(styles.MutedForeground)
		leftStyle = leftStyle.BorderForeground(styles.BorderForeground)
	}

	left := leftStyle.Render(m.milestone.View())
	right := rightStyle.Render(m.task.View())
	return lg.JoinHorizontal(lg.Top, left, right)
}
