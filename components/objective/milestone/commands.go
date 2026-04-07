package milestone

import (
	"notion-project-tui/notion"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// tea.Cmd factory + msg types
// about sending commands out (w/ no context of the model)

func fetchMilestoneTasks(milestones *list.Model, nc *notion.Client) tea.Cmd {
	cmds := []tea.Cmd{}

	for i, m := range milestones.Items() {
		if m, ok := m.(DefaultItem); ok && m.MilestoneStatus == notion.MilestoneUnderDevelopment {
			m.FetchStatus = FetchPending
			milestones.SetItem(i, m)
			cmds = append(cmds, fetchTasksByStatus(m.ID, i, nc))
		}
	}

	return tea.Batch(cmds...)
}

func fetchTasksByStatus(milestoneID string, idx int, nc *notion.Client) tea.Cmd {
	return tea.Batch(
		nc.QueryTasks(milestoneID, "dev", "", idx),
		nc.QueryTasks(milestoneID, "idle", "", idx),
		nc.QueryTasks(milestoneID, "done", "", idx),
		nc.QueryTasks(milestoneID, "archive", "", idx),
	)
}

func emitTaskViewMsg(g notion.TaskGroups) tea.Cmd {
	return func() tea.Msg {
		return TaskViewMsg{Groups: g}
	}
}
