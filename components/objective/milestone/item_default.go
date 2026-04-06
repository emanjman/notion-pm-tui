package milestone

import (
	"notion-project-tui/notion"

	"github.com/charmbracelet/bubbles/list"
)

type DefaultItem struct {
	ID         string
	Name       string
	Status     string
	Progress   float64
	Icon       string
	TaskCount  int
	TaskGroups notion.TaskGroups
	FetchState FetchState
}

var _ list.Item = (*DefaultItem)(nil) // conform

func NewDefaultItem(page notion.MilestonePage) DefaultItem {
	title := notion.ExtractPlainText(page.Properties.Title.Title)
	status := ""
	if page.Properties.Status.Formula.String != nil {
		status = *page.Properties.Status.Formula.String
	}
	progress := 0.0
	if page.Properties.Progress.Formula.Number != nil {
		progress = *page.Properties.Progress.Formula.Number
	}
	icon := ""
	if page.Icon != nil && page.Icon.Emoji != nil {
		icon = *page.Icon.Emoji
	}
	cnt := 0
	if page.Properties.TaskCount.Formula.Number != nil {
		cnt = int(*page.Properties.TaskCount.Formula.Number)
	}

	return DefaultItem{
		ID:         page.ID,
		Name:       title,
		Status:     status,
		Progress:   progress,
		Icon:       icon,
		TaskCount:  cnt,
		TaskGroups: notion.TaskGroups{},
		FetchState: Idle,
	}
}

func (m DefaultItem) FilterValue() string { return m.Name }
