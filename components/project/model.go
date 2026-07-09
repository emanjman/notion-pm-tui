package project

import (
	"notion-project-tui/components/notebook"
	"notion-project-tui/components/objective"
	"notion-project-tui/notion"
	"time"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	projID string

	activeTab Tab

	page *notion.ProjectPage
	keys KeyMap

	width  int
	height int

	help     help.Model
	duration time.Duration

	notion *notion.Client

	objective objective.Model
	notebook  notebook.Model
}

// var _ tea.Model = (*Model)(nil) // conform

func New() Model {
	ntn := notion.NewClient()
	return Model{
		projID:    "",
		activeTab: ObjectiveTab,
		page:      nil,
		keys:      RootKeyMap,
		help:      help.New(),
		notion:    ntn,

		objective: objective.New(ntn),
		notebook:  notebook.New(ntn),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
