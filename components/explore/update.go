package explore

import (
	"notion-project-tui/notion"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case notion.QueryProjectPagesMsg:
		return m.onQueryProjectPages(msg)
	}
	return m, nil
}
