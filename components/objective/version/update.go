package version

import (
	"notion-project-tui/notion"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	case notion.QueryVersionPagesMsg:
		return m.onQueryVersionPages(msg)
	case tea.WindowSizeMsg:
		return m.handleWindow(msg)
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch *m.Mode {
	case NeutralMode:
		switch {
		case key.Matches(msg, m.neutralKeyMap.Next):
			return m.onNeutralNext()
		case key.Matches(msg, m.neutralKeyMap.Prev):
			return m.onNeutralPrev()
		}
	}
	return m, nil
}

func (m Model) handleWindow(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height
	return m, nil
}
