package milestone

import (
	"fmt"
	"io"
	"notion-project-tui/styles"
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

type ItemDelegate struct {
	focused bool
	style   style

	focus *FocusState
}

func NewItemDelegate(focused bool, focus *FocusState) ItemDelegate {
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
			Foreground(styles.MutedForeground).
			PaddingLeft(leftEdgeDistance).
			PaddingRight(borderDistance)
		hsel = lg.NewStyle().
			Foreground(styles.PrimaryForeground).
			Background(styles.SelectedBackground).
			PaddingLeft(leftEdgeDistance).
			PaddingRight(borderDistance)
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

func (d ItemDelegate) Height() int  { return 2 }
func (d ItemDelegate) Spacing() int { return 0 }
func (d ItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

// render items (based on the list item type => header vs milestone)
func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	selected := index == m.Index() && d.focused
	items := m.Items()
	isLast := index == len(items)-1
	nextIsHeader := !isLast && func() bool {
		_, ok := items[index+1].(GroupHeader)
		return ok
	}()
	noBorder := isLast || nextIsHeader

	switch item := item.(type) {
	case GroupHeader:
		header := renderItemHeader(d, item, selected, m.Width())
		fmt.Fprint(w, header)
	case Item:
		milestone := renderItem(d, item, selected, noBorder, m.Width())
		fmt.Fprint(w, milestone)
	case LoadMoreItem:
		fmt.Fprint(w, renderLoadMore(d, item.Loading, selected, noBorder, m.Width()))
	}
}

// -- helper funcs
func createProgressBar(progress float64, width int, baseStyle lg.Style) string {
	wfilled := int(progress * float64(width))
	wempty := width - wfilled

	filled := strings.Repeat("▬", wfilled)
	empty := strings.Repeat("▬", wempty)

	return baseStyle.Foreground(styles.TechForeground).Render(filled) +
		baseStyle.Foreground(lg.Color("#2a2a2a")).Render(empty)
}

func renderLoadMore(d ItemDelegate, loading bool, selected bool, noBorder bool, windowWidth int) string {
	style := d.style.itemContainer.base.Foreground(styles.MutedForeground).PaddingLeft(2)
	if selected {
		style = d.style.itemContainer.selected.Foreground(styles.MutedForeground).PaddingLeft(2)
	}
	if noBorder {
		style = style.Border(lg.NormalBorder(), false)
	}
	text := "..."
	if loading {
		text = "Loading..."
	} else if selected {
		text = "[Enter] to load more..."
	}
	rendered := style.Width(windowWidth).Render(text)
	if noBorder {
		return rendered + "\n" + lg.NewStyle().Render("")
	}
	return rendered
}

func renderItemHeader(d ItemDelegate, item GroupHeader, selected bool, windowWidth int) string {
	style := d.style.header.base
	if selected {
		style = d.style.header.selected
	}

	chevron := "▼"
	if item.Hidden {
		chevron = "▶"
	}

	count := fmt.Sprintf("%d", item.Count)
	if item.HasMore {
		count += "+"
	}
	content := fmt.Sprintf("%s %s (%s)", chevron, item.Label, count)
	label := style.Width(windowWidth).Render(content)
	spacer := lg.NewStyle().Render("")
	return label + "\n" + spacer
}

func renderItem(d ItemDelegate, item Item, selected bool, noBorder bool, windowWidth int) string {
	contStyle := d.style.itemContainer.base
	segStyle := d.style.itemSegment.base
	nameStyle, stateStyle, countStyle := lg.Style{}, lg.Style{}, lg.Style{}

	// handle field highlighting by mode
	if selected {
		if d.focus.Mode == NeutralMode {
			segStyle = d.style.itemSegment.selected
			countStyle = d.style.itemContainer.selected

			// apply select highlight row-wide
			nameStyle, stateStyle, countStyle = segStyle, segStyle, segStyle
		} else {
			nameStyle = nameStyle.Inherit(d.style.itemSegment.selected)
			countStyle = countStyle.Inherit(d.style.itemSegment.selected)
		}
	}

	// apply final field-specific styles
	nameStyle = nameStyle.Foreground(styles.PrimaryForeground)
	countStyle = countStyle.Foreground(styles.MutedForeground)

	if item.Icon == "" {
		item.Icon = "  "
	}
	name := nameStyle.Render(item.Icon + " " + item.Name)
	count := countStyle.Render(fmt.Sprintf("%d", item.TaskCount))

	var state string
	switch item.FetchState {
	case Idle:
		state = stateStyle.Foreground(styles.MutedForeground).Render("◌")
	case Pending:
		state = stateStyle.Foreground(styles.MutedForeground).Render("↻")
	case Failed:
		state = stateStyle.Foreground(lg.Color("#e0af68")).Render("⚠")
	}
	space := segStyle.Render(" ")

	// hide progress bar for completed milestones
	var progress string
	if item.Status != "🎉 complete" { // !hardcode
		// completion := segStyle.
		// 	Foreground(styles.MutedForeground).
		// 	Render(fmt.Sprintf("%.0f%%", item.Progress*100))
		pbar := createProgressBar(item.Progress, windowWidth/4, segStyle)
		progress = pbar + segStyle.Render(" ") + count
	}

	// calculate max title width
	leftOffset, rightOffset := 3, 2
	offset := leftOffset + rightOffset
	nameMaxWidth := windowWidth - lg.Width(progress+space+state) - offset

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

	if noBorder {
		contStyle = contStyle.Border(lg.NormalBorder(), false)
	}

	left := name
	right := progress + space + state
	px := styles.GetPaddingBetween(left, right, windowWidth, contStyle)
	content := left + styles.RenderPadding(segStyle, px) + right

	rendered := contStyle.Width(windowWidth).Render(content)
	if noBorder {
		return rendered + "\n" + lg.NewStyle().Render("")
	}
	return rendered
}
