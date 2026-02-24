package model

import (
	// ! temp, styling ui
	// "fmt"
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

	milestones MilestoneListModel
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
		milestones: NewMilestoneListModel(),
	}
}

func (m ProjectModel) Init() tea.Cmd {
	// return m.client.FetchProject()
	return nil // ! temp, styling ui
}

func (m ProjectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {

		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil

		// send other keymaps to the active tab
		default:
			var cmd tea.Cmd

			switch m.activeTab {

			case MilestonesTab:
				m.milestones, cmd = m.milestones.Update(msg)

			}

			return m, cmd
		}

	// send window size to the milestones model
	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		m.milestones, cmd = m.milestones.Update(msg)
		return m, cmd

	case notion.ProjectMsg:
		// if failed fetch, don't proceed w/ fetching ids
		if msg.Err != nil {
			return m, nil
		}

		// update project data
		m.page = &msg.Data
		m.duration = msg.Duration

		return m, m.client.FetchAllRelationIds(m.page.ID, m.page.Properties.Milestones)

	case notion.RelationIdsMsg:
		// if failed fetch, don't proceed w/ milestones fetch
		if msg.Err != nil {
			return m, nil
		}

		return m, m.client.FetchMilestones(msg.IDs)

	// forward updated milestones model + cmd
	case notion.MilestoneMsg:
		var cmd tea.Cmd
		m.milestones, cmd = m.milestones.Update(msg)
		return m, cmd

	}

	return m, nil
}

func (m ProjectModel) View() string {
	// ! temp, styling ui
	// if m.page == nil {
	// 	return "Loading project..."
	// }

	var view strings.Builder

	// ! temp, styling ui
	// view.WriteString(fmt.Sprintf("Project ID: %s", m.page.ID))
	// view.WriteString("\n\n")
	// view.WriteString(fmt.Sprintf("Fetched in %dms", m.duration.Milliseconds()))
	// view.WriteString("\n\n")

	switch m.activeTab {

	case MilestonesTab:
		view.WriteString(m.milestones.View())
	case OverviewTab:
		view.WriteString("Overview (coming soon)")
	case ProjectNotesTab:
		view.WriteString("Project notes (coming soon)")
	case DebugNotesTab:
		view.WriteString("Debug notes (coming soon)")

	}

	view.WriteString("\n\n")
	view.WriteString(m.help.View(m.keys))

	return view.String()
}
