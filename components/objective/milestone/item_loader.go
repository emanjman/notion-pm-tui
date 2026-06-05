package milestone

import (
	"notion-project-tui/notion"

	"github.com/charmbracelet/bubbles/list"
)

type LoadMoreItem struct {
	Status  notion.MilestoneStatus
	Loading bool
}

var _ list.Item = (*LoadMoreItem)(nil) // conform

func NewLoadMoreItem(status notion.MilestoneStatus, g notion.MilestoneGroup) LoadMoreItem {
	return LoadMoreItem{
		Status:  status,
		Loading: g.Loading,
	}
}

func (_ LoadMoreItem) FilterValue() string { return "" }
