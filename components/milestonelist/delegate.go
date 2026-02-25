package milestonelist

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

// implementation for delegate
type MilestoneListDelegate struct {
	defaultStyle  lg.Style
	selectedStyle lg.Style
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
	case MilestoneListItemHeader:
		chevron := "▼"
		if item.Hidden {
			chevron = "▶"
		}

		row1 = fmt.Sprintf("%s %s", chevron, item.Label)
		row2 = fmt.Sprintf("%d", item.Count)

	case MilestoneListItem:
		var (
			name     = item.Name
			status   = item.Status
			tags     = strings.Join(item.Tags, " · ")
			progress = fmt.Sprintf("%.0f%%", item.Progress*100)
		)

		row1 = padBetween(name, status, m.Width())
		row2 = padBetween(tags, progress, m.Width())
	}

	block := row1 + "\n" + row2

	style := d.defaultStyle
	if selected {
		style = d.selectedStyle
	}

	// write to `w`
	fmt.Fprint(w, style.Width(m.Width()).Render(block))

}

func padBetween(left, right string, windowWidth int) string {
	staticPadding := 4 // left/right padding(2)

	// use lg.Width to only consider visible cells
	padding := windowWidth - lg.Width(left) - lg.Width(right) - staticPadding
	if padding < 0 {
		padding = 0
	}

	return left + strings.Repeat(" ", padding) + right
}

func NewMilestoneListDelegate() MilestoneListDelegate {
	base := lg.NewStyle().
		Border(lg.NormalBorder(), false, false, true, false).
		BorderForeground(lg.Color("238")).
		PaddingLeft(2).
		PaddingRight(2)

	return MilestoneListDelegate{
		defaultStyle: base,

		selectedStyle: base.
			Bold(true).
			Foreground(lg.Color("205")),
	}
}
