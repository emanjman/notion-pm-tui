package projectnotelist

import (
	"notion-project-tui/notion"
	"time"
)

type NoteListItem struct {
	ID          string
	NoteTitle   string
	CreatedTime time.Time
}

func (i NoteListItem) FilterValue() string { return i.NoteTitle }
func (i NoteListItem) Title() string       { return i.NoteTitle }
func (i NoteListItem) Description() string { return i.CreatedTime.Format("Jan 2, 2006") }

func NewNoteListItem(page notion.ProjectNotePage) NoteListItem {
	title := notion.ExtractPlainText(page.Properties.Title.Title)

	createdTime, _ := time.Parse(time.RFC3339, page.CreatedTime)

	return NoteListItem{
		ID:          page.ID,
		NoteTitle:   title,
		CreatedTime: createdTime,
	}
}
