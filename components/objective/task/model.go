package task

import (
	"fmt"
	"notion-project-tui/components/objective/milestone"
	"notion-project-tui/notion"
	listutil "notion-project-tui/util/list"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type Model struct {
	// Milestone notion.SelectedMilestone // milestone these tasks source from
	list    list.Model
	notion  *notion.Client
	loading bool

	groups map[string][]Item // should this be Groupable intf?
	hidden map[string]bool   // should this be Groupable intf?

	ActiveKeyMap    help.KeyMap // for help focus view
	neutralKeyMap   NeutralKeyMap
	selectingKeyMap SelectingKeyMap
	writingKeyMap   WritingKeyMap

	Focus *FocusState

	tempIDCounter int // for generating temp IDs for new tasks
}

var statusOrder = []string{"dev", "idle", "done"}

func New(clt *notion.Client) Model {
	f := FocusState{}

	l := list.New([]list.Item{}, NewItemDelegate(false, &f), 0, 0)
	l.Title = "Tasks"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	l.SetShowFilter(false)
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()

	return Model{
		// Milestone: mstone,
		list:    l,
		notion:  clt,
		loading: true,

		groups: map[string][]Item{},
		hidden: map[string]bool{},

		ActiveKeyMap:    NeutralKeyMapper, // default map view
		neutralKeyMap:   NeutralKeyMapper,
		selectingKeyMap: SelectingKeyMapper,
		writingKeyMap:   WritingKeyMapper,

		Focus: &f,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	case milestone.TaskViewMsg:
		// create list items
		tempItems := make([]Item, len(msg.Tasks))
		for i, page := range msg.Tasks {
			tempItems[i] = NewItem(page)
		}
		m.groups = listutil.GroupByKey(tempItems)
		items := listutil.BuildGroupList(m.groups, m.hidden, statusOrder)
		m.list.SetItems(items)
		m.loading = false

		return m, nil

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		// handle edits by field
		if m.Focus.Mode == WritingMode {
			switch {
			case key.Matches(msg, m.writingKeyMap.Save):
				// update item in list
				if task, ok := m.list.SelectedItem().(Item); ok {
					task.Task = m.Focus.tempTitle.Value()
					m.list.SetItem(m.Focus.taskIdx, task)
					m.updateTaskInGroups(task)
				}

				m.ActiveKeyMap = NeutralKeyMapper
				m.Focus.Mode = NeutralMode

				// todo: send command to update task title in notion
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

				// todo: send command to update task changes (type/priority) in notion
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
				if task, ok := selected.(Item); ok {
					switch m.Focus.field {
					case TaskType:
						task.Type = cycleTypeField(task.Type, 1)
					case TaskPriority:
						task.Priority = cyclePriorityField(task.Priority, 1)
					case TaskTitle:
						m.Focus.Mode = WritingMode
						m.ActiveKeyMap = WritingKeyMapper

						if item, ok := m.list.SelectedItem().(Item); ok {
							m.Focus.tempTitle = initTempTitle(item)
						}
					}

					m.list.SetItem(m.Focus.taskIdx, task)
					m.updateTaskInGroups(task)
					return m, nil
				}
			}

			// consume all keys, don't forward to list navigations
			return m, nil
		}

		if m.Focus.Mode == NeutralMode {
			// Cancel pending delete if any key other than Delete is pressed
			if m.Focus.pendingDelete && !key.Matches(msg, m.neutralKeyMap.Delete) {
				m.Focus.pendingDelete = false
				// Consume the key event to prevent it from being forwarded
				return m, nil
			}

			switch {
			case key.Matches(msg, m.neutralKeyMap.Select):
				selected := m.list.SelectedItem()

				// if selected item is header, toggle + rebuild list
				if header, ok := selected.(listutil.ListItemGroupHeader); ok {
					m.hidden[header.Label] = !m.hidden[header.Label]
					m.list.SetItems(listutil.BuildGroupList(m.groups, m.hidden, statusOrder))
				} else if task, ok := selected.(Item); ok {
					// initialize the focus state
					m.Focus.taskID = task.ID
					m.Focus.taskIdx = m.list.Index()
					m.Focus.field = TaskTitle // default field

					m.ActiveKeyMap = SelectingKeyMapper
					m.Focus.Mode = SelectingMode
				}

				return m, nil

			case key.Matches(msg, m.neutralKeyMap.AddTask):
				// always add to idle group
				m = m.addTask()
				return m, nil

			case key.Matches(msg, m.neutralKeyMap.StatusPrev):
				if task, ok := m.list.SelectedItem().(Item); ok {
					m = m.changeTaskStatus(task, -1)
					// todo: send command to update status in notion
				}
				return m, nil

			case key.Matches(msg, m.neutralKeyMap.StatusNext):
				if task, ok := m.list.SelectedItem().(Item); ok {
					m = m.changeTaskStatus(task, +1)
					// todo: send command to update status in notion
				}
				return m, nil

			case key.Matches(msg, m.neutralKeyMap.Delete):
				if task, ok := m.list.SelectedItem().(Item); ok {
					if m.Focus.pendingDelete && m.Focus.taskID == task.ID {
						// Second press: actually delete
						m = m.deleteTask(task)
					} else {
						// First press: set pending delete
						m.Focus.pendingDelete = true
						m.Focus.taskID = task.ID
						m.Focus.taskIdx = m.list.Index()
					}
				}
				return m, nil
			}
		}
	}

	// forward to children
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	// ! temp, styling
	// if m.loading {
	// 	return "Loading tasks..."
	// }

	containerStyle := lg.NewStyle().PaddingLeft(1)
	return containerStyle.Render(m.list.View())
}

// --- helpers

func (m Model) changeTaskStatus(task Item, delta int) Model {
	newStatus := cycleStatus(task.Status, delta)
	if newStatus == task.Status {
		return m // no change
	}

	// remove from old group
	oldGroup := m.groups[task.Status]
	for i, t := range oldGroup {
		if t.ID == task.ID {
			m.groups[task.Status] = append(oldGroup[:i], oldGroup[i+1:]...)
			break
		}
	}

	// update task status and add to new group
	task.Status = newStatus
	m.groups[newStatus] = append(m.groups[newStatus], task)

	// rebuild list to show change
	m.list.SetItems(listutil.BuildGroupList(m.groups, m.hidden, statusOrder))

	return m
}

func (m Model) updateTaskInGroups(updated Item) Model {
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

func (m Model) addTask() Model {
	defaultStatus := "idle"

	m.tempIDCounter++
	tempID := fmt.Sprintf("temp-%d", m.tempIDCounter)

	// new task with defaults
	newTask := Item{
		ID:       tempID,
		Task:     "",
		Status:   "idle",
		Priority: 0, // p0
		Type:     "feat",
	}

	// Add to group
	m.groups[defaultStatus] = append(m.groups[defaultStatus], newTask)

	// Rebuild list
	items := listutil.BuildGroupList(m.groups, m.hidden, statusOrder)
	m.list.SetItems(items)

	// Find the new task's index in the rebuilt list
	for i, item := range items {
		if task, ok := item.(Item); ok && task.ID == tempID {
			// Select the new task
			m.list.Select(i)

			// Initialize focus state
			m.Focus.taskID = tempID
			m.Focus.taskIdx = i
			m.Focus.field = TaskTitle

			// Initialize text input and enter writing mode
			m.Focus.tempTitle = initTempTitle(newTask)
			m.Focus.Mode = WritingMode
			m.ActiveKeyMap = m.writingKeyMap

			break
		}
	}

	return m
}

func (m Model) deleteTask(task Item) Model {
	// Remove from group
	group := m.groups[task.Status]
	for i, t := range group {
		if t.ID == task.ID {
			m.groups[task.Status] = append(group[:i], group[i+1:]...)
			break
		}
	}

	// Rebuild list
	m.list.SetItems(listutil.BuildGroupList(m.groups, m.hidden, statusOrder))

	// Clear pending delete state
	m.Focus.pendingDelete = false

	// todo: send command to delete task in notion

	return m
}

func (m *Model) SetItemDelegate(d list.ItemDelegate) {
	m.list.SetDelegate(d)
}
