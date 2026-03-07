package milestonelist

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

type variantStyle struct {
	base     lg.Style
	selected lg.Style
}

type style struct {
	itemContainer variantStyle
	itemSegment   variantStyle
	header        variantStyle
}

type MilestoneListDelegate struct {
	focused bool
	style   style

	focus *FocusState
}

func NewMilestoneListDelegate(focused bool, focus *FocusState) MilestoneListDelegate {
	borderDistance := 1
	leftEdgeDistance := 1

	// item container style
	var (
		icbase = lg.NewStyle().
			Border(lg.NormalBorder(), false, false, true, false).
			BorderForeground(styles.BorderForeground).
			PaddingLeft(leftEdgeDistance + 2).
			PaddingRight(borderDistance)
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
			PaddingLeft(leftEdgeDistance).
			PaddingRight(borderDistance)
		hsel = hbase.
			Foreground(styles.PrimaryForeground).
			Underline(true)
	)

	return MilestoneListDelegate{
		focused: focused,
		style: style{
			itemContainer: variantStyle{base: icbase, selected: icsel},
			itemSegment:   variantStyle{base: isbase, selected: issel},
			header:        variantStyle{base: hbase, selected: hsel},
		},
		focus: focus,
	}
}

func (d MilestoneListDelegate) Height() int  { return 3 }
func (d MilestoneListDelegate) Spacing() int { return 0 }
func (d MilestoneListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

// render items (based on the list item type => header vs milestone)
func (d MilestoneListDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	selected := index == m.Index() && d.focused

	switch item := item.(type) {
	case listutil.ListItemGroupHeader:
		header := renderListItemGroupHeader(d, item, selected, m.Width())
		fmt.Fprint(w, header)
	case MilestoneListItem:
		milestone := renderMilestoneListItem(d, item, selected, m.Width())
		fmt.Fprint(w, milestone)
	}
}

// -- helper funcs

func progressBar(progress float64, width int) string {
	s := lg.NewStyle()

	wfilled := int(progress * float64(width))
	wempty := width - wfilled

	filled := strings.Repeat("█", wfilled)
	empty := strings.Repeat("░", wempty)

	return s.Foreground(styles.PrimaryForeground).Render(filled) +
		s.Foreground(styles.MutedForeground).Render(empty)
}

func renderListItemGroupHeader(d MilestoneListDelegate, item listutil.ListItemGroupHeader, selected bool, windowWidth int) string {
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

func renderMilestoneListItem(d MilestoneListDelegate, item MilestoneListItem, selected bool, windowWidth int) string {
	contStyle := d.style.itemContainer.base
	segStyle := d.style.itemSegment.base
	nameStyle, tagStyle := lg.Style{}, lg.Style{}

	// handle field highlighting by mode
	if selected {
		if d.focus.Mode == NeutralMode {
			segStyle = d.style.itemSegment.selected
			contStyle = d.style.itemContainer.selected

			// apply select highlight row-wide
			nameStyle, tagStyle = segStyle, segStyle
		} else {
			// apply select highlight by field
			switch d.focus.field {
			case MilestoneTitle:
				nameStyle = nameStyle.Inherit(d.style.itemSegment.selected)
			case MilestoneTag:
				tagStyle = tagStyle.Inherit(d.style.itemSegment.selected)
			}
		}
	}

	// apply final field-specific styles
	nameStyle = nameStyle.Foreground(styles.PrimaryForeground)
	tagStyle = tagStyle.Foreground(styles.MutedForeground)

	// render each field
	name := nameStyle.Render(item.Name)
	tag := tagStyle.Render(item.Tag)
	progress := segStyle.
		Foreground(styles.MutedForeground).
		Render(fmt.Sprintf("%.0f%%", item.Progress*100))
	bar := segStyle.
		Render(progressBar(item.Progress, windowWidth/3))
	activity := segStyle.
		Foreground(styles.MutedForeground).
		Render(item.LatestActivityLabel)

	// calculate max title width
	leftOffset, rightOffset := 3, 2
	offset := leftOffset + rightOffset
	nameMaxWidth := windowWidth - lg.Width(activity) - offset

	if selected && d.focus.Mode == WritingMode {
		// use textinput component in writing mode
		d.focus.tempTitle.Width = nameMaxWidth
		name = d.focus.tempTitle.View()
	} else if lg.Width(name) > nameMaxWidth {
		// if past the max width, truncate until valid
		n := item.Name
		for lg.Width(n+"...") > nameMaxWidth && len(n) > 0 {
			n = n[:len(n)-1]
		}
		n = n + "..."
		name = nameStyle.Render(n)
	}

	progressBar := progress + segStyle.Render(" ") + bar

	r1px := styles.GetPaddingBetween(name, activity, windowWidth, contStyle)
	r2px := styles.GetPaddingBetween(tag, progressBar, windowWidth, contStyle)
	r1 := name + styles.RenderPadding(segStyle, r1px) + activity
	r2 := tag + styles.RenderPadding(segStyle, r2px) + progressBar

	return contStyle.Width(windowWidth).Render(r1 + "\n" + r2)
}
