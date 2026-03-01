package tasklist

import (
	"fmt"
	"io"
	"notion-project-tui/styles"
	listutil "notion-project-tui/util/list"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type TaskListDelegate struct{}

func NewTaskListDelegate() TaskListDelegate {
	return TaskListDelegate{}
}

func (d TaskListDelegate) Height() int                               { return 2 }
func (d TaskListDelegate) Spacing() int                              { return 0 }
func (d TaskListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d TaskListDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	selected := index == m.Index()

	switch item := item.(type) {
	case listutil.ListItemGroupHeader:
		chevron := "▼"
		if item.Hidden {
			chevron = "▶"
		}

		content := fmt.Sprintf("%s %s (%d)", chevron, item.Label, item.Count)

		style := d.headerStyle(selected).Width(m.Width())
		fmt.Fprint(w, style.Render(content))

	case TaskListItem:
		// Get content style with background (if selected)
		bg := d.taskBaseStyle(selected)

		typ := bg.Foreground(styles.MutedForeground).Render(item.Type)
		space := bg.Render(" ")
		task := bg.Foreground(styles.PrimaryForeground).Render(item.Task)

		priorityColors := []lg.Color{
			styles.MutedForeground, // p0 - none/gray
			lg.Color("#7aa2f7"),    // p1 - blue (low, calm)
			lg.Color("#9ece6a"),    // p2 - green (medium-low)
			lg.Color("#e0af68"),    // p3 - yellow (medium, caution)
			lg.Color("#ff9e64"),    // p4 - orange (high, warning)
			lg.Color("#f7768e"),    // p5 - red (critical, urgent)
		}
		p := item.Priority
		if p < 0 || p >= len(priorityColors) {
			p = 0
		}
		secondaryContent := bg.Foreground(priorityColors[p]).Render(fmt.Sprintf("%d", p))

		primaryContent := typ + space + task

		// padding width (account for task style's horizontal padding)
		taskStyle := d.taskFrameStyle(selected)
		availableWidth := m.Width() - taskStyle.GetHorizontalPadding()
		paddingWidth := availableWidth - lg.Width(primaryContent) - lg.Width(secondaryContent)
		if paddingWidth < 0 {
			paddingWidth = 0
		}

		// build full content with background applied to all segments
		padding := bg.Render(strings.Repeat(" ", paddingWidth))
		content := primaryContent + padding + secondaryContent

		// wrap with task style (border, outer padding, background if selected)
		fmt.Fprint(w, taskStyle.Width(m.Width()).Render(content))
	}
}

// -- styling

// finishing touches that applies borders/padding + additional highlights
func (d TaskListDelegate) taskFrameStyle(selected bool) lg.Style {
	base := lg.NewStyle().
		Border(lg.NormalBorder(), false, false, true, false).
		BorderForeground(styles.BorderForeground).
		PaddingLeft(4).
		PaddingRight(4)

	if selected {
		return base.Background(styles.SelectedBackground)
	}
	return base
}

// set base style that must apply to inner content
func (d TaskListDelegate) taskBaseStyle(selected bool) lg.Style {
	base := lg.NewStyle()
	if selected {
		return base.Background(styles.SelectedBackground)
	}
	return base
}

func (d TaskListDelegate) headerStyle(selected bool) lg.Style {
	base := lg.NewStyle().
		PaddingBottom(1).
		PaddingLeft(2).
		PaddingRight(2)

	if selected {
		return base.Underline(true)
	}
	return base
}
