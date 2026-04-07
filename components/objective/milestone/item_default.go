package milestone

import (
	"log"
	"notion-project-tui/notion"

	"github.com/charmbracelet/bubbles/list"
)

type DefaultItem struct {
	ID              string
	Name            string
	MilestoneStatus notion.MilestoneStatus
	Progress        float64
	Icon            string
	TaskCount       int
	TaskGroups      notion.TaskGroups
	FetchStatus     FetchStatus
}

var _ list.Item = (*DefaultItem)(nil) // conform

func NewDefaultItem(page notion.MilestonePage) DefaultItem {
	title := notion.ExtractPlainText(page.Properties.Title.Title)

	var status notion.MilestoneStatus
	var err error
	if page.Properties.Status.Formula.String != nil {
		temp := *page.Properties.Status.Formula.String
		status, err = notion.MilestoneStatusFromString(temp)
		if err != nil {
			log.Printf(err.Error())
		}
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
		ID:              page.ID,
		Name:            title,
		MilestoneStatus: status,
		Progress:        progress,
		Icon:            icon,
		TaskCount:       cnt,
		TaskGroups:      notion.TaskGroups{},
		FetchStatus:     FetchIdle,
	}
}

func (x DefaultItem) FilterValue() string { return x.Name }
