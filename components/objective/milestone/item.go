package milestone

import (
	"notion-project-tui/notion"
)

type FetchState int

const (
	Idle FetchState = iota
	Pending
	Success
	Failed
)

// implementation for the `list.Item` interface
type Item struct {
	ID       string
	Name     string
	Status   string
	Progress float64
	Icon     string

	// source of truth for this milestone's task data in memory. persists across
	// milestone switches so previously fetched tasks don't need to be refetched.
	TaskGroups notion.TaskGroups
	FetchState FetchState
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

	icon := ""
	if page.Icon != nil && page.Icon.Emoji != nil {
		icon = *page.Icon.Emoji
	}

	return Item{
		ID:       page.ID,
		Name:     title,
		Status:   status,
		Progress: progress,
		Icon:     icon,

		TaskGroups: notion.TaskGroups{},
		FetchState: Idle,
	}
}
