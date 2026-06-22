package version

import (
	"log"
	"notion-project-tui/notion"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) onQueryVersionPages(msg notion.QueryVersionPagesMsg) (Model, tea.Cmd) {
	log.Printf("on query version pages") // !debug

	if msg.Err != nil {
		log.Printf("Error: %s", msg.Err)
		m.err = msg.Err
		m.loading = false
		return m, nil
	}
	m.pages = msg.Pages
	m.loading = false

	// kickoff req to get milestones of this version (should get caught by milestone panel)
	return m, m.FetchInitVersionMilestones()
}
