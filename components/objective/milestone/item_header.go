package milestone

import (
	"notion-project-tui/notion"

	"github.com/charmbracelet/bubbles/list"
)

type GroupHeaderItem struct {
	Status  notion.MilestoneStatus
	Hidden  bool
	Count   int
	HasMore bool
}

var _ list.Item = (*GroupHeaderItem)(nil) // conform

func NewGroupHeaderItem(status notion.MilestoneStatus, g notion.MilestoneGroup) GroupHeaderItem {
	return GroupHeaderItem{
		Status:  status,
		Hidden:  g.Hide,
		Count:   len(g.Milestones),
		HasMore: g.NextCursor != nil,
	}
}

func (x GroupHeaderItem) FilterValue() string { return x.Status.String() }
