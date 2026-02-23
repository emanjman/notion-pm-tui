package views

import "notion-project-tui/notion"
import "notion-project-tui/model"

func NewMilestoneListItem(page notion.MilestonePage) model.MilestoneListItem {
	title := notion.ExtractPlainText(page.Properties.Title.Title)

	status := ""
	if page.Properties.Status.Formula.String != nil {
		status = *page.Properties.Status.Formula.String
	}

	progress := 0.0
	if page.Properties.Progress.Formula.Number != nil {
		progress = *page.Properties.Progress.Formula.Number
	}

	tags := make([]string, len(page.Properties.Tags.MultiSelect))
	for i, tag := range page.Properties.Tags.MultiSelect {
		tags[i] = tag.Name
	}

	return model.MilestoneListItem{
		ID:       page.ID,
		Name:     title,
		Status:   status,
		Progress: progress,
		Tags:     tags,
	}
}
