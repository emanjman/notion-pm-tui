package explore

import (
	"notion-project-tui/components/project"
	"notion-project-tui/notion"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	notion *notion.Client

	list list.Model
	// pages   []notion.ProjectPage
	// pageIdx int

	project project.Model

	err     error
	loading bool
}

var _ tea.Model = (*Model)(nil) // conform

func New() Model {
	return Model{
		notion:  notion.NewClient(),
		loading: true,
	}
}

func (m Model) Init() tea.Cmd {
	return fetchProjects(m.notion)
}
