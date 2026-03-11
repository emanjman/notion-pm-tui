package pagecontent

import (
	"fmt"
	"io"
	"notion-project-tui/notion"
	"notion-project-tui/styles"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type PageContentDelegate struct{}

func NewPageContentDelegate() PageContentDelegate {
	return PageContentDelegate{}
}

func (d PageContentDelegate) Height() int                               { return 1 }
func (d PageContentDelegate) Spacing() int                              { return 0 }
func (d PageContentDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d PageContentDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var s string
	selected := index == m.Index()

	style := lg.NewStyle()
	if selected {
		style = style.Background(styles.SelectedBackground)
	}

	switch item := item.(type) {
	case notion.Block:
		switch item.Type {
		case notion.Divider:
			s = strings.Repeat("-", m.Width())

		case notion.Paragraph:
			s = notion.ExtractPlainText(item.Paragraph.RichText)

		case notion.Heading2:
			s = notion.ExtractPlainText(item.Heading2.RichText)

		case notion.Heading3:
			s = notion.ExtractPlainText(item.Heading3.RichText)

		default:
			s = "Unhandled"
		}
	}

	fmt.Fprint(w, style.Width(m.Width()).Render(s))

}
