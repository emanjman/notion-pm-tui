package task

import (
	"notion-project-tui/notion"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) onWriteKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.writingKeyMap.Save):
		var cmd tea.Cmd

		if task, ok := m.list.SelectedItem().(Item); ok {
			isNew := strings.HasPrefix(task.ID, "temp")
			title := strings.TrimSpace(m.Focus.tempTitle.Value())

			// optimistically update local task title
			m.Focus.prevTitle = task.Task
			task.Task = m.Focus.tempTitle.Value()
			m.list.SetItem(m.Focus.taskIdx, task)
			m.updateTaskInGroups(task)

			switch {
			case isNew && title != "":
				// create on notion; temp id gets reconciled on the response
				cmd = m.notion.AddTaskPage(task.ID, title, m.milestoneID, task.Type, task.Status, task.Priority)
			case isNew:
				// discard an empty brand-new task instead of persisting it
				// (temp task was never on notion, so cmd is nil)
				m, _ = m.deleteTask(task)
			default:
				// existing task: push the title update
				newTitle := notion.TitleProperty{Title: []notion.RichText{
					{Text: notion.TextContent{Content: task.Task}},
				}}
				taskID := task.ID
				cmd = func() tea.Msg {
					err := m.notion.UpdatePageProperties(
						taskID,
						map[string]any{"task": newTitle},
					)
					return UpdateTitleMsg{Err: err}
				}
			}
		}

		m.ActiveKeyMap = NormalKeyMapper
		m.Focus.Mode = NormalMode

		return m, cmd

	default:
		var cmd tea.Cmd
		m.Focus.tempTitle, cmd = m.Focus.tempTitle.Update(msg)
		return m, cmd
	}
}
