package tasklist

import (
	"fmt"
	"io"
	listutil "notion-project-tui/util/list"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type delegateStyle struct {
	defaultStyle  lg.Style
	selectedStyle lg.Style
}

type TaskListDelegate struct {
	task   delegateStyle
	header delegateStyle
}

func NewTaskListDelegate() TaskListDelegate {
	taskBase := lg.NewStyle().
		Border(lg.NormalBorder(), false, false, true, false).
		BorderForeground(lg.Color("236")).
		PaddingLeft(4).
		PaddingRight(4)

	headerBase := lg.NewStyle().
		Border(lg.NormalBorder(), false, false, true, false).
		BorderForeground(lg.Color("236")).
		PaddingLeft(2).
		PaddingRight(2)

	return TaskListDelegate{
		task: delegateStyle{
			defaultStyle: taskBase,
			selectedStyle: taskBase.
				Foreground(lg.Color("205")),
		},

		header: delegateStyle{
			defaultStyle: headerBase,
			selectedStyle: headerBase.
				Foreground(lg.Color("205")),
		},
	}
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

		style := d.header.defaultStyle
		if selected {
			style = d.header.selectedStyle
		}
		style.Width(m.Width())

		fmt.Fprint(w, style.Render(content))

	case TaskListItem:
		var (
			task     = item.Task
			typename = item.Type
			priority = item.Priority
		)

		style := d.task.defaultStyle
		if selected {
			style = d.task.selectedStyle
		}

		left := fmt.Sprintf("[%s] %s", typename, task)
		right := fmt.Sprintf("%d", priority)
		content := padBetween(left, right, m.Width(), style)

		fmt.Fprint(w, style.Width(m.Width()).Render(content))
	}
}

func padBetween(left, right string, windowWidth int, style lg.Style) string {
	// use lg.Width to only consider visible cells
	padding := windowWidth - lg.Width(left) - lg.Width(right) - style.GetHorizontalPadding()
	if padding < 0 {
		padding = 0
	}

	return left + strings.Repeat(" ", padding) + right
}
