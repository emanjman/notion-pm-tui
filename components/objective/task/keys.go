package task

import (
	"notion-project-tui/notion"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type NeutralKeyMap struct {
	Up         key.Binding
	Down       key.Binding
	JumpUp     key.Binding // jump up 5
	JumpDown   key.Binding // jump down 5
	Select     key.Binding // enter focus (select) mode
	StatusPrev key.Binding // cycle status backward
	StatusNext key.Binding // cycle status forward
	AddTask    key.Binding // add new task to idle group
	Delete     key.Binding // delete task (requires confirmation)
}

var NeutralKeyMapper = NeutralKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	JumpUp: key.NewBinding(
		key.WithKeys("ctrl+k"),
		key.WithHelp("ctrl+k", "jump up 5"),
	),
	JumpDown: key.NewBinding(
		key.WithKeys("ctrl+j"),
		key.WithHelp("ctrl+j", "jump down 5"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select task"),
	),
	StatusPrev: key.NewBinding(
		key.WithKeys("<", "shift+,"),
		key.WithHelp("<", "prev status"),
	),
	StatusNext: key.NewBinding(
		key.WithKeys(">", "shift+."),
		key.WithHelp(">", "next status"),
	),
	AddTask: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add task"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete task"),
	),
}

func (k NeutralKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select}
}

func (k NeutralKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select},
		{k.StatusPrev, k.StatusNext, k.AddTask, k.Delete},
	}
}

// ---

type SelectingKeyMap struct {
	Left   key.Binding // prev field
	Right  key.Binding // next field
	Select key.Binding // cycle select-options or enter rewrite mode
	Exit   key.Binding // send off changes to notion (server)
}

var SelectingKeyMapper = SelectingKeyMap{
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("<-/h", "prev field"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("->/l", "right field"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "enter edit mode"),
	),
	Exit: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "save + exit"),
	),
}

func (k SelectingKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Left, k.Right, k.Select}
}

func (k SelectingKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Left, k.Right, k.Select},
		{k.Exit},
	}
}

// ---

type WritingKeyMap struct {
	Save key.Binding // update list item (client)
}

var WritingKeyMapper = WritingKeyMap{
	Save: key.NewBinding(
		key.WithKeys("enter", "esc"),
		key.WithHelp("enter/esc", "save changes"),
	),
}

func (k WritingKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Save}
}

func (k WritingKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Save},
		{},
	}
}

// ---

func (m Model) handleKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch m.Focus.Mode {
	case WritingMode:
		return m.onWritingKey(msg)
	case SelectingMode:
		return m.onSelectingKey(msg)
	case NeutralMode:
		return m.onNeutralKey(msg)
	}
	return m, nil
}

func (m Model) onWritingKey(msg tea.KeyMsg) (Model, tea.Cmd) {
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
				cmd = m.notion.AddTask(task.ID, title, m.milestoneID, task.Status, task.Type, task.Priority)
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

		m.ActiveKeyMap = NeutralKeyMapper
		m.Focus.Mode = NeutralMode

		return m, cmd

	default:
		var cmd tea.Cmd
		m.Focus.tempTitle, cmd = m.Focus.tempTitle.Update(msg)
		return m, cmd
	}
}

func (m Model) onSelectingKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.selectingKeyMap.Exit):
		m.Focus.Mode = NeutralMode
		m.ActiveKeyMap = NeutralKeyMapper

		if task, ok := m.list.SelectedItem().(Item); ok {
			var typeOpt notion.SelectItem
			for _, opt := range m.typeOptions {
				if opt.Name == task.Type {
					typeOpt = opt
					break
				}
			}
			taskID := m.Focus.taskID
			return m, func() tea.Msg {
				err := m.notion.UpdatePageProperties(taskID, map[string]any{
					"type": map[string]any{"select": typeOpt},
				})
				return UpdateSelectionsMsg{Err: err}
			}
		}
		return m, nil

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

	case key.Matches(msg, m.selectingKeyMap.Select):
		if task, ok := m.list.SelectedItem().(Item); ok {
			switch m.Focus.field {
			case TaskType:
				m.Focus.prevType = task.Type
				task.Type = cycleTypeField(task.Type, 1, m.typeOptions)
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

	// consume all keys, don't forward to list navigation
	return m, nil
}

func (m Model) onNeutralKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	// cancel pending delete if any key other than Delete is pressed
	if m.Focus.pendingDelete && !key.Matches(msg, m.neutralKeyMap.Delete) {
		m.Focus.pendingDelete = false
		return m, nil
	}

	switch {
	case key.Matches(msg, m.neutralKeyMap.Select):
		selected := m.list.SelectedItem()
		if header, ok := selected.(GroupHeader); ok {
			return m, func() tea.Msg {
				return notion.ToggleTaskGroupMsg{Status: header.Label}
			}
		} else if loadMore, ok := selected.(LoadMoreItem); ok && !loadMore.Loading {
			return m, func() tea.Msg {
				return notion.QueryMoreTaskPagesMsg{Status: loadMore.Status}
			}
		} else if task, ok := selected.(Item); ok {
			m.Focus.taskID = task.ID
			m.Focus.taskIdx = m.list.Index()
			m.Focus.field = TaskTitle

			m.ActiveKeyMap = SelectingKeyMapper
			m.Focus.Mode = SelectingMode
		}
		return m, nil

	case key.Matches(msg, m.neutralKeyMap.AddTask):
		m = m.addTask()
		return m, nil

	case key.Matches(msg, m.neutralKeyMap.JumpUp):
		m.list.Select(max(0, m.list.Index()-5))
		return m, nil

	case key.Matches(msg, m.neutralKeyMap.JumpDown):
		m.list.Select(min(len(m.list.Items())-1, m.list.Index()+5))
		return m, nil

	case key.Matches(msg, m.neutralKeyMap.StatusPrev):
		if task, ok := m.list.SelectedItem().(Item); ok {
			m = m.changeTaskStatus(task, -1)
		}
		return m, nil

	case key.Matches(msg, m.neutralKeyMap.StatusNext):
		if task, ok := m.list.SelectedItem().(Item); ok {
			m = m.changeTaskStatus(task, +1)
		}
		return m, nil

	case key.Matches(msg, m.neutralKeyMap.Delete):
		if task, ok := m.list.SelectedItem().(Item); ok {
			if m.Focus.pendingDelete && m.Focus.taskID == task.ID {
				return m.deleteTask(task)
			}
			m.Focus.pendingDelete = true
			m.Focus.taskID = task.ID
			m.Focus.taskIdx = m.list.Index()
		}
		return m, nil
	}

	// forward unmatched keys to list (e.g. up/down navigation)
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}
