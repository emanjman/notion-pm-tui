package explore

import (
	lg "github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	style := lg.NewStyle()

	if m.loading {
		return style.Render("Loading...")
	}

	switch m.focus {
	case ProjectFocus:
		return m.project.View()
	case SelectFocus:
		return m.list.View()
	default:
		return "Unhandled view"
	}
}
