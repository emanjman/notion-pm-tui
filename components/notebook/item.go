package notebook

import (
	"notion-project-tui/notion"
	"time"
)

type ItemState int

const (
	Idle ItemState = iota
	Pending
	Success
	Failed
)

type Item struct {
	ID           string
	Title        string
	CreatedDate  time.Time
	CreatedLabel string

	Content string // defined on fetch
	State   ItemState
}

func (x Item) FilterValue() string { return x.Title }

func NewItem(page notion.NotePage) Item {
	label := ""
	if page.Properties.CreatedLabel.Formula.String != nil {
		label = *page.Properties.CreatedLabel.Formula.String
	}

	date, err := time.Parse(time.RFC3339Nano, page.Properties.CreatedDate.CreatedTime)
	if err != nil {
		date = time.Time{}
	}

	return Item{
		ID:           page.ID,
		Title:        notion.ExtractPlainText(page.Properties.Title.Title),
		CreatedDate:  date,
		CreatedLabel: label,
		Content:      "",
		State:        Idle,
	}
}
