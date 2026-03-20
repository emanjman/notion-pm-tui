package task

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

type ItemDelegate struct {
	focused bool
	style   style

	focus *FocusState
}

func NewItemDelegate(focused bool, focus *FocusState) ItemDelegate {
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

	return ItemDelegate{
		focused: focused,
		style: style{
			itemContainer: variantStyle{base: icbase, selected: icsel},
			itemSegment:   variantStyle{base: isbase, selected: issel},
			header:        variantStyle{base: hbase, selected: hsel},
		},
		focus: focus,
	}
}

func (d ItemDelegate) Height() int                               { return 2 }
func (d ItemDelegate) Spacing() int                              { return 0 }
func (d ItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	selected := index == m.Index() && d.focused

	switch item := item.(type) {
	case listutil.ListItemGroupHeader:
		header := renderItemHeader(d, item, selected, m.Width())
		fmt.Fprint(w, header)
	case Item:
		task := renderItem(d, item, selected, m.Width())
		fmt.Fprint(w, task)
	}
}

// -- helper funcs

func renderItemHeader(d ItemDelegate, item listutil.ListItemGroupHeader, selected bool, windowWidth int) string {
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

type priority struct {
	fg       lg.Color
	severity string
}

var priorities = []priority{
	{fg: styles.MutedForeground, severity: "none"},
	{fg: lg.Color("#7aa2f7"), severity: "minimal"},
	{fg: lg.Color("#9ece6a"), severity: "low"},
	{fg: lg.Color("#e0af68"), severity: "medium"},
	{fg: lg.Color("#ff9e64"), severity: "high"},
	{fg: lg.Color("#f7768e"), severity: "critical"},
}

func renderItem(d ItemDelegate, item Item, selected bool, windowWidth int) string {
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

		// if deleting curr item, override w/ red bg
		if d.focus.pendingDelete && d.focus.taskID == item.ID {
			contStyle = contStyle.Background(styles.ErrorBackground)
			segStyle = segStyle.Background(styles.ErrorBackground)
			typStyle = typStyle.Background(styles.ErrorBackground)
			titleStyle = titleStyle.Background(styles.ErrorBackground)
			priorityStyle = priorityStyle.Background(styles.ErrorBackground)
		}
	}

	// guard against unhandled priority values
	safePriorityIdx := item.Priority
	if safePriorityIdx < 0 || safePriorityIdx >= len(priorities) {
		safePriorityIdx = 0
	}

	// apply final field-specific styles
	typStyle = typStyle.Foreground(styles.MutedForeground)
	titleStyle = titleStyle.Foreground(styles.PrimaryForeground)
	priorityStyle = priorityStyle.
		Foreground(priorities[safePriorityIdx].fg).
		Padding(0, 1)

	// render each field
	typ := typStyle.Render(item.Type)
	space := segStyle.Render(" ")

	// handle empty title with placeholder
	renderedTitle := item.Task
	if renderedTitle == "" {
		renderedTitle = "Enter task..."
		titleStyle = titleStyle.Foreground(styles.MutedForeground)
	}
	title := titleStyle.Render(renderedTitle)
	priority := priorityStyle.Render(fmt.Sprintf("%s [%d]", priorities[safePriorityIdx].severity, safePriorityIdx))

	// calculate max title width
	leftOffset, rightOffset, fieldSpacing := 3, 2, 2
	offset := leftOffset + rightOffset + fieldSpacing
	titleMaxWidth := windowWidth - lg.Width(typ) - lg.Width(priority) - (offset)

	if selected && d.focus.Mode == WritingMode {
		// use textinput component in writing mode
		d.focus.tempTitle.Width = titleMaxWidth
		title = d.focus.tempTitle.View()
	} else if lg.Width(title) > titleMaxWidth {
		// if past the max width, truncate until valid
		t := renderedTitle
		for lg.Width(t+"...") > titleMaxWidth && len(t) > 0 {
			t = t[:len(t)-1]
		}
		t = t + "..."
		title = titleStyle.Render(t)
	}

	left := typ + space + title
	right := priority

	px := styles.GetPaddingBetween(left, right, windowWidth, contStyle)
	content := left + styles.RenderPadding(segStyle, px) + right

	return contStyle.Width(windowWidth).Render(content)
}
