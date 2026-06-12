package version

import (
	"notion-project-tui/notion"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	notion *notion.Client
	projID string
	err    error

	ActiveKeyMap  help.KeyMap // for help focus view
	neutralKeyMap NeutralKeyMap
}

func New(n *notion.Client, projID string) Model {
	return Model{
		notion: n,
		projID: projID,
		err:    nil,

		ActiveKeyMap:  NeutralKeyMapper, // default map view
		neutralKeyMap: NeutralKeyMapper,
	}
}

func (m Model) Init() tea.Cmd {
	return fetchVersions(m.projID, m.notion)
}
