package explore

import (
	"notion-project-tui/notion"

	tea "github.com/charmbracelet/bubbletea"
)

func fetchProjects(ntn *notion.Client) tea.Cmd {
	return ntn.QueryProjectPages("")
}
