package milestonelist

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

// implementation for delegate
type MilestoneListDelegate struct {
	milestone delegateStyle
	header    delegateStyle
}

func NewMilestoneListDelegate() MilestoneListDelegate {
	milestoneBase := lg.NewStyle().
		Border(lg.NormalBorder(), false, false, true, false).
		BorderForeground(lg.Color("236")).
		PaddingLeft(4).
		PaddingRight(4)

	headerBase := lg.NewStyle().
		Border(lg.NormalBorder(), false, false, true, false).
		BorderForeground(lg.Color("236")).
		PaddingLeft(2).
		PaddingRight(2)

	return MilestoneListDelegate{
		milestone: delegateStyle{
			defaultStyle: milestoneBase,
			selectedStyle: milestoneBase.
				Foreground(lg.Color("205")),
		},

		header: delegateStyle{
			defaultStyle: headerBase,
			selectedStyle: headerBase.
				Foreground(lg.Color("205")),
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
	row1 := ""
	row2 := ""
	selected := index == m.Index()

	switch item := item.(type) {
	case listutil.ListItemGroupHeader:
		chevron := "▼"
		if item.Hidden {
			chevron = "▶"
		}

		style := d.header.defaultStyle
		if selected {
			style = d.header.selectedStyle
		}

		row1 = fmt.Sprintf("%s %s", chevron, item.Label)
		row2 = fmt.Sprintf("%d", item.Count)
		block := row1 + "\n" + row2

		fmt.Fprint(w, style.Width(m.Width()).Render(block))

	case MilestoneListItem:
		var (
			name     = item.Name
			status   = item.Status
			tags     = strings.Join(item.Tags, " · ")
			progress = fmt.Sprintf("%.0f%%", item.Progress*100)
		)

		style := d.milestone.defaultStyle
		if selected {
			style = d.milestone.selectedStyle
		}

		row1 = padBetween(name, status, m.Width(), style)
		row2 = padBetween(tags, progress, m.Width(), style)
		block := row1 + "\n" + row2

		// write to `w`
		fmt.Fprint(w, style.Width(m.Width()).Render(block))

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
