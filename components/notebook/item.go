package notebook

import (
	"notion-project-tui/notion"
	"time"
)

type ContentState int

const (
	Idle ContentState = iota
	Pending
	Success
	Failed
)

type Item struct {
	ID           string
	Title        string
	CreatedDate  time.Time
	CreatedLabel string
	Icon         string

	Content       string // defined on fetch
	Markdown      string // defined on fetch
	BlocksState   ContentState
	MarkdownState ContentState
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

	icon := ""
	if page.Icon != nil && page.Icon.Emoji != nil {
		icon = *page.Icon.Emoji
	}

	return Item{
		ID:           page.ID,
		Title:        notion.ExtractPlainText(page.Properties.Title.Title),
		CreatedDate:  date,
		CreatedLabel: label,
		Icon:         icon,

		Content:       "",
		Markdown:      "",
		BlocksState:   Idle,
		MarkdownState: Idle,
	}
}
