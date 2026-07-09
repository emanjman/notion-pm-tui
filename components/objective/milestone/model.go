package milestone

import (
	"notion-project-tui/notion"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	notion         *notion.Client
	versionID      string
	list           list.Model
	err            error
	pendingFetches int
	groups         notion.MilestoneGroups

	Mode   *Mode        // ptr so the list delegate sees mode switches live (shared heap)
	Edit   *EditModeCtx // todo: do these NEED to be ptrs?
	Delete *DeleteModeCtx

	ActiveKeyMap  help.KeyMap // for help focus view
	neutralKeyMap NeutralKeyMap
	editKeyMap    EditKeyMap
	deleteKeyMap  DeleteKeyMap
}

func New(n *notion.Client) Model {
	mode := NeutralMode
	edit := EditModeCtx{}
	del := DeleteModeCtx{}

	l := list.New([]list.Item{}, NewItemDelegate(true, &mode, &edit), 0, 0)
	l.Title = "Milestones"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	l.SetShowFilter(false)
	l.SetShowPagination(false)
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()

	return Model{
		notion:         n,
		versionID:      "", // initially, no version selected, injected after version selected
		pendingFetches: 3,
		list:           l,
		err:            nil,
		groups:         notion.MilestoneGroups{},

		ActiveKeyMap:  NeutralKeyMapper, // default map view
		neutralKeyMap: NeutralKeyMapper,
		editKeyMap:    EditKeyMapper,
		deleteKeyMap:  DeleteKeyMapper,

		Mode:   &mode,
		Edit:   &edit,
		Delete: &del,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

// dispatch messages to handlers
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	// handle queries
	case notion.QueryMilestonePagesMsg:
		return m.handleQueryMilestonePages(msg)
	case notion.QueryMoreMilestonePagesMsg:
		return m.handleQueryMoreMilestonePages(msg)
	case notion.QueryTaskPagesMsg:
		return m.handleQueryTaskPages(msg)
	case notion.QueryMoreTaskPagesMsg:
		return m.handleQueryMoreTaskPages(msg)

	// handle mutations
	case notion.AddMilestonePageMsg:
		return m.handleAddMilestonePage(msg)
	case UpdateNotionTitleMsg:
		return m.handleUpdateNotionTitle(msg)
	case notion.ToggleTaskGroupMsg:
		return m.handleToggleTaskGroup(msg)
	case TrashMilestonePageMsg:
		return m.handleTrashMilestonePage(msg)

	// handle keys
	case tea.KeyMsg:
		switch *m.Mode {
		case NeutralMode:
			switch {
			// navigation
			case key.Matches(msg, m.neutralKeyMap.Down):
				return m.handleNeutralDown()
			case key.Matches(msg, m.neutralKeyMap.Up):
				return m.handleNeutralUp()
			case key.Matches(msg, m.neutralKeyMap.JumpDown):
				return m.handleNeutralJumpDown()
			case key.Matches(msg, m.neutralKeyMap.JumpUp):
				return m.handleNeutralJumpUp()

			// change modes
			case key.Matches(msg, m.neutralKeyMap.Rename):
				return m.handleNeutralRename()
			case key.Matches(msg, m.neutralKeyMap.Add):
				return m.handleNeutralAdd()
			case key.Matches(msg, m.neutralKeyMap.Delete):
				return m.handleNeutralDelete()

			// dynamic: change mode, launch fetches
			case key.Matches(msg, m.neutralKeyMap.Select):
				return m.handleNeutralSelect()
			}

		case EditMode:
			switch {
			case key.Matches(msg, m.editKeyMap.Save):
				return m.handleEditSave()
			default:
				var cmd tea.Cmd
				m.Edit.titleInput, cmd = m.Edit.titleInput.Update(msg)
				return m, cmd
			}

		case DeleteMode:
			switch {
			case key.Matches(msg, m.deleteKeyMap.Cancel):
				return m.handleDeleteCancel()
			case key.Matches(msg, m.deleteKeyMap.Confirm):
				return m.handleDeleteConfirm()
			default:
				return m.handleDeleteCancel()
			}
		}

	// handle window
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
	}

	// otherwise, handle from children
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.pendingFetches > 0 {
		return "Loading milestones..."
	}
	return m.list.View()
}
