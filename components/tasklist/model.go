package tasklist

import (
	"notion-project-tui/notion"
	listutil "notion-project-tui/util/list"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type EditState struct {
	active    bool
	taskID    string
	taskIndex int
	field     FieldIndex
	subActive bool

	TextInput    textinput.Model
	TempType     string
	TempPriority int
}

type FieldIndex int

const (
	TypeField = iota
	TaskField
	PriorityField
)
const fieldCnt = 3

func cycleType(current string, delta int) string {
	opts := notion.TypeSelectValues
	for i, t := range opts {
		if t == current {
			n := len(opts)
			return opts[((i+delta)%n+n)%n]
		}
	}
	return opts[0]
}

func cyclePriority(current, delta int) int {
	const n = 6
	return ((current+delta)%n + n) % n
}

func commitSubEdit(m TaskListModel) TaskListModel {
	items := m.list.Items()
	item, ok := items[m.EditState.taskIndex].(TaskListItem)
	if !ok {
		return m
	}
	switch m.EditState.field {
	case TypeField:
		item.Type = m.EditState.TempType
	case PriorityField:
		item.Priority = m.EditState.TempPriority
	case TaskField:
		item.Task = m.EditState.TextInput.Value()
	}
	m.list.SetItem(m.EditState.taskIndex, item)
	return m
}

// ? required structure for grouping (groups, hidden, ...) could be an interface
type TaskListModel struct {
	Milestone notion.SelectedMilestone
	list      list.Model
	loading   bool
	groups    map[string][]TaskListItem
	hidden    map[string]bool
	Keys      KeyMap
	EditKeys  EditKeyMap
	client    *notion.Client
	EditState *EditState
	// todo: cached milestones
}

var statusOrder = []string{"dev", "idle", "done"}

func NewTaskListModel(milestone notion.SelectedMilestone, c *notion.Client) TaskListModel {
	edit := &EditState{}

	l := list.New([]list.Item{}, NewTaskListDelegate(false, edit), 0, 0)
	l.Title = "Tasks"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)

	m := TaskListModel{
		Milestone: milestone,
		list:      l,
		loading:   true,
		groups:    listutil.GroupByKey(mockTaskItems()),
		hidden:    map[string]bool{},
		Keys:      DefaultKeyMap,
		EditKeys:  DefaultEditKeyMap,
		client:    c,
		EditState: edit,
	}
	m.list.SetItems(listutil.BuildGroupList(m.groups, m.hidden, statusOrder))

	return m
}

func (m TaskListModel) IsCapturingTextInput() bool {
	return m.EditState.active && m.EditState.subActive && m.EditState.field == TaskField
}

func (m TaskListModel) Update(msg tea.Msg) (TaskListModel, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		if m.EditState.active {

			// --- Level 2: sub-edit engaged ---
			if m.EditState.subActive {
				switch m.EditState.field {

				case TypeField:
					switch {
					case key.Matches(msg, m.EditKeys.PrevField):
						m.EditState.TempType = cycleType(m.EditState.TempType, -1)
					case key.Matches(msg, m.EditKeys.NextField):
						m.EditState.TempType = cycleType(m.EditState.TempType, +1)
					case key.Matches(msg, m.EditKeys.EnableEdit), key.Matches(msg, m.EditKeys.Exit):
						m = commitSubEdit(m)
						m.EditState.subActive = false
					}

				case PriorityField:
					switch {
					case key.Matches(msg, m.EditKeys.PrevField):
						m.EditState.TempPriority = cyclePriority(m.EditState.TempPriority, -1)
					case key.Matches(msg, m.EditKeys.NextField):
						m.EditState.TempPriority = cyclePriority(m.EditState.TempPriority, +1)
					case key.Matches(msg, m.EditKeys.EnableEdit), key.Matches(msg, m.EditKeys.Exit):
						m = commitSubEdit(m)
						m.EditState.subActive = false
					}

				case TaskField:
					if key.Matches(msg, m.EditKeys.EnableEdit) || key.Matches(msg, m.EditKeys.Exit) {
						m = commitSubEdit(m)
						m.EditState.subActive = false
					} else {
						var cmd tea.Cmd
						m.EditState.TextInput, cmd = m.EditState.TextInput.Update(msg)
						return m, cmd
					}
				}

				return m, nil
			}

			// --- Level 1: field selection ---
			switch {

			case key.Matches(msg, m.EditKeys.Exit):
				m.EditState.active = false
				// todo: send some command to save notion changes
				return m, nil

			case key.Matches(msg, m.EditKeys.PrevField):
				if m.EditState.field == TypeField {
					m.EditState.field = fieldCnt - 1
				} else {
					m.EditState.field = (m.EditState.field - 1) % fieldCnt
				}
				return m, nil

			case key.Matches(msg, m.EditKeys.NextField):
				m.EditState.field = (m.EditState.field + 1) % fieldCnt
				return m, nil

			case key.Matches(msg, m.EditKeys.EnableEdit):
				item, ok := m.list.Items()[m.list.Index()].(TaskListItem)
				if !ok {
					return m, nil
				}
				m.EditState.taskIndex = m.list.Index()

				switch m.EditState.field {
				case TypeField:
					m.EditState.TempType = item.Type
				case PriorityField:
					m.EditState.TempPriority = item.Priority
				case TaskField:
					ti := textinput.New()
					ti.SetValue(item.Task)
					ti.CursorEnd()
					ti.Width = m.list.Width() - lg.Width(item.Type) - 1 - 3 - 7
					ti.Focus()
					m.EditState.TextInput = ti
				}

				m.EditState.subActive = true
				return m, nil
			}

			// consume all keys, don't forward to list navigation
			return m, nil

		} else {
			switch {
			case key.Matches(msg, m.Keys.Select):
				selected := m.list.SelectedItem()

				// if selected item is header, toggle + rebuild list
				if header, ok := selected.(listutil.ListItemGroupHeader); ok {
					m.hidden[header.Label] = !m.hidden[header.Label]
					m.list.SetItems(listutil.BuildGroupList(m.groups, m.hidden, statusOrder))
				} else if task, ok := selected.(TaskListItem); ok {
					m.EditState.taskID = task.ID
					m.EditState.active = true
					m.EditState.field = TaskField
				}
				return m, nil
			}
		}

	case notion.MilestoneSelectedMsg:
		m.Milestone.ID = msg.Milestone.ID
		m.Milestone.TasksPropID = msg.Milestone.TasksPropID

		return m, m.client.FetchTaskRelationIds(m.Milestone.ID, m.Milestone.TasksPropID)

	case notion.TaskRelationIdsMsg:
		if msg.Err != nil {
			return m, nil
		}

		return m, m.client.FetchTasks(msg.IDs)

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
