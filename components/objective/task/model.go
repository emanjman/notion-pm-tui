package task

import (
	"fmt"
	"notion-project-tui/notion"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type UpdateTitleMsg struct{ Err error }
type UpdateSelectionsMsg struct{ Err error }
type UpdateStatusMsg struct{ Err error }

type Model struct {
	// Milestone notion.SelectedMilestone // milestone these tasks source from
	list    list.Model
	notion  *notion.Client
	loading bool

	typeOptions []notion.SelectItem

	// id of the milestone backing the current task list; new tasks hang off this
	// via the @milestone relation. set on every milestone switch (MilestoneTasksMsg).
	milestoneID string

	// working copy of the current milestone's tasks, used for local mutations
	// (add, delete, status change) and list rendering. rebuilt on every
	// milestone switch via MilestoneTasksMsg — not the source of truth for persistence.
	// note: local mutations here do not sync back to the milestone's TaskGroups.
	groups map[string][]Item

	ActiveKeyMap    help.KeyMap // for help focus view
	neutralKeyMap   NeutralKeyMap
	selectingKeyMap SelectingKeyMap
	writingKeyMap   WritingKeyMap

	Focus *FocusState

	tempIDCounter int // for generating temp IDs for new tasks
}

func New(clt *notion.Client) Model {
	f := FocusState{}

	l := list.New([]list.Item{}, NewItemDelegate(false, &f), 0, 0)
	l.Title = "Tasks"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	l.SetShowFilter(false)
	l.SetShowPagination(false)
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()

	return Model{
		// Milestone: mstone,
		list:    l,
		notion:  clt,
		loading: true,

		groups: map[string][]Item{},

		ActiveKeyMap:    NeutralKeyMapper, // default map view
		neutralKeyMap:   NeutralKeyMapper,
		selectingKeyMap: SelectingKeyMapper,
		writingKeyMap:   WritingKeyMapper,

		Focus: &f,
	}
}

func (m Model) Init() tea.Cmd {
	return m.notion.FetchTaskTypeOptions()
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
	m.list.SetItems(m.buildTaskList(notion.TaskGroups{}))

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
	m.list.SetItems(m.buildTaskList(notion.TaskGroups{}))
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
	items := m.buildTaskList(notion.TaskGroups{})
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
	m.list.SetItems(m.buildTaskList(notion.TaskGroups{}))

	// Clear pending delete state
	m.Focus.pendingDelete = false

	// todo: send command to delete task in notion

	return m
}

func (m *Model) SetItemDelegate(d list.ItemDelegate) {
	m.list.SetDelegate(d)
}

func (m Model) ClearTasks() Model {
	m.groups = map[string][]Item{}
	m.list.SetItems([]list.Item{})
	m.loading = true
	return m
}

func (m Model) buildTaskList(groups notion.TaskGroups) []list.Item {
	var items []list.Item
	for _, status := range notion.TaskStatusOrder {
		group, ok := m.groups[status]
		if !ok || len(group) == 0 {
			continue
		}
		hasMore := groups[status].NextCursor != nil
		items = append(items, GroupHeader{
			Label:   status,
			Hidden:  groups[status].Hide,
			Count:   len(group),
			HasMore: hasMore,
		})
		if !groups[status].Hide {
			for _, item := range group {
				items = append(items, item)
			}
			if groups[status].NextCursor != nil {
				items = append(items, LoadMoreItem{Status: status, Loading: groups[status].Loading})
			}
		}
	}
	return items
}
