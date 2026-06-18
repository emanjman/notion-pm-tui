package version

import (
	"notion-project-tui/notion"

	tea "github.com/charmbracelet/bubbletea"
)

func fetchVersions(projID string, ntn *notion.Client) tea.Cmd {
	return ntn.QueryVersionPages(projID, "")
}

// init kickoff to get milestones; queried by milestone status
func fetchInitVersionMilestones(versionID string, ntn *notion.Client) tea.Cmd {
	return tea.Batch(
		ntn.QueryMilestonePages(versionID, notion.MilestoneUnderDevelopment, ""),
		ntn.QueryMilestonePages(versionID, notion.MilestoneIdle, ""),
		ntn.QueryMilestonePages(versionID, notion.MilestoneComplete, ""),
	)
}
