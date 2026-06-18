package version

import (
	"notion-project-tui/notion"
	"notion-project-tui/styles"
	"strings"

	lg "github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	style := lg.NewStyle().Width(m.width).Height(m.height)

	if m.loading {
		return style.Render("Loading...")
	}

	versions := make([]string, len(m.pages))
	for i, pg := range m.pages {

		var txtStyle lg.Style
		if i == m.pageIdx {
			txtStyle = lg.NewStyle().Foreground(styles.PrimaryForeground)
		} else {
			txtStyle = lg.NewStyle().Foreground(styles.MutedForeground)
		}

		txt := notion.ExtractPlainText(pg.Properties.Title.Title)
		versions[i] = txtStyle.Render(txt)
	}

	return style.Render(strings.Join(versions[:], "   "))
}
