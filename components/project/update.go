package project

import (
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	case tea.WindowSizeMsg:
		return m.handleWindow(msg)
	default:
		return m.handleDefault(msg)
	}
}

// spill messages into children to be handled at that lvl
func (m Model) handleDefault(msg tea.Msg) (Model, tea.Cmd) {
	log.Printf("project update catches default") // !debug
	var objCmd, ovrCmd, noteCmd tea.Cmd
	m.objective, objCmd = m.objective.Update(msg)
	// m.overview, ovrCmd = m.overview.Update(msg)
	m.notebook, noteCmd = m.notebook.Update(msg) // todo: handle related deep bug here
	return m, tea.Batch(objCmd, ovrCmd, noteCmd)
}
