package projectnotelist

import "github.com/charmbracelet/bubbles/list"

type ProjectNotesListModel struct {
	list    list.Model
	loading bool
	error   error
}

// func NewProjectNotesListModel() ProjectNotesListModel {
// }
