package objective

import (
	"notion-project-tui/components/objective/milestone"
	"notion-project-tui/components/objective/task"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	if m.InFocusMode() {
		var cmd tea.Cmd
		// delegate keys to focused child
		switch m.focus {
		case MilestonePanel:
			m.milestone, cmd = m.milestone.Update(msg)
		case TaskPanel:
			m.task, cmd = m.task.Update(msg)
		}
		return m, cmd
	} else {
		switch {
		case key.Matches(msg, m.keys.LeftFocus):
			return m.onLeftFocus(msg)
		case key.Matches(msg, m.keys.RightFocus):
			return m.onRightFocus(msg)
		}
	}
	return m, nil
}

func (m Model) onLeftFocus(msg tea.KeyMsg) (Model, tea.Cmd) {
	m.focus = MilestonePanel
	m.task.SetItemDelegate(task.NewItemDelegate(false, m.task.Focus))
	m.milestone.SetItemDelegate(milestone.NewItemDelegate(true, m.milestone.Mode, m.milestone.Edit))
	return m, nil
}

func (m Model) onRightFocus(msg tea.KeyMsg) (Model, tea.Cmd) {
	m.focus = TaskPanel

	m.task.SetItemDelegate(task.NewItemDelegate(true, m.task.Focus))
	m.milestone.SetItemDelegate(milestone.NewItemDelegate(false, m.milestone.Mode, m.milestone.Edit))

	return m, nil
}
