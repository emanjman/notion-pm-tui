package project

import (
	// ! temp, styling ui
	// "fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"notion-project-tui/notion"
)

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
