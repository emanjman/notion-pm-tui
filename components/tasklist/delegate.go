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

type TaskListDelegate struct {
	focused bool
}

func NewTaskListDelegate(flag bool) TaskListDelegate {
	return TaskListDelegate{focused: flag}
}

func (d TaskListDelegate) Height() int                               { return 2 }
func (d TaskListDelegate) Spacing() int                              { return 0 }
func (d TaskListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d TaskListDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	selected := index == m.Index() && d.focused

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
		baseStyle := d.taskBaseStyle(selected)

		typ := baseStyle.Foreground(styles.MutedForeground).Render(item.Type)
		space := baseStyle.Render(" ")
		task := baseStyle.Foreground(styles.PrimaryForeground).Render(item.Task)

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

		primaryContent := typ + space + task
		secondaryContent := baseStyle.Foreground(priorityColors[p]).Render(fmt.Sprintf("%d", p))

		// padding width (account for task style's horizontal padding)
		frameStyle := d.taskFrameStyle(selected)
		availableWidth := m.Width() - frameStyle.GetHorizontalPadding()
		paddingWidth := availableWidth - lg.Width(primaryContent) - lg.Width(secondaryContent)
		if paddingWidth < 0 {
			paddingWidth = 0
		}

		// build full content with background applied to all segments
		padding := baseStyle.Render(strings.Repeat(" ", paddingWidth))
		content := primaryContent + padding + secondaryContent

		// wrap with frame style (border, outer padding, background if selected)
		fmt.Fprint(w, frameStyle.Width(m.Width()).Render(content))
	}
}

// -- styling

// finishing touches that applies borders/padding + additional highlights
func (d TaskListDelegate) taskFrameStyle(selected bool) lg.Style {
	s := lg.NewStyle().
		Border(lg.NormalBorder(), false, false, true, false).
		BorderForeground(styles.BorderForeground).
		PaddingLeft(4).
		PaddingRight(4)

	if selected {
		return s.Background(styles.SelectedBackground)
	}
	return s
}

// set base style that must apply to inner content
func (d TaskListDelegate) taskBaseStyle(selected bool) lg.Style {
	s := lg.NewStyle()
	if selected {
		return s.Background(styles.SelectedBackground)
	}
	return s
}

func (d TaskListDelegate) headerStyle(selected bool) lg.Style {
	s := lg.NewStyle().
		PaddingBottom(1).
		PaddingLeft(2).
		PaddingRight(2).
		Foreground(styles.MutedForeground)

	if selected {
		return s.Foreground(styles.PrimaryForeground).Underline(true)
	}
	return s
}
