package projectnotes

import (
	"notion-project-tui/components/pagecontent"
	"notion-project-tui/components/projectnotelist"
)

type Panel int

const (
	NotesPanel Panel = iota
	ContentPanel
)

type ProjectNotesModel struct {
	focus   Panel
	notes   projectnotelist.ProjectNotesListModel
	content pagecontent.PageContentModel
	// keys       KeyMap
}
