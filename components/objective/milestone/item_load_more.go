package milestone

import (
	"notion-project-tui/notion"

	"github.com/charmbracelet/bubbles/list"
)

type LoadMoreItem struct {
	Status  string
	Loading bool
}

var _ list.Item = (*LoadMoreItem)(nil) // conform

func NewLoadMoreItem(status string, g notion.MilestoneGroup) LoadMoreItem {
	return LoadMoreItem{
		Status:  status,
		Loading: g.Loading,
	}
}

func (_ LoadMoreItem) FilterValue() string { return "" }
