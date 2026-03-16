package milestone

import (
	"notion-project-tui/notion"
)

// implementation for the `list.Item` interface
type Item struct {
	ID                  string
	TasksPropID         string
	Name                string
	Status              string
	LatestActivityLabel string
	Progress            float64
	Tag                 string
}

func (m Item) FilterValue() string { return m.Name }
func (m Item) GroupKey() string    { return m.Status } // conform Groupable

func NewItem(page notion.MilestonePage) Item {
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

	return Item{
		ID:                  page.ID,
		TasksPropID:         page.Properties.Tasks.ID,
		Name:                title,
		Status:              status,
		LatestActivityLabel: label,
		Progress:            progress,
		Tag:                 tag,
	}
}
