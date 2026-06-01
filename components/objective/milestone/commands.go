package milestone

import (
	"notion-project-tui/notion"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// init kickoff to get milestones; queried by milestone status
func fetchMilestonesByStatus(projID string, nc *notion.Client) tea.Cmd {
	// todo: eventually depend on versionID (over projID)
	return tea.Batch(
		nc.QueryMilestones(projID, notion.MilestoneUnderDevelopment, ""),
		nc.QueryMilestones(projID, notion.MilestoneIdle, ""),
		nc.QueryMilestones(projID, notion.MilestoneComplete, ""),
	)
}

// fetch tasks in milestones where status:under-development
func fetchInitMilestoneTasks(milestones *list.Model, nc *notion.Client) tea.Cmd {
	cmds := []tea.Cmd{}
	for i, mstone := range milestones.Items() {
		if mstone, ok := mstone.(DefaultItem); ok && mstone.MilestoneStatus == notion.MilestoneUnderDevelopment {
			mstone.FetchStatus = FetchPending
			milestones.SetItem(i, mstone)
			cmds = append(cmds, fetchInitTasks(mstone.ID, i, nc))
		}
	}
	return tea.Batch(cmds...)
}

// fetch init batch of tasks across all statuses of a milestone
func fetchInitTasks(milestoneID string, idx int, nc *notion.Client) tea.Cmd {
	// todo: i sus we update this later to use enums
	return tea.Batch(
		nc.QueryTasks(milestoneID, "dev", "", idx),
		nc.QueryTasks(milestoneID, "idle", "", idx),
		nc.QueryTasks(milestoneID, "done", "", idx),
		nc.QueryTasks(milestoneID, "archive", "", idx),
	)
}

// hit handler; refresh task panel w/ latest milestone groups
func refreshMilestoneTasks(g notion.TaskGroups) tea.Cmd {
	return func() tea.Msg {
		return MilestoneTasksMsg{Groups: g}
	}
}

// hit handler; fetch more milestones for passed milestone status
func loadMoreMilestones(s notion.MilestoneStatus) tea.Cmd {
	return func() tea.Msg {
		return notion.FetchMoreMilestonesMsg{Status: s}
	}
}

// send req to update miletone title prop on notion-server
func updateNotionMilestoneTitle(nc *notion.Client, milestoneID, title string) tea.Cmd {
	return func() tea.Msg {
		newTitle := notion.TitleProperty{Title: []notion.RichText{
			{Text: notion.TextContent{Content: title}},
		}}
		err := nc.UpdatePageProperties(milestoneID, map[string]any{"name": newTitle})
		return UpdateNotionTitleMsg{Err: err}
	}
}
