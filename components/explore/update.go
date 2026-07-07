package explore

import (
	"notion-project-tui/components/project"
	"notion-project-tui/notion"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// reserve msgs to active project
	if m.project != nil {
		proj, cmd := m.project.Update(msg)
		m.project = &proj
		return m, cmd
	}

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
				proj := project.New(item.ID)
				cmd := proj.Init()
				m.project = &proj

				return m, cmd
			}
		default:
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		}
	case notion.QueryProjectPagesMsg:
		return m.onQueryProjectPages(msg)
	}

	return m, nil
}
