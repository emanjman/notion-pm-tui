package explore

import (
	"notion-project-tui/notion"

	"github.com/charmbracelet/bubbles/list"
)

// list friendly version of the notion-page
type DefaultItem struct {
	ID    string
	Title string
	Icon  string
}

var _ list.Item = (*DefaultItem)(nil) // conform

func NewDefaultItem(pg notion.ProjectPage) DefaultItem {
	title := notion.ExtractPlainText(pg.Properties.Title.Title)

	icon := ""
	if pg.Icon != nil && pg.Icon.Emoji != nil {
		icon = *pg.Icon.Emoji
	}

	return DefaultItem{
		ID:    pg.ID,
		Title: title,
		Icon:  icon,
	}
}

func (x DefaultItem) FilterValue() string { return x.Title }
