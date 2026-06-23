package milestone

import (
	"notion-project-tui/notion"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

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
func fetchInitTasks(milestoneID string, idx int, ntn *notion.Client) tea.Cmd {
	// todo: i sus we update this later to use enums
	return tea.Batch(
		ntn.QueryTasks(milestoneID, "dev", "", idx),
		ntn.QueryTasks(milestoneID, "idle", "", idx),
		ntn.QueryTasks(milestoneID, "done", "", idx),
		ntn.QueryTasks(milestoneID, "archive", "", idx),
	)
}

// hit handler; refresh task panel w/ latest milestone groups
func refreshMilestoneTasks(milestoneID string, g notion.TaskGroups) tea.Cmd {
	return func() tea.Msg {
		return MilestoneTasksMsg{MilestoneID: milestoneID, Groups: g}
	}
}

// hit handler; fetch more milestones for passed milestone status
func emitQueryMoreMilestonePages(s notion.MilestoneStatus) tea.Cmd {
	return func() tea.Msg {
		return notion.QueryMoreMilestonePagesMsg{Status: s}
	}
}

// send req to update miletone title prop on notion-server
func updateNotionMilestoneTitle(ntn *notion.Client, milestoneID, title string) tea.Cmd {
	return func() tea.Msg {
		newTitle := notion.TitleProperty{Title: []notion.RichText{
			{Text: notion.TextContent{Content: title}},
		}}
		err := ntn.UpdatePageProperties(milestoneID, map[string]any{"name": newTitle})
		return UpdateNotionTitleMsg{Err: err}
	}
}

// hit handler; reconcile trash page req
func emitTrashMilestonePage(err error) tea.Cmd {
	return func() tea.Msg {
		return TrashMilestonePageMsg{Err: err}
	}
}
