package explore

import (
	"notion-project-tui/styles"
)

// renders proj title string; highlight fg if selected
func renderItem(d ItemDelegate, proj DefaultItem, selected bool) string {
	if selected {
		d.style = d.style.Foreground(styles.PrimaryForeground)
	}

	return d.style.Render(proj.Title)
}
