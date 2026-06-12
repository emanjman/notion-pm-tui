package app

import (
	// ! temp, styling ui
	// "fmt"
	"notion-project-tui/components/notebook"
	"notion-project-tui/components/objective"
	"notion-project-tui/notion"
	"notion-project-tui/styles"
	"notion-project-tui/util/keymap"
	"os"
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
	NotebookTab
	BugsTab
	TechTab
)
const tabCount = 5

var labels = []string{
	"Objective (n%)",
	"Overview",
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
		objective: objective.New(c, os.Getenv("NOTION_HOOP_ARCHIVES_ID"), "P%60%60s"),
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

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if m.objective.ChildPriorityMode() {
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
				// case OverviewTab:
				// 	m.overview, cmd = m.overview.Update(msg)
				case NotebookTab:
					m.notebook, cmd = m.notebook.Update(msg)
				}

				return m, cmd
			}
		}

	// send window size to the milestones model
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		childHeight := msg.Height - 6 // header/footer row, help bar, spacing

		var objCmd, ovrCmd, noteCmd tea.Cmd
		m.objective, objCmd = m.objective.Update(tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: childHeight,
		})

		// m.overview, ovrCmd = m.overview.Update(tea.WindowSizeMsg{
		// 	Width:  msg.Width,
		// 	Height: childHeight,
		// })

		m.notebook, noteCmd = m.notebook.Update(tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: childHeight,
		})

		return m, tea.Batch(objCmd, ovrCmd, noteCmd)

	// todo: on msg, start forwarding the commands to children??? use default?
	case notion.ProjectMsg:
		// if failed fetch, don't proceed w/ fetching ids
		if msg.Err != nil {
			return m, nil
		}

		// update project data
		m.page = &msg.Data
		m.duration = msg.Duration

		return m, nil

	// spill messages into children to be handled at that lvl
	default:
		var objCmd, ovrCmd, noteCmd tea.Cmd
		m.objective, objCmd = m.objective.Update(msg)
		// m.overview, ovrCmd = m.overview.Update(msg)
		m.notebook, noteCmd = m.notebook.Update(msg)
		return m, tea.Batch(objCmd, ovrCmd, noteCmd)
	}
}

func (m Model) View() string {
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
		// main = m.overview.View()
		main = "Overview (coming soon)"
	case NotebookTab:
		main = m.notebook.View()
	case BugsTab:
		main = "Debug notes (coming soon)"
	case TechTab:
		main = "Tech notes (coming soon)"
	}

	tabDivider := lg.NewStyle().
		Foreground(styles.BorderForeground).
		SetString("|")
	s.WriteString(
		lg.NewStyle().
			Padding(1, 2).
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
			Padding(1, 2).
			Width(m.width).
			Render(help))

	return s.String()
}

func (m Model) getActiveKeyMap() help.KeyMap {
	switch m.activeTab {

	case ObjectiveTab:
		return m.objective.KeyMap()

	// todo: handle other tabs
	case OverviewTab:
		return m.objective.KeyMap() // todo: change
	case NotebookTab:
		return m.notebook.ActiveKeyMap
	case BugsTab:
		return m.objective.KeyMap()
	case TechTab:
		return m.objective.KeyMap()

	default:
		// return nil
		return m.objective.KeyMap()

	}

}
