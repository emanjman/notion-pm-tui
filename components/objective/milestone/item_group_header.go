package milestone

import (
	"notion-project-tui/notion"

	"github.com/charmbracelet/bubbles/list"
)

type GroupHeaderItem struct {
	Label   string
	Hidden  bool
	Count   int
	HasMore bool
}

var _ list.Item = (*GroupHeaderItem)(nil) // conform

func NewGroupHeaderItem(label string, g notion.MilestoneGroup) GroupHeaderItem {
	return GroupHeaderItem{
		Label:   label,
		Hidden:  g.Hide,
		Count:   len(g.Milestones),
		HasMore: g.NextCursor != nil,
	}
}

func (h GroupHeaderItem) FilterValue() string { return "" }
