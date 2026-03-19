package notebook

import (
	"fmt"
	"io"
	"notion-project-tui/styles"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type variantStyle struct {
	base    lg.Style
	focused lg.Style
}

type style struct {
	itemContainer variantStyle
	itemSegment   variantStyle
}

type ItemDelegate struct {
	sectionFocused bool
	style          style
}

func NewItemDelegate(sectionFocused bool) ItemDelegate {
	borderDistance := 0
	rightEdgeDistance := 3

	var (
		contBase = lg.NewStyle().
				Border(lg.NormalBorder(), false, false, true, false).
				BorderForeground(styles.BorderForeground).
				PaddingLeft(borderDistance + 2).
				PaddingRight(rightEdgeDistance)
		contFocus = contBase.
				Background(styles.SelectedBackground)
	)

	var (
		segBase  = lg.NewStyle()
		segFocus = segBase.
				Background(styles.SelectedBackground)
	)

	return ItemDelegate{
		sectionFocused: sectionFocused,
		style: style{
			itemContainer: variantStyle{base: contBase, focused: contFocus},
			itemSegment:   variantStyle{base: segBase, focused: segFocus},
		},
	}
}

func (d ItemDelegate) Height() int                               { return 2 }
func (d ItemDelegate) Spacing() int                              { return 0 }
func (d ItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

var statusColors = []lg.Color{
	styles.MutedForeground, // idle
	lg.Color("#24a7ff"),    //  pending
	lg.Color("#24ff7b"),    // success
	lg.Color("#ff244c"),    // failed
}

func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	focused := index == m.Index() && d.sectionFocused

	switch item := item.(type) {
	case Item:
		contStyle := d.style.itemContainer.base
		segStyle := d.style.itemSegment.base
		titleStyle, dateStyle, stateStyle := lg.Style{}, lg.Style{}, lg.Style{}

		if focused {
			segStyle = d.style.itemSegment.focused
			contStyle = d.style.itemContainer.focused
			titleStyle, dateStyle, stateStyle = segStyle, segStyle, segStyle
		}

		// apply field-specific styles
		titleStyle = titleStyle.Foreground(styles.PrimaryForeground)
		dateStyle = dateStyle.Foreground(styles.MutedForeground)
		stateStyle = stateStyle.Foreground(statusColors[item.State])

		// render each field
		title := titleStyle.Render(item.Title)

		created := dateStyle.Render(item.CreatedLabel)
		state := stateStyle.Render("●")
		space := segStyle.Render(" ")

		left := title
		right := created + space + state

		px := styles.GetPaddingBetween(left, right, m.Width(), contStyle)
		content := left + styles.RenderPadding(segStyle, px) + right

		fmt.Fprint(w, contStyle.Width(m.Width()).Render(content))
	}
}
