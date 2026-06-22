package version

import (
	"notion-project-tui/notion"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	notion  *notion.Client
	projID  string
	loading bool
	err     error

	pages       []notion.VersionPage
	PageIdx     int
	HesiPageIdx int // hesitating val to diff check

	Mode *Mode // in case we need to support more modes later

	ActiveKeyMap  help.KeyMap // for help focus view
	neutralKeyMap NeutralKeyMap

	width  int
	height int
}

func New(n *notion.Client, projID string) Model {
	mode := NeutralMode

	return Model{
		notion:  n,
		projID:  projID,
		loading: true,
		err:     nil,

		PageIdx:     0,
		HesiPageIdx: 0,

		Mode: &mode,

		ActiveKeyMap:  NeutralKeyMapper, // default map view
		neutralKeyMap: NeutralKeyMapper,

		width:  0,
		height: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return fetchVersions(m.projID, m.notion)
}
