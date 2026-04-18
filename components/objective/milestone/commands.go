package milestone

import (
	"notion-project-tui/notion"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// tea.Cmd factory + msg types
// about sending commands out (w/ no context of the model)

// init task fetches for all under-development milestones
func fetchInitMilestoneTasks(milestones *list.Model, nc *notion.Client) tea.Cmd {
	cmds := []tea.Cmd{}

	for i, m := range milestones.Items() {
		if m, ok := m.(DefaultItem); ok && m.MilestoneStatus == notion.MilestoneUnderDevelopment {
			m.FetchStatus = FetchPending
			milestones.SetItem(i, m)
			cmds = append(cmds, fetchInitTasks(m.ID, i, nc))
		}
	}

	return tea.Batch(cmds...)
}

// fetch initial set of tasks across all statuses for a milestone
func fetchInitTasks(milestoneID string, idx int, nc *notion.Client) tea.Cmd {
	return tea.Batch(
		nc.QueryTasks(milestoneID, "dev", "", idx),
		nc.QueryTasks(milestoneID, "idle", "", idx),
		nc.QueryTasks(milestoneID, "done", "", idx),
		nc.QueryTasks(milestoneID, "archive", "", idx),
	)
}

// refresh task panel w/ latest milestone groups
func refreshMilestoneTasks(g notion.TaskGroups) tea.Cmd {
	return func() tea.Msg {
		return MilestoneTasksMsg{Groups: g}
	}
}

func fetchMilestonesByStatus(projID string, nc *notion.Client) tea.Cmd {
	return tea.Batch(
		nc.QueryMilestones(projID, notion.MilestoneUnderDevelopment, ""),
		nc.QueryMilestones(projID, notion.MilestoneIdle, ""),
		nc.QueryMilestones(projID, notion.MilestoneComplete, ""),
	)

}
