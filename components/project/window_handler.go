package project

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleWindow(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	// todo: never handled bc window size msg occurs before a project even exists, it's not re-sent
	// todo: so a proj model should always exist, but sustain a null internal state until resolved
	log.Printf("project update catches window size msg") // !debug

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
}
