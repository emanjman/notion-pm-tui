package note

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ItemDelegate struct {
	focused bool
}

func NewItemDelegate(focused bool) ItemDelegate {
	return ItemDelegate{
		focused: focused,
	}
}

func (d ItemDelegate) Height() int                               { return 1 }
func (d ItemDelegate) Spacing() int                              { return 0 }
func (d ItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	// todo: implement hover styling
	switch item := item.(type) {
	case Item:
		fmt.Fprint(w, item.Title)
	default:
		fmt.Fprint(w, "Unhandled")
	}
}
