package milestone

import (
	"notion-project-tui/notion"

	tea "github.com/charmbracelet/bubbletea"
)

// dispatch messages to handlers
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case notion.MilestonePagesMsg:
		return m.onMilestonePages(msg)
	case notion.FetchMoreMilestonesMsg:
		return m.onFetchMoreMilestones(msg)
	case notion.FetchMoreTasksMsg:
		return m.onFetchMoreTasksByStatus(msg)
	case notion.ToggleTaskGroupMsg:
		return m.onToggleTaskGroup(msg)
	case notion.TaskQueryMsg:
		return m.onTaskQuery(msg)
	case tea.KeyMsg:
		return m.handleKey(msg)
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
	}

	// otherwise, handle from children
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}
