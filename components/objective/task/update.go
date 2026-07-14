package task

import (
	"log"
	"notion-project-tui/components/objective/milestone"
	"notion-project-tui/notion"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	case UpdateTitleMsg:
		if msg.Err != nil {
			if task, ok := m.list.SelectedItem().(Item); ok {
				log.Printf("[ERROR] updating task title in Notion, reverting...")
				task.Task = m.Focus.prevTitle
				m.list.SetItem(m.Focus.taskIdx, task)
				m.updateTaskInGroups(task)
			}
			return m, nil
		}
	case UpdateSelectionsMsg:
		if msg.Err != nil {
			if task, ok := m.list.SelectedItem().(Item); ok {
				task.Type = m.Focus.prevType
				task.Priority = m.Focus.prevPriority
				m.list.SetItem(m.Focus.taskIdx, task)
				m.updateTaskInGroups(task)
			}
		}
		return m, nil

	case UpdateStatusMsg:
		return m.onUpdateStatus(msg)

	case notion.AddTaskPageMsg:
		return m.onAddTaskPage(msg)

	case DeleteTaskMsg:
		return m.onDeleteTask(msg)

	case notion.TaskTypeOptionsMsg:
		if msg.Err != nil {
			log.Printf("[ERROR] fetching task type options: %v", msg.Err)
			return m, nil
		}
		m.typeOptions = msg.Options
		m.loading = false
		// log.Printf("[INFO] task type options loaded: %+v", m.typeOptions)
		return m, nil

	case milestone.MilestoneTasksMsg:
		// rebuild working copy from the milestone's TaskGroups on each milestone switch
		m.milestoneID = msg.MilestoneID
		m.groups = map[notion.TaskStatus][]Item{}
		for status, group := range msg.Groups {
			items := make([]Item, len(group.Tasks))
			for i, page := range group.Tasks {
				items[i] = NewItem(page)
			}
			m.groups[status] = items
		}
		m.list.SetItems(m.buildTaskList(msg.Groups))
		m.loading = false
		return m, nil

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	// forward to children
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) handleKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch m.Focus.Mode {
	case EditMode:
		return m.onEditKey(msg)
	case SelectMode:
		return m.onSelectKey(msg)
	case NormalMode:
		return m.onNormalKey(msg)
	}
	return m, nil
}

// if the notion status update failed, move the task back to its prior group
func (m Model) onUpdateStatus(msg UpdateStatusMsg) (Model, tea.Cmd) {
	if msg.Err == nil {
		return m, nil
	}
	log.Printf("[ERROR] update task status failed, reverting: %v", msg.Err)

	// find the task by id across groups; move it back to PrevStatus
	for status, group := range m.groups {
		for i, t := range group {
			if t.ID == msg.TaskID {
				m.groups[status] = append(group[:i], group[i+1:]...)
				t.Status = msg.PrevStatus
				m.groups[msg.PrevStatus] = append(m.groups[msg.PrevStatus], t)
				m.list.SetItems(m.buildTaskList(notion.TaskGroups{}))
				return m, nil
			}
		}
	}
	return m, nil
}

// if the notion trash failed, restore the optimistically-deleted task
func (m Model) onDeleteTask(msg DeleteTaskMsg) (Model, tea.Cmd) {
	if msg.Err != nil {
		log.Printf("[ERROR] delete task failed, restoring: %v", msg.Err)
		m.groups[msg.Task.Status] = append(m.groups[msg.Task.Status], msg.Task)
		m.list.SetItems(m.buildTaskList(notion.TaskGroups{}))
	}
	return m, nil
}

// reconcile the optimistic task-creation w/ result of actual notion-page creation
func (m Model) onAddTaskPage(msg notion.AddTaskPageMsg) (Model, tea.Cmd) {
	if msg.Err != nil {
		log.Printf("[ERROR] add task failed: %v", msg.Err)
		// drop the optimistic item that never made it to notion
		// (still has the temp id, so deleteTask won't hit notion)
		for _, group := range m.groups {
			for _, t := range group {
				if t.ID == msg.TempID {
					return m.deleteTask(t)
				}
			}
		}
		return m, nil
	}

	// swap the temp id for the real notion page id in groups + list
	for status, group := range m.groups {
		for i, t := range group {
			if t.ID == msg.TempID {
				t.ID = msg.Page.ID
				m.groups[status][i] = t
				m.list.SetItems(m.buildTaskList(notion.TaskGroups{}))

				// keep focus tracking pointed at the real id
				if m.Focus.taskID == msg.TempID {
					m.Focus.taskID = msg.Page.ID
				}
				return m, nil
			}
		}
	}
	return m, nil
}
