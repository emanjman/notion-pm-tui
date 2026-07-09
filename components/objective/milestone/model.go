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

	ActiveKeyMap help.KeyMap // for help focus view
	normalKeyMap NormalKeyMap
	editKeyMap   EditKeyMap
	deleteKeyMap DeleteKeyMap
}

func New(n *notion.Client) Model {
	mode := NormalMode
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

		ActiveKeyMap: NormalKeyMapper, // default map view
		normalKeyMap: NormalKeyMapper,
		editKeyMap:   EditKeyMapper,
		deleteKeyMap: DeleteKeyMapper,

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
		case NormalMode:
			switch {
			// navigation
			case key.Matches(msg, m.normalKeyMap.Down):
				return m.handleNormalDown()
			case key.Matches(msg, m.normalKeyMap.Up):
				return m.handleNormalUp()
			case key.Matches(msg, m.normalKeyMap.JumpDown):
				return m.handleNormalJumpDown()
			case key.Matches(msg, m.normalKeyMap.JumpUp):
				return m.handleNormalJumpUp()

			// change modes
			case key.Matches(msg, m.normalKeyMap.Rename):
				return m.handleNormalRename()
			case key.Matches(msg, m.normalKeyMap.Add):
				return m.handleNormalAdd()
			case key.Matches(msg, m.normalKeyMap.Delete):
				return m.handleNormalDelete()

			// dynamic: change mode, launch fetches
			case key.Matches(msg, m.normalKeyMap.Select):
				return m.handleNormalSelect()
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
