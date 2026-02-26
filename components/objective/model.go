package objective

import (
	"notion-project-tui/components/milestonelist"
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
	// tasks      tasklist.TaskListModel
	keys KeyMap
}

func NewObjectiveModel() ObjectiveModel {
	return ObjectiveModel{
		focus:      MilestonesPanel,
		milestones: milestonelist.NewMilestoneListModel(),
		keys:       DefaultKeyMap,
	}
}

func (m ObjectiveModel) Init() tea.Cmd {
	return nil
}

func (m ObjectiveModel) Update(msg tea.Msg) (ObjectiveModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.LeftFocus):
			m.focus = MilestonesPanel
			// todo: would need to tell tasks what milestone is selected
		case key.Matches(msg, m.keys.RightFocus):
			m.focus = TasksPanel
		default:
			// forward to active panel
			var cmd tea.Cmd

			switch m.focus {
			case MilestonesPanel:
				m.milestones, cmd = m.milestones.Update(msg)

			case TasksPanel:
				// m.tasks, cmd = m.tasks.Update(msg)

				return m, cmd
			}
		}

	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		leftWidth := msg.Width * 40 / 100
		// rightWidth := msg.Width - leftWidth

		m.milestones, cmd = m.milestones.Update(tea.WindowSizeMsg{
			Width:  leftWidth,
			Height: msg.Height,
		})

		// todo: apply task width
		// m.tasks, cmd = m.tasks.Update(tea.WindowSizeMsg{
		// 	Width:  rightWidth,
		// 	Height: msg.Height,
		// })

		return m, cmd
	}

	return m, nil
}

func (m ObjectiveModel) View() string {
	left := m.milestones.View()
	// right := m.tasks.View()
	right := ""
	return lg.JoinHorizontal(lg.Top, left, right)
}

func (m ObjectiveModel) KeyMap() help.KeyMap {
	switch m.focus {
	case MilestonesPanel:
		// return m.milestones.Keys
		return keymap.JoinedKeyMap{Primary: m.keys, Secondary: m.milestones.Keys}

	case TasksPanel:
		// return SharedKeyMap{shared: m.keys, focus: m.tasks.Keys}
		// todo: use tasks keys
		return keymap.JoinedKeyMap{Primary: m.keys, Secondary: m.milestones.Keys}
	}

	return nil
}
