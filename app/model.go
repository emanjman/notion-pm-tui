package app

import (
	// ! temp, styling ui
	// "fmt"
	"notion-project-tui/components/objective"
	"notion-project-tui/notion"
	"notion-project-tui/styles"
	"notion-project-tui/util/keymap"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type Tab int

// enum representation for better readability
const (
	ObjectiveTab = iota
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

	objective objective.ObjectiveModel
	// overview     views.PageContentModel
	// projectNotes views.NotesListModel
	// debugNotes   views.NotesListModel
}

func InitProjectModel() ProjectModel {
	client := notion.NewClient()
	return ProjectModel{
		activeTab: 0,
		page:      nil,
		keys:      RootKeyMap,
		help:      help.New(),
		client:    client,
		objective: objective.NewObjectiveModel(client),
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

			case ObjectiveTab:
				m.objective, cmd = m.objective.Update(msg)

			}

			return m, cmd
		}

	// send window size to the milestones model
	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		m.objective, cmd = m.objective.Update(msg)
		return m, cmd

	case notion.ProjectMsg:
		// if failed fetch, don't proceed w/ fetching ids
		if msg.Err != nil {
			return m, nil
		}

		// update project data
		m.page = &msg.Data
		m.duration = msg.Duration

		return m, m.client.FetchMilestoneRelationIds(m.page.ID, m.page.Properties.Milestones.ID)

	case notion.MilestoneRelationIdsMsg:
		// if failed fetch, don't proceed w/ milestones fetch
		if msg.Err != nil {
			return m, nil
		}

		return m, m.client.FetchMilestones(msg.IDs)

	// forward updated milestones model + cmd
	case notion.MilestoneMsg:
		var cmd tea.Cmd
		m.objective, cmd = m.objective.Update(msg)
		return m, cmd

	default:
		var cmd tea.Cmd
		m.objective, cmd = m.objective.Update(msg)
		return m, cmd
	}
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

	case ObjectiveTab:
		view.WriteString(m.objective.View())
	case OverviewTab:
		view.WriteString("Overview (coming soon)")
	case ProjectNotesTab:
		view.WriteString("Project notes (coming soon)")
	case DebugNotesTab:
		view.WriteString("Debug notes (coming soon)")

	}

	view.WriteString("\n")

	help := m.help.View(keymap.JoinedKeyMap{
		Primary:   RootKeyMap,
		Secondary: m.getActiveKeyMap(),
	})

	view.WriteString(
		lg.NewStyle().
			BorderTop(true).
			BorderStyle(lg.NormalBorder()).
			BorderForeground(styles.BorderForeground).
			Render(help))

	return view.String()
}

func (m ProjectModel) getActiveKeyMap() help.KeyMap {
	switch m.activeTab {

	case ObjectiveTab:
		return m.objective.KeyMap()

	// todo... handle other tabs

	default:
		return nil

	}

}
