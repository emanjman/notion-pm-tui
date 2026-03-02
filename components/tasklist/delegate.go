package tasklist

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"io"
	"notion-project-tui/styles"
	listutil "notion-project-tui/util/list"
)

type variantStyle struct {
	base     lg.Style
	selected lg.Style
}

type style struct {
	itemContainer variantStyle
	itemSegment   variantStyle
	header        variantStyle
}

type TaskListDelegate struct {
	focused bool
	style   style
}

func NewTaskListDelegate(focused bool) TaskListDelegate {
	// item container style
	var (
		icbase = lg.NewStyle().
			Border(lg.NormalBorder(), false, false, true, false).
			BorderForeground(styles.BorderForeground).
			PaddingLeft(4).
			PaddingRight(4)
		icsel = icbase.
			Background(styles.SelectedBackground)
	)

	// item segment style
	var (
		isbase = lg.NewStyle()
		issel  = isbase.
			Background(styles.SelectedBackground)
	)

	// header style
	var (
		hbase = lg.NewStyle().
			PaddingBottom(1).
			PaddingLeft(2).
			PaddingRight(2).
			Foreground(styles.MutedForeground)
		hsel = hbase.
			Foreground(styles.PrimaryForeground).
			Underline(true)
	)

	return TaskListDelegate{
		focused: focused,
		style: style{
			itemContainer: variantStyle{base: icbase, selected: icsel},
			itemSegment:   variantStyle{base: isbase, selected: issel},
			header:        variantStyle{base: hbase, selected: hsel},
		},
	}
}

func (d TaskListDelegate) Height() int                               { return 2 }
func (d TaskListDelegate) Spacing() int                              { return 0 }
func (d TaskListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d TaskListDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	selected := index == m.Index() && d.focused

	switch item := item.(type) {
	case listutil.ListItemGroupHeader:
		style := d.style.header.base
		if selected {
			style = d.style.header.selected
		}

		chevron := "▼"
		if item.Hidden {
			chevron = "▶"
		}

		content := fmt.Sprintf("%s %s (%d)", chevron, item.Label, item.Count)

		fmt.Fprint(w, style.Width(m.Width()).Render(content))

	case TaskListItem:
		segStyle := d.style.itemSegment.base
		contStyle := d.style.itemContainer.base
		if selected {
			segStyle = d.style.itemSegment.selected
			contStyle = d.style.itemContainer.selected
		}

		typ := segStyle.Foreground(styles.MutedForeground).Render(item.Type)
		space := segStyle.Render(" ")
		task := segStyle.Foreground(styles.PrimaryForeground).Render(item.Task)

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

		left := typ + space + task
		right := segStyle.Foreground(priorityColors[p]).Render(fmt.Sprintf("[%d]", p))

		px := styles.GetPaddingBetween(left, right, m.Width(), contStyle)
		content := left + styles.RenderPadding(segStyle, px) + right

		// wrap with frame style (border, outer padding, background if selected)
		fmt.Fprint(w, contStyle.Width(m.Width()).Render(content))
	}
}
