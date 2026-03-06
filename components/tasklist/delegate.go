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

	focus *FocusState
}

func NewTaskListDelegate(focused bool, focus *FocusState) TaskListDelegate {
	borderDistance := 0
	rightEdgeDistance := 3

	// item container style
	var (
		icbase = lg.NewStyle().
			Border(lg.NormalBorder(), false, false, true, false).
			BorderForeground(styles.BorderForeground).
			PaddingLeft(borderDistance + 2).
			PaddingRight(rightEdgeDistance)
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
			Foreground(styles.MutedForeground).
			PaddingLeft(borderDistance).
			PaddingRight(rightEdgeDistance)
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
		focus: focus,
	}
}

func (d TaskListDelegate) Height() int                               { return 2 }
func (d TaskListDelegate) Spacing() int                              { return 0 }
func (d TaskListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d TaskListDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	selected := index == m.Index() && d.focused

	switch item := item.(type) {
	case listutil.ListItemGroupHeader:
		header := renderListItemGroupHeader(d, item, selected, m.Width())
		fmt.Fprint(w, header)
	case TaskListItem:
		task := renderTaskListItem(d, item, selected, m.Width())
		fmt.Fprint(w, task)
	}
}

// -- helper funcs

func renderListItemGroupHeader(d TaskListDelegate, item listutil.ListItemGroupHeader, selected bool, windowWidth int) string {
	style := d.style.header.base
	if selected {
		style = d.style.header.selected
	}

	chevron := "▼"
	if item.Hidden {
		chevron = "▶"
	}

	content := fmt.Sprintf("%s %s (%d)", chevron, item.Label, item.Count)
	return style.Width(windowWidth).Render(content)
}

var priorityColors = []lg.Color{
	styles.MutedForeground, // p0 - none/gray
	lg.Color("#7aa2f7"),    // p1 - blue (low, calm)
	lg.Color("#9ece6a"),    // p2 - green (medium-low)
	lg.Color("#e0af68"),    // p3 - yellow (medium, caution)
	lg.Color("#ff9e64"),    // p4 - orange (high, warning)
	lg.Color("#f7768e"),    // p5 - red (critical, urgent)
}

func renderTaskListItem(d TaskListDelegate, item TaskListItem, selected bool, windowWidth int) string {
	contStyle := d.style.itemContainer.base
	segStyle := d.style.itemSegment.base
	typStyle, titleStyle, priorityStyle := lg.Style{}, lg.Style{}, lg.Style{}

	// handle field highlighting by mode
	if selected {
		if d.focus.Mode == NeutralMode {
			segStyle = d.style.itemSegment.selected
			contStyle = d.style.itemContainer.selected

			// apply select highlight row-wide
			typStyle, titleStyle, priorityStyle = segStyle, segStyle, segStyle
		} else {
			// apply select highlight by field
			switch d.focus.field {
			case TaskType:
				typStyle = typStyle.Inherit(d.style.itemSegment.selected)
			case TaskTitle:
				titleStyle = titleStyle.Inherit(d.style.itemSegment.selected)
			case TaskPriority:
				priorityStyle = priorityStyle.Inherit(d.style.itemSegment.selected)
			}
		}
	}

	// guard against unhandled priority values
	safePriorityIdx := item.Priority
	if safePriorityIdx < 0 || safePriorityIdx >= len(priorityColors) {
		safePriorityIdx = 0
	}

	// apply final field-specific styles
	typStyle = typStyle.Foreground(styles.MutedForeground)
	titleStyle = titleStyle.Foreground(styles.PrimaryForeground)
	priorityStyle = priorityStyle.Foreground(priorityColors[safePriorityIdx])

	// render each field
	typ := typStyle.Render(item.Type)
	space := segStyle.Render(" ")
	task := titleStyle.Render(item.Task)
	priority := priorityStyle.Render(fmt.Sprintf("[%d]", safePriorityIdx))

	left := typ + space + task
	right := priority

	px := styles.GetPaddingBetween(left, right, windowWidth, contStyle)
	content := left + styles.RenderPadding(segStyle, px) + right

	return contStyle.Width(windowWidth).Render(content)
}
