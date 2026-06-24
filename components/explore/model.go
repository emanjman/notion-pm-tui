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

	project       *project.Model
	neutralKeyMap NeutralKeyMap

	err     error
	loading bool
}

var _ tea.Model = (*Model)(nil) // conform

func New() Model {
	l := list.New([]list.Item{}, NewItemDelegate(), 0, 0)
	l.Title = "Explore Projects"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetShowFilter(false)
	l.SetShowPagination(false)
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()

	return Model{
		notion:        notion.NewClient(),
		list:          l,
		project:       nil,
		loading:       true,
		err:           nil,
		neutralKeyMap: NeutralKeyMapper,
	}
}

func (m Model) Init() tea.Cmd {
	return fetchProjects(m.notion)
}
