package project

import (
	"notion-project-tui/styles"
	"notion-project-tui/util/keymap"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	lg "github.com/charmbracelet/lipgloss"
)

var labels = []string{
	"Objective (n%)",
	"Notebook (n)",
	"Bug Report (n)",
	"Technology (n)",
}

func (m Model) View() string {
	var s strings.Builder

	lg.JoinHorizontal(lg.Top)
	headers := make([]string, len(labels))

	for i := range labels {
		base := lg.NewStyle().Padding(0, 2)

		tabStyle := base.Foreground(styles.MutedForeground)
		if int(m.activeTab) == i {
			tabStyle = base.
				Foreground(styles.PrimaryForeground).
				Background(styles.SelectedBackground)
		}
		headers[i] = tabStyle.Render(labels[i])
	}
	main := ""

	switch m.activeTab {

	case ObjectiveTab:
		main = m.objective.View()
	case NotebookTab:
		main = m.notebook.View()
	case BugsTab:
		main = "Debug notes (coming soon)"
	case TechTab:
		main = "Tech notes (coming soon)"
	}

	tabDivider := lg.NewStyle().
		Foreground(styles.BorderForeground).
		SetString("|")
	s.WriteString(
		lg.NewStyle().
			Padding(1, 2).
			Width(m.width).
			Render(strings.Join(headers, tabDivider.String())))

	s.WriteString("\n")
	s.WriteString(main)
	s.WriteString("\n")

	help := m.help.View(keymap.JoinedKeyMap{
		Primary:   RootKeyMap,
		Secondary: m.getActiveKeyMap(),
	})
	s.WriteString(
		lg.NewStyle().
			Padding(1, 2).
			Width(m.width).
			Render(help))

	return s.String()
}

func (m Model) getActiveKeyMap() help.KeyMap {
	switch m.activeTab {

	case ObjectiveTab:
		return m.objective.KeyMap()

	// todo: handle other tabs
	case NotebookTab:
		return m.notebook.ActiveKeyMap
	case BugsTab:
		return m.objective.KeyMap()
	case TechTab:
		return m.objective.KeyMap()

	default:
		// return nil
		return m.objective.KeyMap()
	}
}
