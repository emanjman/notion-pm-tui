package version

import (
	"log"
	"notion-project-tui/notion"

	tea "github.com/charmbracelet/bubbletea"
)

func fetchVersions(projID string, ntn *notion.Client) tea.Cmd {
	return ntn.QueryVersionPages(projID, "")
}

// init kickoff to get milestones; queried by milestone status
func fetchInitVersionMilestones(versionID string, ntn *notion.Client) tea.Cmd {
	log.Printf("kick off from version") // !debug
	return tea.Batch(
		ntn.QueryMilestonePages(versionID, notion.MilestoneUnderDevelopment, "", notion.VersionChange),
		ntn.QueryMilestonePages(versionID, notion.MilestoneIdle, "", notion.VersionChange),
		ntn.QueryMilestonePages(versionID, notion.MilestoneComplete, "", notion.VersionChange),
	)
}
