package notebook

import "notion-project-tui/notion"

type ItemState string

const (
	Idle    ItemState = "⚫"
	Pending ItemState = "🔵"
	Success ItemState = "🟢"
	Failed  ItemState = "🔴"
)

type Item struct {
	ID           string
	Title        string
	CreatedLabel string

	Content string // defined on fetch
	State   ItemState
}

func (x Item) FilterValue() string { return x.Title }

func NewItem(page notion.NotePage) Item {
	createdLabel := ""
	if page.Properties.CreatedLabel.Formula.String != nil {
		createdLabel = *page.Properties.CreatedLabel.Formula.String
	}

	return Item{
		ID:           page.ID,
		Title:        notion.ExtractPlainText(page.Properties.Title.Title),
		CreatedLabel: createdLabel,
		Content:      "",
		State:        Idle,
	}
}
