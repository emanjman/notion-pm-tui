package milestone

import (
	"fmt"
	lg "github.com/charmbracelet/lipgloss"
)

func renderItemHeader(d ItemDelegate, item GroupHeaderItem, selected bool, windowWidth int) string {
	style := d.style.header.base
	if selected {
		style = d.style.header.selected
	}

	chevron := "▼"
	if item.Hidden {
		chevron = "▶"
	}

	count := fmt.Sprintf("%d", item.Count)
	if item.HasMore {
		count += "+"
	}
	content := fmt.Sprintf("%s %s (%s)", chevron, item.Status.String(), count)
	label := style.Width(windowWidth).Render(content)
	spacer := lg.NewStyle().Render("")
	return label + "\n" + spacer
}
