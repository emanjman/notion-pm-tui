package projectnotelist

import (
	"fmt"
	"io"
	"notion-project-tui/styles"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type NoteListDelegate struct {
	focused bool
}

func NewNoteListDelegate(focused bool) NoteListDelegate {
	return NoteListDelegate{
		focused: focused,
	}
}

func (d NoteListDelegate) Height() int  { return 2 }
func (d NoteListDelegate) Spacing() int { return 0 }
func (d NoteListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d NoteListDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	note, ok := item.(NoteListItem)
	if !ok {
		return
	}

	selected := index == m.Index() && d.focused

	contStyle := lg.NewStyle().
		PaddingLeft(3).
		PaddingRight(1).
		Border(lg.NormalBorder(), false, false, true, false).
		BorderForeground(styles.BorderForeground)

	titleStyle := lg.NewStyle().Foreground(styles.PrimaryForeground)
	dateStyle := lg.NewStyle().Foreground(styles.MutedForeground)

	if selected {
		contStyle = contStyle.Background(styles.SelectedBackground)
		titleStyle = titleStyle.Background(styles.SelectedBackground)
		dateStyle = dateStyle.Background(styles.SelectedBackground)
	}

	// truncate title if too wide
	maxWidth := m.Width() - contStyle.GetHorizontalPadding() - contStyle.GetHorizontalBorderSize()
	title := note.NoteTitle
	if lg.Width(title) > maxWidth && maxWidth > 3 {
		for lg.Width(title+"...") > maxWidth && len(title) > 0 {
			title = title[:len(title)-1]
		}
		title = title + "..."
	}

	r1 := titleStyle.Render(title)
	r2 := dateStyle.Render(note.CreatedTime.Format("Jan 2, 2006"))

	fmt.Fprint(w, contStyle.Width(m.Width()).Render(r1+"\n"+r2))
}
