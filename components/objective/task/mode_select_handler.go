package task

import (
	"notion-project-tui/notion"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) onSelectKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.selectingKeyMap.Exit):
		m.Focus.Mode = NormalMode
		m.ActiveKeyMap = NormalKeyMapper

		if task, ok := m.list.SelectedItem().(Item); ok {
			var typeOpt notion.SelectItem
			for _, opt := range m.typeOptions {
				if opt.Name == task.Type {
					typeOpt = opt
					break
				}
			}
			taskID := m.Focus.taskID
			priority := strconv.Itoa(task.Priority)
			return m, func() tea.Msg {
				err := m.notion.UpdatePageProperties(taskID, map[string]any{
					"type":     map[string]any{"select": typeOpt},
					"priority": map[string]any{"select": map[string]any{"name": priority}},
				})
				return UpdateSelectionsMsg{Err: err}
			}
		}
		return m, nil

	case key.Matches(msg, m.selectingKeyMap.Left):
		if m.Focus.field == TaskType {
			m.Focus.field = _SelectedFieldCount - 1
		} else {
			m.Focus.field = (m.Focus.field - 1) % _SelectedFieldCount
		}
		return m, nil

	case key.Matches(msg, m.selectingKeyMap.Right):
		m.Focus.field = (m.Focus.field + 1) % _SelectedFieldCount
		return m, nil

	case key.Matches(msg, m.selectingKeyMap.Select):
		if task, ok := m.list.SelectedItem().(Item); ok {
			switch m.Focus.field {
			case TaskType:
				m.Focus.prevType = task.Type
				task.Type = cycleTypeField(task.Type, 1, m.typeOptions)
			case TaskPriority:
				m.Focus.prevPriority = task.Priority
				task.Priority = cyclePriorityField(task.Priority, 1)
			case TaskTitle:
				m.Focus.Mode = WriteMode
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

	// consume all keys, don't forward to list navigation
	return m, nil
}
