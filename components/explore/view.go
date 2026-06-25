package explore

import (
	lg "github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	style := lg.NewStyle()

	if m.loading {
		return style.Render("Loading...")
	}

	if m.project != nil {
		return m.project.View()
	}
	return m.list.View()
}
