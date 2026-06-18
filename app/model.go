package app

import (
	// ! temp, styling ui
	// "fmt"
	"notion-project-tui/components/notebook"
	"notion-project-tui/components/objective"
	"notion-project-tui/notion"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

type Tab int

// enum representation for better readability
const (
	ObjectiveTab Tab = iota
	// OverviewTab
	NotebookTab
	BugsTab
	TechTab
)
const tabCount = 4 // todo: can inject labels/cnt into model itself?

var labels = []string{
	"Objective (n%)",
	// "Overview",
	"Notebook (n)",
	"Bug Report (n)",
	"Technology (n)",
}

type Model struct {
	activeTab Tab

	page *notion.ProjectPage
	keys KeyMap

	width  int
	height int

	help     help.Model
	duration time.Duration

	notion *notion.Client

	objective objective.Model
	// overview  overview.Model
	notebook notebook.Model
	// debugNotes   views.NotesListModel
}

var _ tea.Model = (*Model)(nil) // conform

func New() Model {
	c := notion.NewClient()
	return Model{
		activeTab: ObjectiveTab,
		page:      nil,
		keys:      RootKeyMap,
		help:      help.New(),
		notion:    notion.NewClient(),
		objective: objective.New(c, os.Getenv("NOTION_HOOP_ARCHIVES_ID")),
		// overview:  overview.New(c),

		notebook: notebook.New(c, os.Getenv("NOTION_HOOP_ARCHIVES_ID"), "%7BGKi"),
	}
}

func (m Model) Init() tea.Cmd {
	// return m.client.FetchProject()
	// return nil // ! temp, styling ui

	return tea.Batch(
		// m.client.FetchProject(),
		m.objective.Init(),
		// m.overview.Init(),
		m.notebook.Init(),
	)
}
