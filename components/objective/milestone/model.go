package milestone

import (
	"notion-project-tui/notion"
	listutil "notion-project-tui/util/list"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type TaskViewMsg struct {
	Tasks []notion.TaskPage
}

type Model struct {
	notion *notion.Client
	projID string
	propID string

	list    list.Model
	err     error
	loading bool
	groups  map[string][]Item // header to items
	hidden  map[string]bool   // hidden group

	ActiveKeyMap    help.KeyMap // for help focus view
	neutralKeyMap   NeutralKeyMap
	selectingKeyMap SelectingKeyMap
	writingKeyMap   WritingKeyMap

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
		groups:  listutil.GroupByKey(mockItems()),
		hidden:  map[string]bool{},

		ActiveKeyMap:    NeutralKeyMapper, // default map view
		neutralKeyMap:   NeutralKeyMapper,
		selectingKeyMap: SelectingKeyMapper,
		writingKeyMap:   WritingKeyMapper,

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
				cmds = append(cmds, m.fetchTaskRelationIDs(i, item))
			}
		}
		return m, tea.Batch(cmds...)

	case notion.TaskIDsMsg:
		if msg.Err != nil {
			m.err = msg.Err
			m.loading = false
			return m, nil
		}
		item := m.list.Items()[msg.MilestoneIdx]
		if _, ok := item.(Item); ok {
			return m, func() tea.Msg {
				pages, err := notion.FetchPages[notion.TaskPage](m.notion, msg.IDs)
				return notion.TaskPagesMsg{Pages: pages, Err: err, MilestoneIdx: msg.MilestoneIdx}
			}
		}
		return m, nil

	case notion.TaskPagesMsg:
		if msg.Err != nil {
			m.err = msg.Err
			m.loading = false
			return m, nil
		}
		item := m.list.Items()[msg.MilestoneIdx]
		if mstone, ok := item.(Item); ok {
			mstone.Tasks = msg.Pages
			mstone.FetchState = Success
			m.list.SetItem(msg.MilestoneIdx, mstone)
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

		if m.Focus.Mode == SelectingMode {
			switch {

			// on exit, save updates via notion api
			case key.Matches(msg, m.selectingKeyMap.Exit):
				m.Focus.Mode = NeutralMode
				m.ActiveKeyMap = NeutralKeyMapper

				// todo: send command to update milestone changes in notion
				return m, nil

			// switch between fields (vertical navigation)
			case key.Matches(msg, m.selectingKeyMap.Up):
				if m.Focus.field == MilestoneTitle {
					m.Focus.field = fieldCnt - 1
				} else {
					m.Focus.field = (m.Focus.field - 1) % fieldCnt
				}
				return m, nil
			case key.Matches(msg, m.selectingKeyMap.Down):
				m.Focus.field = (m.Focus.field + 1) % fieldCnt
				return m, nil

			// enter field: cycle select or enter writing mode
			case key.Matches(msg, m.selectingKeyMap.Select):
				selected := m.list.SelectedItem()
				if milestone, ok := selected.(Item); ok {
					switch m.Focus.field {
					case MilestoneTag:
						// cycle tag, stay in selecting mode
						milestone.Tag = cycleTagField(milestone.Tag, 1)
						m.list.SetItem(m.Focus.milestoneIdx, milestone)
						m.updateMilestoneInGroups(milestone)
					case MilestoneTitle:
						// enter writing mode for title
						m.Focus.Mode = WritingMode
						m.ActiveKeyMap = WritingKeyMapper

						if item, ok := m.list.SelectedItem().(Item); ok {
							m.Focus.tempTitle = initTempTitle(item)
						}
					}

					return m, nil
				}
			}

			// consume all keys, don't forward to list navigations
			return m, nil
		}

		if m.Focus.Mode == NeutralMode {
			switch {
			case key.Matches(msg, m.neutralKeyMap.Select):
				selected := m.list.SelectedItem()

				// if selected item is header, toggle + rebuild list
				if header, ok := selected.(listutil.ListItemGroupHeader); ok {
					m.hidden[header.Label] = !m.hidden[header.Label]
					m.list.SetItems(listutil.BuildGroupList(m.groups, m.hidden, statusOrder))
				} else if milestone, ok := selected.(Item); ok {
					// initialize the focus state
					m.Focus.milestoneID = milestone.ID
					m.Focus.milestoneIdx = m.list.Index()
					m.Focus.field = MilestoneTitle // default field

					m.ActiveKeyMap = SelectingKeyMapper
					m.Focus.Mode = SelectingMode
				}
				return m, nil

			case key.Matches(msg, m.neutralKeyMap.Down):
				m.list.CursorDown()
				return m, func() tea.Msg {
					return TaskViewMsg{Tasks: m.getCurrTasks()}
				}
			case key.Matches(msg, m.neutralKeyMap.Up):
				m.list.CursorUp()
				return m, func() tea.Msg {
					return TaskViewMsg{Tasks: m.getCurrTasks()}
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

	containerStyle := lg.NewStyle().PaddingRight(1)
	return containerStyle.Render(m.list.View())
}

func (m Model) getCurrTasks() []notion.TaskPage {
	item := m.list.SelectedItem()

	switch item := item.(type) {
	case listutil.ListItemGroupHeader:
		mstone := m.groups[item.Label][0]
		return mstone.Tasks
	case Item:
		return item.Tasks
	}
	return []notion.TaskPage{}
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

func (m Model) fetchTaskRelationIDs(idx int, mstone Item) tea.Cmd {
	return func() tea.Msg {
		ids, err := m.notion.FetchRelationIDs(mstone.ID, mstone.TasksPropID)
		return notion.TaskIDsMsg{IDs: ids, Err: err, MilestoneIdx: idx}
	}
}
