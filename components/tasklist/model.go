package tasklist

import (
	"notion-project-tui/notion"
	listutil "notion-project-tui/util/list"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// ? required structure for grouping (groups, hidden, ...) could be an interface
type TaskListModel struct {
	Milestone notion.SelectedMilestone
	list      list.Model
	loading   bool
	groups    map[string][]TaskListItem
	hidden    map[string]bool
	Keys      KeyMap
	client    *notion.Client
	// todo: cached milestones
}

var statusOrder = []string{"dev", "idle", "done"}

func NewTaskListModel(milestone notion.SelectedMilestone, c *notion.Client) TaskListModel {
	l := list.New([]list.Item{}, NewTaskListDelegate(), 0, 0)
	l.Title = "Tasks"
	l.SetShowHelp(false)

	m := TaskListModel{
		Milestone: milestone,
		list:      l,
		loading:   true,
		groups:    listutil.GroupByKey(mockTaskItems()),
		hidden:    map[string]bool{},
		Keys:      DefaultKeyMap,
		client:    c,
	}
	m.list.SetItems(listutil.BuildGroupList(m.groups, m.hidden, statusOrder))

	return m
}

func (m TaskListModel) Update(msg tea.Msg) (TaskListModel, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Select):
			selected := m.list.SelectedItem()

			// if selected item is header, toggle + rebuild list
			if header, ok := selected.(listutil.ListItemGroupHeader); ok {
				m.hidden[header.Label] = !m.hidden[header.Label]
				m.list.SetItems(listutil.BuildGroupList(m.groups, m.hidden, statusOrder))
			}

			return m, nil
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

	return m.list.View()
}

