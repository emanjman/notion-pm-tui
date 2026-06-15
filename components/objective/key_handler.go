package objective

import (
	"notion-project-tui/components/objective/milestone"
	"notion-project-tui/components/objective/task"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	// reserve keys at child
	if m.ChildPriorityMode() {
		return m.onChild(msg)
	}
	// switch panels
	switch {
	case key.Matches(msg, m.keys.FocusVersions):
		return m.onPanelFocus(VersionPanel)
	case key.Matches(msg, m.keys.UnfocusVersions):
		if m.panel == VersionPanel {
			return m.onPanelFocus(MilestonePanel)
		}
	case key.Matches(msg, m.keys.FocusMilestones):
		return m.onPanelFocus(MilestonePanel)
	case key.Matches(msg, m.keys.FocusTasks):
		return m.onPanelFocus(TaskPanel)
	}
	// otherwise, handle keys at the respective child-level
	return m.onChild(msg)
}

func (m Model) onPanelFocus(panel Panel) (Model, tea.Cmd) {
	mfocus, tfocus := false, false

	switch panel {
	case VersionPanel:
		m.panel = VersionPanel
		// vfocus = true
	case MilestonePanel:
		m.panel = MilestonePanel
		mfocus = true
	case TaskPanel:
		m.panel = TaskPanel
		tfocus = true
	}

	md := milestone.NewItemDelegate(mfocus, m.milestone.Mode, m.milestone.Edit)
	m.milestone.SetItemDelegate(md)
	td := task.NewItemDelegate(tfocus, m.task.Focus)
	m.task.SetItemDelegate(td)

	return m, nil
}

// handle key at child-lvl
func (m Model) onChild(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.panel {
	case VersionPanel:
		m.version, cmd = m.version.Update(msg)
	case MilestonePanel:
		m.milestone, cmd = m.milestone.Update(msg)
	case TaskPanel:
		m.task, cmd = m.task.Update(msg)
	}
	return m, cmd
}

// child in mode that needs to reserve all keys, e.g. typing
func (m Model) ChildPriorityMode() bool {
	mstoneTakesPriority := *m.milestone.Mode > milestone.NeutralMode
	taskTakesPriority := m.task.Focus.Mode > task.NeutralMode
	return mstoneTakesPriority || taskTakesPriority
}
