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
				m.list.SetItem(m.Focus.taskIdx, task)
				m.updateTaskInGroups(task)
			}
		}
		return m, nil

	case notion.TaskTypeOptionsMsg:
		if msg.Err != nil {
			log.Printf("[ERROR] fetching task type options: %v", msg.Err)
			return m, nil
		}
		m.typeOptions = msg.Options
		m.loading = false
		log.Printf("[INFO] task type options loaded: %+v", m.typeOptions)
		return m, nil

	case milestone.MilestoneTasksMsg:
		// rebuild working copy from the milestone's TaskGroups on each milestone switch
		m.groups = map[string][]Item{}
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
