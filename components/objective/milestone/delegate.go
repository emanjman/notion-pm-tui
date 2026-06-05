package milestone

import (
	"fmt"
	"io"
	"notion-project-tui/styles"

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

	mode *Mode
	edit *EditModeCtx
}

var _ list.ItemDelegate = (*ItemDelegate)(nil) // compile-time compliance

func NewItemDelegate(focused bool, mode *Mode, edit *EditModeCtx) ItemDelegate {
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
		mode: mode,
		edit: edit,
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
		_, ok := items[index+1].(GroupHeaderItem)
		return ok
	}()
	noBorder := isLast || nextIsHeader

	switch item := item.(type) {
	case GroupHeaderItem:
		header := renderItemHeader(d, item, selected, m.Width())
		fmt.Fprint(w, header)
	case DefaultItem:
		mstone := renderItem(d, item, selected, noBorder, m.Width())
		fmt.Fprint(w, mstone)
	case LoadMoreItem:
		fmt.Fprint(w, renderLoadMore(d, item.Loading, selected, noBorder, m.Width()))
	}
}
