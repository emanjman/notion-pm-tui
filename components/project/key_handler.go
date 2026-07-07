package project

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	log.Printf("project catches key msg") // !debug
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
}
