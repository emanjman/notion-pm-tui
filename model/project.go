package model

import (
	"fmt"
	"notion-project-tui/notion"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Tab int

// enum representation for better readability
const (
	MilestonesTab Tab = iota
	OverviewTab
	ProjectNotesTab
	DebugNotesTab
)

type ProjectModel struct {
	activeTab Tab
	// milestones   views.MilestonesListModel
	// overview     views.PageContentModel
	// projectNotes views.NotesListModel
	// debugNotes   views.NotesListModel
	page *notion.ProjectPage
	keys KeyMap

	help help.Model
}

func InitProjectModel() ProjectModel {
	return ProjectModel{
		activeTab: 0,
		page:      nil,
		keys:      DefaultKeyMap,
		help:      help.New(),
	}
}

func (m ProjectModel) Init() tea.Cmd {
	return notion.NewClient().FetchProjectById()
}

func (m ProjectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {

		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Up):
			return m, nil // todo: handle nav

		case key.Matches(msg, m.keys.Down):
			return m, nil // todo: handle nav

		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil

		}

	case notion.ProjectMsg:
		m.page = &msg.Data
		return m, nil

	}

	return m, nil
}

func (m ProjectModel) View() string {
	if m.page == nil {
		return "Loading..."
	}

	return fmt.Sprintf("Project ID: %s\n\n%s", m.page.ID, m.help.View(m.keys))
}
