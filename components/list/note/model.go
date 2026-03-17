package note

import (
	"notion-project-tui/notion"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	list    list.Model
	loading bool
	error   error
	notion  *notion.Client
}

func New(notion *notion.Client) Model {
	return Model{
		list:    list.New([]list.Item{}, NewItemDelegate(true), 0, 0),
		loading: false,
		error:   nil,
		notion:  notion,
	}
}

func (m Model) Init() tea.Cmd {
	// todo: kickoff fetch to relation ids
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	// todo: handle msg on fetched relation ids (request for relations)

	// todo: handle msg on fetched relations (build list)

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.list.View()
}
