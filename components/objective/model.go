package objective

import (
	"notion-project-tui/components/milestonelist"
	"notion-project-tui/components/tasklist"
	"notion-project-tui/notion"
	"notion-project-tui/styles"
	"notion-project-tui/util/keymap"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type Panel int

const (
	MilestonesPanel Panel = iota
	TasksPanel
)

type ObjectiveModel struct {
	focus      Panel
	milestones milestonelist.MilestoneListModel
	tasks      tasklist.TaskListModel
	keys       KeyMap
}

func NewObjectiveModel(c *notion.Client) ObjectiveModel {
	milestones := milestonelist.NewMilestoneListModel()
	tasks := tasklist.NewTaskListModel(milestones.SelectedMilestone(), c)

	return ObjectiveModel{
		focus:      MilestonesPanel,
		milestones: milestones,
		tasks:      tasks,
		keys:       DefaultKeyMap,
	}
}

func (m ObjectiveModel) Init() tea.Cmd {
	return nil
}

func (m ObjectiveModel) Update(msg tea.Msg) (ObjectiveModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		if m.tasks.Focus.Mode != tasklist.NeutralMode {
			var cmd tea.Cmd
			m.tasks, cmd = m.tasks.Update(msg)
			return m, cmd
		} else {
			switch {
			case key.Matches(msg, m.keys.LeftFocus):
				m.tasks.Milestone = m.milestones.SelectedMilestone()
				m.focus = MilestonesPanel

				m.tasks.SetItemDelegate(tasklist.NewTaskListDelegate(false, m.tasks.Focus))
				m.milestones.SetItemDelegate(milestonelist.NewMilestoneListDelegate(true, m.milestones.Focus))

				return m, nil

			case key.Matches(msg, m.keys.RightFocus):
				m.focus = TasksPanel

				m.tasks.SetItemDelegate(tasklist.NewTaskListDelegate(true, m.tasks.Focus))
				m.milestones.SetItemDelegate(milestonelist.NewMilestoneListDelegate(false, m.milestones.Focus))

				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		leftWidth := msg.Width * 30 / 100
		rightWidth := msg.Width - leftWidth - 1 // account for dividing border

		m.milestones, cmd = m.milestones.Update(tea.WindowSizeMsg{
			Width:  leftWidth,
			Height: msg.Height,
		})
		m.tasks, cmd = m.tasks.Update(tea.WindowSizeMsg{
			Width:  rightWidth,
			Height: msg.Height,
		})

		return m, cmd
	}

	// key presses go to active panel only; data messages go to both
	var milestoneCmd, taskCmd tea.Cmd

	if _, isKey := msg.(tea.KeyMsg); isKey {
		switch m.focus {
		case MilestonesPanel:
			m.milestones, milestoneCmd = m.milestones.Update(msg)
			return m, milestoneCmd
		case TasksPanel:
			m.tasks, taskCmd = m.tasks.Update(msg)
			return m, taskCmd
		}
	} else {
		m.milestones, milestoneCmd = m.milestones.Update(msg)
		m.tasks, taskCmd = m.tasks.Update(msg)
	}

	return m, tea.Batch(milestoneCmd, taskCmd)

}

func (m ObjectiveModel) View() string {
	left := lg.NewStyle().
		BorderRight(true).
		BorderStyle(lg.NormalBorder()).
		BorderForeground(styles.BorderForeground).
		Render(m.milestones.View())
	right := m.tasks.View()
	return lg.JoinHorizontal(lg.Top, left, right)
}

func (m ObjectiveModel) KeyMap() help.KeyMap {
	switch m.focus {
	case MilestonesPanel:
		return keymap.JoinedKeyMap{Primary: m.keys, Secondary: m.milestones.ActiveKeyMap}
	case TasksPanel:
		return keymap.JoinedKeyMap{Primary: m.keys, Secondary: m.tasks.ActiveKeyMap}
	}
	return nil
}
