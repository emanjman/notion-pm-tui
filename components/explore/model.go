package explore

import (
	tea "github.com/charmbracelet/bubbletea"
	"notion-project-tui/components/project"
	"notion-project-tui/notion"
)

type Model struct {
	notion  *notion.Client
	project project.Model
}

var _ tea.Model = (*Model)(nil) // conform

func New() Model {
	ntn := notion.NewClient()
	return Model{notion: ntn}
}

// initiate project fetch
// update will handle the selection to kickoff the selected project
func (m Model) Init() tea.Cmd {
	return nil
}
