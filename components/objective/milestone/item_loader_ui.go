package milestone

import (
	lg "github.com/charmbracelet/lipgloss"
	"notion-project-tui/styles"
)

func renderLoadMore(d ItemDelegate, loading bool, selected bool, noBorder bool, windowWidth int) string {
	style := d.style.itemContainer.base.Foreground(styles.MutedForeground).PaddingLeft(2)
	if selected {
		style = d.style.itemContainer.selected.Foreground(styles.MutedForeground).PaddingLeft(2)
	}
	if noBorder {
		style = style.Border(lg.NormalBorder(), false)
	}
	text := "..."
	if loading {
		text = "Loading..."
	} else if selected {
		text = "[Enter] to load more..."
	}
	rendered := style.Width(windowWidth).Render(text)
	if noBorder {
		return rendered + "\n" + lg.NewStyle().Render("")
	}
	return rendered
}
