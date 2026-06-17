package version

import (
	"strings"

	lg "github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	style := lg.NewStyle().Width(m.width).Height(m.height)

	tempVersions := [2]string{"v1 - Project Management", "v2 - Final Fantasy"}

	return style.Render(strings.Join(tempVersions[:], "  "))
}
