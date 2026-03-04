package tasklist

import (
	"notion-project-tui/notion"
	listutil "notion-project-tui/util/list"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type EditState struct {
	active bool
	taskID string
	field  FieldIndex
}

type FieldIndex int

const (
	TypeField = iota
	TaskField
	PriorityField
)
const fieldCnt = 3

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
	EditState EditState
	// todo: cached milestones
}

var statusOrder = []string{"dev", "idle", "done"}

func NewTaskListModel(milestone notion.SelectedMilestone, c *notion.Client) TaskListModel {
	edit := EditState{}

	l := list.New([]list.Item{}, NewTaskListDelegate(false, &edit), 0, 0)
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

func (m TaskListModel) Update(msg tea.Msg) (TaskListModel, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		if m.EditState.active {
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
				// todo: enter 2-deep edit mode

			}

			// consume all keys, don't forward to list navigations
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
