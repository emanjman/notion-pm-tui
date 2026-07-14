package task

import (
	"notion-project-tui/notion"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) onEditKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.editKeyMap.Save):
		var cmd tea.Cmd

		if task, ok := m.list.SelectedItem().(Item); ok {
			isNew := strings.HasPrefix(task.ID, "temp")
			title := strings.TrimSpace(m.EditCtx.titleInput.Value())

			// optimistically update local task title
			m.EditCtx.titleBackup = task.Task
			task.Task = title
			m.list.SetItem(m.EditCtx.taskIdx, task)
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
		*m.Mode = NormalMode

		return m, cmd

	default:
		var cmd tea.Cmd
		m.EditCtx.titleInput, cmd = m.EditCtx.titleInput.Update(msg)
		return m, cmd
	}
}
