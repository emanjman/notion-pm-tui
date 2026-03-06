package tasklist

import (
	"notion-project-tui/notion"
	listutil "notion-project-tui/util/list"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type TaskListModel struct {
	Milestone notion.SelectedMilestone // milestone these tasks source from
	list      list.Model
	notion    *notion.Client
	loading   bool

	groups map[string][]TaskListItem // should this be Groupable intf?
	hidden map[string]bool           // should this be Groupable intf?

	ActiveKeyMap    help.KeyMap // for help focus view
	neutralKeyMap   NeutralKeyMap
	selectingKeyMap SelectingKeyMap
	writingKeyMap   WritingKeyMap

	Focus *FocusState

	// todo: cached milestones
}

var statusOrder = []string{"dev", "idle", "done"}

func NewTaskListModel(mstone notion.SelectedMilestone, clt *notion.Client) TaskListModel {
	f := FocusState{}

	l := list.New([]list.Item{}, NewTaskListDelegate(false, &f), 0, 0)
	l.Title = "Tasks"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)

	m := TaskListModel{
		Milestone: mstone,
		list:      l,
		notion:    clt,
		loading:   true,

		groups: listutil.GroupByKey(mockTaskItems()), // ! fetch actual tasks
		hidden: map[string]bool{},

		ActiveKeyMap:    NeutralKeyMapper, // default map view
		neutralKeyMap:   NeutralKeyMapper,
		selectingKeyMap: SelectingKeyMapper,
		writingKeyMap:   WritingKeyMapper,

		Focus: &f,
	}

	m.list.SetItems(listutil.BuildGroupList(m.groups, m.hidden, statusOrder))

	return m
}

func (m TaskListModel) Update(msg tea.Msg) (TaskListModel, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		// handle edits by field
		if m.Focus.Mode == WritingMode {
			switch {
			case key.Matches(msg, m.writingKeyMap.Save):
				// update item in list
				if task, ok := m.list.SelectedItem().(TaskListItem); ok {
					task.Task = m.Focus.tempTitle.Value()
					m.list.SetItem(m.list.Index(), task)
				}

				m.ActiveKeyMap = NeutralKeyMapper
				m.Focus.Mode = NeutralMode

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

				// todo: send command to update task changes in notion
				return m, nil

			// switch between fields
			case key.Matches(msg, m.selectingKeyMap.Left):
				if m.Focus.field == TaskType {
					m.Focus.field = fieldCnt - 1
				} else {
					m.Focus.field = (m.Focus.field - 1) % fieldCnt
				}
				return m, nil
			case key.Matches(msg, m.selectingKeyMap.Right):
				m.Focus.field = (m.Focus.field + 1) % fieldCnt
				return m, nil

			// enter field edit mode, catch all keys from root; handle per field
			case key.Matches(msg, m.selectingKeyMap.Select):
				selected := m.list.SelectedItem()
				if task, ok := selected.(TaskListItem); ok {
					switch m.Focus.field {
					case TaskType:
						task.Type = cycleTypeField(task.Type, 1)
					case TaskPriority:
						task.Priority = cyclePriorityField(task.Priority, 1)
					case TaskTitle:
						m.Focus.Mode = WritingMode
						m.ActiveKeyMap = WritingKeyMapper

						if item, ok := m.list.SelectedItem().(TaskListItem); ok {
							m.Focus.tempTitle = initTempTitle(item)
						}
					}

					m.list.SetItem(m.Focus.taskIdx, task)
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
				} else if task, ok := selected.(TaskListItem); ok {
					// initialize the focus state
					m.Focus.taskID = task.ID
					m.Focus.taskIdx = m.list.Index()
					m.Focus.field = TaskTitle // default field

					m.ActiveKeyMap = SelectingKeyMapper
					m.Focus.Mode = SelectingMode
				}

				return m, nil
			}
		}

	case notion.MilestoneSelectedMsg:
		m.Milestone.ID = msg.Milestone.ID
		m.Milestone.TasksPropID = msg.Milestone.TasksPropID

		return m, m.notion.FetchTaskRelationIds(m.Milestone.ID, m.Milestone.TasksPropID)

	case notion.TaskRelationIdsMsg:
		if msg.Err != nil {
			return m, nil
		}

		return m, m.notion.FetchTasks(msg.IDs)

	case notion.TaskMsg:
		if msg.Err != nil {
			return m, nil
		}

		// create list items
		tempItems := make([]TaskListItem, len(msg.Data))
		for i, page := range msg.Data {
			tempItems[i] = NewTaskListItem(page)
		}

		m.groups = listutil.GroupByKey(tempItems)
		items := listutil.BuildGroupList(m.groups, m.hidden, statusOrder)

		m.list.SetItems(items)
		m.loading = false
	}

	// forward rest of commands to children models (list)
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg) // handles up/down nav
	return m, cmd
}

func (m TaskListModel) View() string {
	// ! temp, styling
	// if m.loading {
	// 	return "Loading tasks..."
	// }

	containerStyle := lg.NewStyle().PaddingLeft(1)
	return containerStyle.Render(m.list.View())
}

func (m *TaskListModel) SetItemDelegate(d list.ItemDelegate) {
	m.list.SetDelegate(d)
}
