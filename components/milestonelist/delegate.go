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
}

func NewMilestoneListDelegate(focused bool) MilestoneListDelegate {
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

	return MilestoneListDelegate{
		focused: focused,
		style: style{
			itemContainer: variantStyle{base: icbase, selected: icsel},
			itemSegment:   variantStyle{base: isbase, selected: issel},
			header:        variantStyle{base: hbase, selected: hsel},
		},
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

	case MilestoneListItem:
		segStyle := d.style.itemSegment.base
		contStyle := d.style.itemContainer.base
		if selected {
			segStyle = d.style.itemSegment.selected
			contStyle = d.style.itemContainer.selected
		}

		var (
			name     = segStyle.Foreground(styles.PrimaryForeground).Render(item.Name)
			tags     = segStyle.Foreground(styles.MutedForeground).Render(strings.Join(item.Tags, " · "))
			progress = segStyle.Render(fmt.Sprintf("%.0f%%", item.Progress*100))
			bar      = segStyle.Render(progressBar(item.Progress, int(m.Width())/3))
		)

		style := d.style.itemContainer.base
		if selected {
			style = d.style.itemContainer.selected
		}

		r1px := styles.GetPaddingBetween(name, progress, m.Width(), contStyle)
		r2px := styles.GetPaddingBetween(tags, bar, m.Width(), contStyle)
		r1 := name + styles.RenderPadding(segStyle, r1px) + progress
		r2 := tags + styles.RenderPadding(segStyle, r2px) + bar

		fmt.Fprint(w, style.Width(m.Width()).Render(r1+"\n"+r2))
	}
}

func progressBar(progress float64, width int) string {
	s := lg.NewStyle()

	wfilled := int(progress * float64(width))
	wempty := width - wfilled

	filled := strings.Repeat("█", wfilled)
	empty := strings.Repeat("░", wempty)

	return s.Foreground(styles.PrimaryForeground).Render(filled) +
		s.Foreground(styles.MutedForeground).Render(empty)
}
