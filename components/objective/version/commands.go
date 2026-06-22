package version

import (
	"notion-project-tui/notion"

	tea "github.com/charmbracelet/bubbletea"
)

func fetchVersions(projID string, ntn *notion.Client) tea.Cmd {
	return ntn.QueryVersionPages(projID, "")
}

// init kickoff to get milestones; queried by milestone status
func (m Model) FetchInitVersionMilestones() tea.Cmd {
	versionID := m.pages[m.pageIdx].ID
	ntn := m.notion
	return tea.Batch(
		ntn.QueryMilestonePages(versionID, notion.MilestoneUnderDevelopment, "", notion.VersionChange),
		ntn.QueryMilestonePages(versionID, notion.MilestoneIdle, "", notion.VersionChange),
		ntn.QueryMilestonePages(versionID, notion.MilestoneComplete, "", notion.VersionChange),
	)
}
