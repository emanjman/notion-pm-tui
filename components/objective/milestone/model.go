package milestone

import (
	"notion-project-tui/notion"
	listutil "notion-project-tui/util/list"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type TaskViewMsg struct {
	Groups notion.TaskGroups
}

type Model struct {
	notion *notion.Client
	projID string
	propID string

	list    list.Model
	err     error
	loading bool
	groups map[string][]Item // milestones grouped by milestone status
	hidden map[string]bool

	ActiveKeyMap  help.KeyMap // for help focus view
	neutralKeyMap NeutralKeyMap
	writingKeyMap WritingKeyMap

	Focus *FocusState
}

var statusOrder = []string{"🚧 under development", "😴 idle", "🎉 complete"}

func New(n *notion.Client, projID, propID string) Model {
	f := FocusState{}

	l := list.New([]list.Item{}, NewItemDelegate(true, &f), 0, 0)
	l.Title = "Milestones"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	l.SetShowFilter(false)
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()

	m := Model{
		notion: n,
		projID: projID,
		propID: propID,

		list:    l,
		err:     nil,
		loading: true,
		groups:  map[string][]Item{},
		hidden:  map[string]bool{},

		ActiveKeyMap:  NeutralKeyMapper, // default map view
		neutralKeyMap: NeutralKeyMapper,
		writingKeyMap: WritingKeyMapper,

		Focus: &f,
	}
	m.list.SetItems(listutil.BuildGroupList(m.groups, m.hidden, statusOrder))

	return m
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		ids, err := m.notion.FetchRelationIDs(m.projID, m.propID)
		return notion.MilestoneIDsMsg{IDs: ids, Err: err}
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	case notion.MilestoneIDsMsg:
		if msg.Err != nil {
			m.err = msg.Err
			m.loading = false
			return m, nil
		}
		return m, func() tea.Msg {
			pages, err := notion.FetchPages[notion.MilestonePage](m.notion, msg.IDs)
			return notion.MilestonePagesMsg{Pages: pages, Err: err}
		}

	case notion.MilestonePagesMsg:
		if msg.Err != nil {
			m.err = msg.Err
			m.loading = false
			return m, nil
		}

		// create the list items
		tempItems := make([]Item, len(msg.Pages))
		for i, pg := range msg.Pages {
			tempItems[i] = NewItem(pg)
		}

		m.groups = listutil.GroupByKey(tempItems)
		items := listutil.BuildGroupList(m.groups, m.hidden, statusOrder)

		m.list.SetItems(items)
		m.loading = false

		cmds := []tea.Cmd{}
		for i, item := range m.list.Items() {
			if item, ok := item.(Item); ok && item.Status == "🚧 under development" {
				item.FetchState = Pending
				m.list.SetItem(i, item)
				cmds = append(cmds, m.queryTasksByStatus(i, item.ID, "dev", ""))
				cmds = append(cmds, m.queryTasksByStatus(i, item.ID, "idle", ""))
				cmds = append(cmds, m.queryTasksByStatus(i, item.ID, "done", ""))
				cmds = append(cmds, m.queryTasksByStatus(i, item.ID, "archive", ""))
			}
		}
		return m, tea.Batch(cmds...)

	case notion.FetchMoreTasksMsg:
		item := m.list.SelectedItem()
		if mstone, ok := item.(Item); ok {
			group := mstone.TaskGroups[msg.Status]
			if group.NextCursor != nil && !group.Loading {
				idx := m.list.Index()
				cursor := *group.NextCursor
				group.Loading = true
				mstone.TaskGroups[msg.Status] = group
				m.list.SetItem(idx, mstone)
				return m, tea.Batch(
					m.queryTasksByStatus(idx, mstone.ID, msg.Status, cursor),
					func() tea.Msg { return TaskViewMsg{Groups: mstone.TaskGroups} },
				)
			}
		}
		return m, nil

	case notion.ToggleTaskGroupMsg:
		item := m.list.SelectedItem()
		if mstone, ok := item.(Item); ok {
			group := mstone.TaskGroups[msg.Status]
			group.Hide = !group.Hide
			mstone.TaskGroups[msg.Status] = group
			m.updateMilestoneInGroups(mstone)
			return m, func() tea.Msg {
				return TaskViewMsg{Groups: mstone.TaskGroups}
			}
		}
		return m, nil

	case notion.TaskQueryMsg:
		if msg.Err != nil {
			m.err = msg.Err
			return m, nil
		}
		item := m.list.Items()[msg.MilestoneIdx]
		if mstone, ok := item.(Item); ok {
			group := mstone.TaskGroups[msg.Status]
			group.Tasks = append(group.Tasks, msg.Pages...)
			group.NextCursor = msg.NextCursor
			group.Loading = false
			mstone.TaskGroups[msg.Status] = group
			mstone.FetchState = Success
			m.list.SetItem(msg.MilestoneIdx, mstone)

			return m, func() tea.Msg {
				return TaskViewMsg{Groups: mstone.TaskGroups}
			}
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		// handle title editing (writing mode)
		if m.Focus.Mode == WritingMode {
			switch {
			case key.Matches(msg, m.writingKeyMap.Save):
				// update item in list
				if milestone, ok := m.list.SelectedItem().(Item); ok {
					milestone.Name = m.Focus.tempTitle.Value()
					m.list.SetItem(m.list.Index(), milestone)
					m.updateMilestoneInGroups(milestone)
				}

				m.ActiveKeyMap = NeutralKeyMapper
				m.Focus.Mode = NeutralMode

				// todo: send command to update milestone title in notion
				return m, nil

			// forward all keys into the textinput model
			default:
				var cmd tea.Cmd
				m.Focus.tempTitle, cmd = m.Focus.tempTitle.Update(msg)
				return m, cmd
			}
		}

		if m.Focus.Mode == NeutralMode {
			switch {
			case key.Matches(msg, m.neutralKeyMap.Select):
				selected := m.list.SelectedItem()

				// if selected item is header, toggle + rebuild list
				if header, ok := selected.(listutil.ListItemGroupHeader); ok {
					m.hidden[header.Label] = !m.hidden[header.Label]
					m.list.SetItems(listutil.BuildGroupList(m.groups, m.hidden, statusOrder))
				}
				return m, nil

			case key.Matches(msg, m.neutralKeyMap.Rename):
				if milestone, ok := m.list.SelectedItem().(Item); ok {
					m.Focus.milestoneID = milestone.ID
					m.Focus.milestoneIdx = m.list.Index()
					m.Focus.tempTitle = initTempTitle(milestone)
					m.Focus.Mode = WritingMode
					m.ActiveKeyMap = WritingKeyMapper
				}
				return m, nil

			case key.Matches(msg, m.neutralKeyMap.Down):
				m.list.CursorDown()
				return m, func() tea.Msg {
					return TaskViewMsg{Groups: m.getCurrTaskGroups()}
				}
			case key.Matches(msg, m.neutralKeyMap.Up):
				m.list.CursorUp()
				return m, func() tea.Msg {
					return TaskViewMsg{Groups: m.getCurrTaskGroups()}
				}
			case key.Matches(msg, m.neutralKeyMap.JumpDown):
				m.list.Select(min(len(m.list.Items())-1, m.list.Index()+5))
				return m, func() tea.Msg {
					return TaskViewMsg{Groups: m.getCurrTaskGroups()}
				}
			case key.Matches(msg, m.neutralKeyMap.JumpUp):
				m.list.Select(max(0, m.list.Index()-5))
				return m, func() tea.Msg {
					return TaskViewMsg{Groups: m.getCurrTaskGroups()}
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg) // handles up/down nav
	return m, cmd
}

// just forward the list.View()
func (m Model) View() string {
	if m.loading {
		return "Loading milestones..."
	}

	return m.list.View()
}

func (m Model) getCurrTaskGroups() notion.TaskGroups {
	item := m.list.SelectedItem()

	switch item := item.(type) {
	case listutil.ListItemGroupHeader:
		mstone := m.groups[item.Label][0]
		return mstone.TaskGroups
	case Item:
		return item.TaskGroups
	}
	return notion.TaskGroups{}
}

func (m *Model) SetItemDelegate(d list.ItemDelegate) {
	m.list.SetDelegate(d)
}

func (m Model) updateMilestoneInGroups(updated Item) Model {
	group := m.groups[updated.Status]

	// overwrite task in m.groups
	for i, t := range group {
		if t.ID == updated.ID {
			m.groups[updated.Status][i] = updated
			break
		}
	}

	// then rebuild item list
	m.list.SetItems(listutil.BuildGroupList(m.groups, m.hidden, statusOrder))
	return m
}

func (m Model) queryTasksByStatus(idx int, milestoneID, status, cursor string) tea.Cmd {
	return m.notion.QueryTasks(milestoneID, status, cursor, idx)
}
