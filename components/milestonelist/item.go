package milestonelist

import (
	"notion-project-tui/notion"
)

// implementation for the `list.Item` interface
type MilestoneListItem struct {
	ID                  string
	TasksPropID         string
	Name                string
	Status              string
	LatestActivityLabel string
	Progress            float64
	Tag                 string
}

// func (m MilestoneListItem) Title() string       { return m.Name }
// func (m MilestoneListItem) Description() string { return m.Status }
func (m MilestoneListItem) FilterValue() string { return m.Name }
func (m MilestoneListItem) GroupKey() string    { return m.Status } // conform Groupable

func NewMilestoneListItem(page notion.MilestonePage) MilestoneListItem {
	title := notion.ExtractPlainText(page.Properties.Title.Title)

	status := ""
	if page.Properties.Status.Formula.String != nil {
		status = *page.Properties.Status.Formula.String
	}

	progress := 0.0
	if page.Properties.Progress.Formula.Number != nil {
		progress = *page.Properties.Progress.Formula.Number
	}

	tag := page.Properties.Tags.Select.Name

	label := ""
	if page.Properties.LatestActivityLabel.Formula.String != nil {
		label = *page.Properties.LatestActivityLabel.Formula.String
	}

	return MilestoneListItem{
		ID:                  page.ID,
		TasksPropID:         page.Properties.Tasks.ID,
		Name:                title,
		Status:              status,
		LatestActivityLabel: label,
		Progress:            progress,
		Tag:                 tag,
	}
}
