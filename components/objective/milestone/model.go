package milestone

import (
	"notion-project-tui/notion"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type TaskViewMsg struct {
	Groups notion.TaskGroups
}

type Model struct {
	notion         *notion.Client
	projID         string
	list           list.Model
	err            error
	pendingFetches int
	groups         notion.MilestoneGroups
	ActiveKeyMap   help.KeyMap // for help focus view
	neutralKeyMap  NeutralKeyMap
	writingKeyMap  WritingKeyMap
	Focus          *FocusState
}

func New(n *notion.Client, projID string) Model {
	f := FocusState{}

	l := list.New([]list.Item{}, NewItemDelegate(true, &f), 0, 0)
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
		projID:         projID,
		pendingFetches: 3,
		list:           l,
		err:            nil,
		groups:         notion.MilestoneGroups{},
		ActiveKeyMap:   NeutralKeyMapper, // default map view
		neutralKeyMap:  NeutralKeyMapper,
		writingKeyMap:  WritingKeyMapper,
		Focus:          &f,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.notion.QueryMilestones(m.projID, notion.MilestoneUnderDevelopment, ""),
		m.notion.QueryMilestones(m.projID, notion.MilestoneIdle, ""),
		m.notion.QueryMilestones(m.projID, notion.MilestoneComplete, ""),
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	case notion.MilestonePagesMsg:
		if msg.Err != nil {
			m.err = msg.Err
			m.pendingFetches--
			return m, nil
		}

		// append incoming pages into the correct status group
		group := m.groups[msg.Status]
		group.Milestones = append(group.Milestones, msg.Pages...)
		group.NextCursor = msg.NextCursor
		m.groups[msg.Status] = group

		m.pendingFetches--

		// only render + kick off task fetches once all 3 status batches have arrived
		if m.pendingFetches > 0 {
			return m, nil
		}

		m.list.SetItems(m.buildMilestoneList())

		cmds := []tea.Cmd{}
		for i, item := range m.list.Items() {
			if item, ok := item.(DefaultItem); ok && item.MilestoneStatus == notion.MilestoneUnderDevelopment {
				item.FetchStatus = FetchPending
				m.list.SetItem(i, item)
				cmds = append(cmds, m.queryTasksByStatus(i, item.ID, "dev", ""))
				cmds = append(cmds, m.queryTasksByStatus(i, item.ID, "idle", ""))
				cmds = append(cmds, m.queryTasksByStatus(i, item.ID, "done", ""))
				cmds = append(cmds, m.queryTasksByStatus(i, item.ID, "archive", ""))
			}
		}
		return m, tea.Batch(cmds...)

	case notion.FetchMoreMilestonesMsg:
		group := m.groups[msg.Status]
		if group.NextCursor != nil && !group.Loading {
			cursor := *group.NextCursor
			group.Loading = true
			m.groups[msg.Status] = group
			m.list.SetItems(m.buildMilestoneList())
			return m, m.notion.QueryMilestones(m.projID, msg.Status, cursor)

		}
		return m, nil

	case notion.FetchMoreTasksMsg:
		item := m.list.SelectedItem()
		if mstone, ok := item.(DefaultItem); ok {
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
		if mstone, ok := item.(DefaultItem); ok {
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
		if mstone, ok := item.(DefaultItem); ok {
			group := mstone.TaskGroups[msg.Status]
			group.Tasks = append(group.Tasks, msg.Pages...)
			group.NextCursor = msg.NextCursor
			group.Loading = false
			mstone.TaskGroups[msg.Status] = group
			mstone.FetchStatus = FetchSuccess
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
				if milestone, ok := m.list.SelectedItem().(DefaultItem); ok {
					milestone.Name = m.Focus.tempTitle.Value()
					m.list.SetItem(m.list.Index(), milestone)
					m.updateMilestoneInGroups(milestone)
				}

				m.ActiveKeyMap = NeutralKeyMapper
				m.Focus.Mode = NeutralMode

				// todo: send command to update milestone title in notion
				return m, nil

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

				if header, ok := selected.(GroupHeaderItem); ok {
					// toggle header
					group := m.groups[header.Status]
					group.Hide = !group.Hide
					m.groups[header.Status] = group
					m.list.SetItems(m.buildMilestoneList())

				} else if loadMore, ok := selected.(LoadMoreItem); ok && !loadMore.Loading {
					// load more milestones
					return m, func() tea.Msg {
						return notion.FetchMoreMilestonesMsg{Status: loadMore.Status}
					}
				} else if mstone, ok := selected.(DefaultItem); ok {
					// fetch tasks for curr milestone
					switch mstone.FetchStatus {
					case FetchIdle:
						idx := m.list.Index()
						mstone.FetchStatus = FetchPending
						m.list.SetItem(idx, mstone)
						if mstone.TaskCount > 0 {
							return m, tea.Batch(
								m.queryTasksByStatus(idx, mstone.ID, "dev", ""),
								m.queryTasksByStatus(idx, mstone.ID, "idle", ""),
								m.queryTasksByStatus(idx, mstone.ID, "done", ""),
								m.queryTasksByStatus(idx, mstone.ID, "archive", ""),
							)
						}
						return m, nil
					// todo: is this necessary
					case FetchSuccess:
						return m, func() tea.Msg {
							return TaskViewMsg{Groups: mstone.TaskGroups}
						}
					}
				}
				return m, nil

			case key.Matches(msg, m.neutralKeyMap.Rename):
				if milestone, ok := m.list.SelectedItem().(DefaultItem); ok {
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
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.pendingFetches > 0 {
		return "Loading milestones..."
	}
	return m.list.View()
}

func (m Model) buildMilestoneList() []list.Item {
	var items []list.Item
	for _, status := range notion.MilestoneStatusOrder() {
		group, ok := m.groups[status]
		if !ok || len(group.Milestones) == 0 {
			continue
		}
		items = append(items, NewGroupHeaderItem(status, group))
		if !group.Hide {
			for _, pg := range group.Milestones {
				items = append(items, NewDefaultItem(pg))
			}
			if group.NextCursor != nil {
				items = append(items, NewLoadMoreItem(status, group))
			}
		}
	}
	return items
}

func (m Model) getCurrTaskGroups() notion.TaskGroups {
	item := m.list.SelectedItem()

	switch item := item.(type) {
	case GroupHeaderItem:
		group := m.groups[item.Status]
		if len(group.Milestones) > 0 {
			return NewDefaultItem(group.Milestones[0]).TaskGroups
		}
	case DefaultItem:
		return item.TaskGroups
	}
	return notion.TaskGroups{}
}

func (m *Model) SetItemDelegate(d list.ItemDelegate) {
	m.list.SetDelegate(d)
}

func (m Model) updateMilestoneInGroups(item DefaultItem) Model {
	group := m.groups[item.MilestoneStatus]

	for i, pg := range group.Milestones {
		if pg.ID == item.ID {
			// sync the name back onto the page (only field editable locally)
			group.Milestones[i].Properties.Title.Title[0].PlainText = item.Name
			break
		}
	}

	m.groups[item.MilestoneStatus] = group
	m.list.SetItems(m.buildMilestoneList())
	return m
}

func (m Model) queryTasksByStatus(idx int, milestoneID, status, cursor string) tea.Cmd {
	return m.notion.QueryTasks(milestoneID, status, cursor, idx)
}
