package explore

import (
	"log"
	"notion-project-tui/notion"

	"github.com/charmbracelet/bubbles/list"
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

	var items []list.Item
	for _, pg := range msg.Pages {
		items = append(items, NewDefaultItem(pg))
	}
	m.list.SetItems(items)
	m.loading = false

	return m, nil
}
