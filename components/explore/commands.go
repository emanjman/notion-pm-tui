package explore

import (
	"notion-project-tui/notion"

	tea "github.com/charmbracelet/bubbletea"
)

func fetchProjects(ntn *notion.Client) tea.Cmd {
	return ntn.QueryProjectPages("")
}

// trickle down projectID to be caught by dependent child models
func emitProjectID(ID string) tea.Cmd {
	return func() tea.Msg {
		return notion.ProjectIDMsg{ID: ID}
	}
}
