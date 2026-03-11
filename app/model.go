package app

import (
	// ! temp, styling ui
	// "fmt"
	"notion-project-tui/components/objective"
	"notion-project-tui/components/overview"
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
	ObjectiveTab Tab = iota
	OverviewTab
	ProjectNotesTab
	DebugNotesTab
)
const tabCount = 4

var labels = []string{"Objective (n%)", "Overview", "Project Notes (n)", "Debug Notes (n)"}

type ProjectModel struct {
	activeTab Tab

	page *notion.ProjectPage
	keys KeyMap

	width  int
	height int

	help     help.Model
	duration time.Duration

	client *notion.Client

	objective objective.ObjectiveModel
	overview  overview.OverviewModel
	// projectNotes views.NotesListModel
	// debugNotes   views.NotesListModel
}

func InitProjectModel() ProjectModel {
	client := notion.NewClient()
	return ProjectModel{
		activeTab: ObjectiveTab,
		page:      nil,
		keys:      RootKeyMap,
		help:      help.New(),
		client:    client,
		objective: objective.NewObjectiveModel(client),
		overview:  overview.NewOverviewModel(client),
	}
}

func (m ProjectModel) Init() tea.Cmd {
	// return m.client.FetchProject()
	// return nil // ! temp, styling ui

	return tea.Batch(
		// m.client.FetchProject(),
		m.overview.Init(),
	)
}

func (m ProjectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if m.objective.InFocusMode() {
			// forward all keys if in writing mode
			var cmd tea.Cmd
			m.objective, cmd = m.objective.Update(msg)
			return m, cmd
		} else {
			switch {

			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit
			case key.Matches(msg, m.keys.Help):
				m.help.ShowAll = !m.help.ShowAll
				return m, nil
			case key.Matches(msg, m.keys.Next):
				m.activeTab = (m.activeTab + 1) % tabCount
				return m, nil
			case key.Matches(msg, m.keys.Prev):
				if m.activeTab == 0 {
					m.activeTab = tabCount - 1
				} else {
					m.activeTab = (m.activeTab - 1) % tabCount
				}
				return m, nil

			// send other keymaps to the active tab
			default:
				var cmd tea.Cmd

				// todo: handle for other tabs
				switch m.activeTab {
				case ObjectiveTab:
					m.objective, cmd = m.objective.Update(msg)
				case OverviewTab:
					m.overview, cmd = m.overview.Update(msg)
				}

				return m, cmd
			}
		}

	// send window size to the milestones model
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		childHeight := msg.Height - 5 // header row, help bar, spacing

		var objCmd, ovrCmd tea.Cmd
		m.objective, objCmd = m.objective.Update(tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: childHeight,
		})

		m.overview, ovrCmd = m.overview.Update(tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: childHeight,
		})

		return m, tea.Batch(objCmd, ovrCmd)

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
		var objCmd, ovrCmd tea.Cmd
		m.objective, objCmd = m.objective.Update(msg)
		m.overview, ovrCmd = m.overview.Update(msg)
		return m, tea.Batch(objCmd, ovrCmd)
	}
}

func (m ProjectModel) View() string {
	// ! temp, styling ui
	// if m.page == nil {
	// 	return "Loading project..."
	// }

	var s strings.Builder

	// ! temp, styling ui
	// view.WriteString(fmt.Sprintf("Project ID: %s", m.page.ID))
	// view.WriteString("\n\n")
	// view.WriteString(fmt.Sprintf("Fetched in %dms", m.duration.Milliseconds()))
	// view.WriteString("\n\n")

	lg.JoinHorizontal(lg.Top)
	headers := make([]string, len(labels))

	for i := range labels {
		base := lg.NewStyle().Padding(0, 2)

		tabStyle := base.Foreground(styles.MutedForeground)
		if int(m.activeTab) == i {
			tabStyle = base.
				Foreground(styles.PrimaryForeground).
				Background(styles.SelectedBackground)
		}
		headers[i] = tabStyle.Render(labels[i])
	}
	main := ""

	switch m.activeTab {

	case ObjectiveTab:
		main = m.objective.View()
	case OverviewTab:
		main = m.overview.View()
	case ProjectNotesTab:
		main = "Project notes (coming soon)"
	case DebugNotesTab:
		main = "Debug notes (coming soon)"

	}

	tabDivider := lg.NewStyle().
		Foreground(styles.BorderForeground).
		SetString("|")

	s.WriteString(
		lg.NewStyle().
			PaddingLeft(2).
			PaddingTop(1).
			BorderBottom(true).
			BorderStyle(lg.ThickBorder()).
			BorderForeground(styles.BorderForeground).
			Width(m.width).
			Render(strings.Join(headers, tabDivider.String())))

	s.WriteString("\n")
	s.WriteString(main)
	s.WriteString("\n")

	help := m.help.View(keymap.JoinedKeyMap{
		Primary:   RootKeyMap,
		Secondary: m.getActiveKeyMap(),
	})

	s.WriteString(
		lg.NewStyle().
			BorderTop(true).
			BorderStyle(lg.NormalBorder()).
			BorderForeground(styles.BorderForeground).
			Width(m.width).
			Render(help))

	return s.String()
}

func (m ProjectModel) getActiveKeyMap() help.KeyMap {
	switch m.activeTab {

	case ObjectiveTab:
		return m.objective.KeyMap()

	// todo: handle other tabs
	case OverviewTab:
		return m.objective.KeyMap()
	case ProjectNotesTab:
		return m.objective.KeyMap()
	case DebugNotesTab:
		return m.objective.KeyMap()

	default:
		// return nil
		return m.objective.KeyMap()

	}

}
