package explore

import (
	"log"
	"notion-project-tui/notion"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) onQueryProjectPages(msg notion.QueryProjectPagesMsg) (Model, tea.Cmd) {
	log.Printf("on query project pages") // !debug

	if msg.Err != nil {
		log.Printf("Error: %s", msg.Err)
		m.err = msg.Err
		m.loading = false
		return m, nil
	}
	// todo: resolve to add pages to list
	// m.pages = msg.Pages
	m.loading = false

	return m, nil
}
