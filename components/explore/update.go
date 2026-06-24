package explore

import (
	"log"
	"notion-project-tui/notion"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height)
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.neutralKeyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.neutralKeyMap.Select):
			if item, ok := m.list.SelectedItem().(DefaultItem); ok {
				proj := m.project.New(item.ID)
				m.project = &proj
				return m, nil
			}
		default:
			log.Printf("is catching default keys")
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		}
	case notion.QueryProjectPagesMsg:
		return m.onQueryProjectPages(msg)
	}
	return m, nil
}
