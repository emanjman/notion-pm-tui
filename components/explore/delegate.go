package explore

import (
	"fmt"
	"io"
	"notion-project-tui/styles"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type ItemDelegate struct {
	style lg.Style
}

var _ list.ItemDelegate = (*ItemDelegate)(nil) // compile-time compliance

func NewItemDelegate() ItemDelegate {
	// setup base style as muted color
	style := lg.NewStyle().Foreground(styles.MutedForeground)
	return ItemDelegate{
		style: style,
	}
}

func (d ItemDelegate) Render(w io.Writer, m list.Model, idx int, item list.Item) {
	selected := idx == m.Index()

	switch item := item.(type) {
	case DefaultItem:
		fmt.Fprint(w, renderItem(d, item, selected))
	}
}
func (d ItemDelegate) Height() int                               { return 2 }
func (d ItemDelegate) Spacing() int                              { return 0 }
func (d ItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
