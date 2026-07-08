package explore

import (
	"notion-project-tui/notion"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// handle focus-agnostic msgs
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		// handle list model
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height)

		// handle proj model
		var cmd tea.Cmd
		m.project, cmd = m.project.Update(msg)

		return m, cmd
	}

	switch m.focus {
	case SelectFocus:
		switch msg := msg.(type) {
		case notion.QueryProjectPagesMsg:
			return m.onQueryProjectPages(msg)
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.neutralKeyMap.Quit):
				return m, tea.Quit
			case key.Matches(msg, m.neutralKeyMap.Select):
				return m.onProjectSelect(msg)
			default:
				var cmd tea.Cmd
				m.list, cmd = m.list.Update(msg)
				return m, cmd
			}
		}
	case ProjectFocus:
		// forward msgs to be handled by project model
		var cmd tea.Cmd
		m.project, cmd = m.project.Update(msg)
		return m, cmd
	}

	return m, nil
}

// send msg w/ projID to get caught by all moedls that
func (m Model) onProjectSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if item, ok := m.list.SelectedItem().(DefaultItem); ok {
		m.focus = ProjectFocus
		return m, emitProjectID(item.ID)
	}
	return m, nil
}
