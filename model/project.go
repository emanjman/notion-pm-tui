package model

import (
	"fmt"
	"notion-project-tui/notion"
	"strings"
	"time"

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

	page *notion.ProjectPage
	keys KeyMap

	help     help.Model
	duration time.Duration

	client *notion.Client

	milestones MilestonesModel
	// overview     views.PageContentModel
	// projectNotes views.NotesListModel
	// debugNotes   views.NotesListModel
}

func InitProjectModel() ProjectModel {
	return ProjectModel{
		activeTab:  0,
		page:       nil,
		keys:       DefaultKeyMap,
		help:       help.New(),
		client:     notion.NewClient(),
		milestones: NewMilestonesModel(),
	}
}

func (m ProjectModel) Init() tea.Cmd {
	return m.client.FetchProject()
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
		// if failed fetch, don't proceed w/ milestones fetch
		if msg.Err != nil {
			return m, nil
		}

		m.page = &msg.Data
		m.duration = msg.Duration

		milestoneIds, err := m.client.FetchAllRelationIds(m.page.ID, m.page.Properties.Milestones)
		if err != nil {
			return m, nil
		}

		// get all milestones after getting all ids
		return m, m.client.FetchMilestones(milestoneIds)

	case notion.MilestoneMsg:
		var cmd tea.Cmd
		// forward updated milestones model + cmd
		m.milestones, cmd = m.milestones.Update(msg)
		return m, cmd

	}

	return m, nil
}

func (m ProjectModel) View() string {
	if m.page == nil {
		return "Loading..."
	}

	var view strings.Builder

	view.WriteString(fmt.Sprintf("Project ID: %s", m.page.ID))
	view.WriteString("\n\n")
	view.WriteString(fmt.Sprintf("Fetched in %dms", m.duration.Milliseconds()))
	view.WriteString("\n\n")
	view.WriteString(m.help.View(m.keys))

	return view.String()
}
