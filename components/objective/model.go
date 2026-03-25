package objective

import (
	"notion-project-tui/components/objective/milestone"
	"notion-project-tui/components/objective/task"
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
	MilestonePanel Panel = iota
	TaskPanel
)

type Model struct {
	projID           string
	milestonesPropID string
	loading          bool
	err              error
	focus            Panel
	notion           *notion.Client
	milestone        milestone.Model
	task             task.Model
	keys             KeyMap
}

func New(n *notion.Client, projID, milestonesPropID string) Model {
	ms := milestone.New(n, projID, milestonesPropID)
	t := task.New(n)

	return Model{
		projID:           projID,
		milestonesPropID: milestonesPropID,
		loading:          true,
		err:              nil,
		focus:            MilestonePanel,
		notion:           n,
		milestone:        ms,
		task:             t,
		keys:             DefaultKeyMap,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.milestone.Init(), m.task.Init())
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.InFocusMode() {
			var cmd tea.Cmd
			if m.focus == MilestonePanel {
				m.milestone, cmd = m.milestone.Update(msg)
			} else if m.focus == TaskPanel {
				m.task, cmd = m.task.Update(msg)
			}
			return m, cmd
		} else {
			switch {
			case key.Matches(msg, m.keys.LeftFocus):
				// m.task.Milestone = m.milestone.SelectedMilestone()
				m.focus = MilestonePanel

				m.task.SetItemDelegate(task.NewItemDelegate(false, m.task.Focus))
				m.milestone.SetItemDelegate(milestone.NewItemDelegate(true, m.milestone.Focus))

				return m, nil

			case key.Matches(msg, m.keys.RightFocus):
				m.focus = TaskPanel

				m.task.SetItemDelegate(task.NewItemDelegate(true, m.task.Focus))
				m.milestone.SetItemDelegate(milestone.NewItemDelegate(false, m.milestone.Focus))

				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		var mstoneCmd, taskCmd tea.Cmd
		leftWidth := msg.Width * 30 / 100
		rightWidth := msg.Width - leftWidth - 1 // account for dividing border

		m.milestone, mstoneCmd = m.milestone.Update(tea.WindowSizeMsg{
			Width:  leftWidth,
			Height: msg.Height,
		})
		m.task, taskCmd = m.task.Update(tea.WindowSizeMsg{
			Width:  rightWidth,
			Height: msg.Height,
		})

		return m, tea.Batch(mstoneCmd, taskCmd)
	}

	// key presses go to active panel only; data messages go to both
	var milestoneCmd, taskCmd tea.Cmd

	if _, isKey := msg.(tea.KeyMsg); isKey {
		switch m.focus {
		case MilestonePanel:
			m.milestone, milestoneCmd = m.milestone.Update(msg)
			return m, milestoneCmd
		case TaskPanel:
			m.task, taskCmd = m.task.Update(msg)
			return m, taskCmd
		}
	} else {
		m.milestone, milestoneCmd = m.milestone.Update(msg)
		m.task, taskCmd = m.task.Update(msg)
	}

	return m, tea.Batch(milestoneCmd, taskCmd)

}

func (m Model) View() string {
	left := lg.NewStyle().
		BorderRight(true).
		BorderStyle(lg.NormalBorder()).
		BorderForeground(styles.BorderForeground).
		Render(m.milestone.View())
	right := m.task.View()
	return lg.JoinHorizontal(lg.Top, left, right)
}

func (m Model) KeyMap() help.KeyMap {
	switch m.focus {
	case MilestonePanel:
		return keymap.JoinedKeyMap{Primary: m.keys, Secondary: m.milestone.ActiveKeyMap}
	case TaskPanel:
		return keymap.JoinedKeyMap{Primary: m.keys, Secondary: m.task.ActiveKeyMap}
	}
	return nil
}

func (m Model) InFocusMode() bool {
	return m.milestone.Focus.Mode > 0 || m.task.Focus.Mode > 0
}
