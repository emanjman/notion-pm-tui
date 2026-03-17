package projectnotelist

import (
	"notion-project-tui/notion"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type ProjectNotesListModel struct {
	list    list.Model
	loading bool
	error   error
	focused bool
}

func NewProjectNotesListModel() ProjectNotesListModel {
	l := list.New([]list.Item{}, NewNoteListDelegate(true), 0, 0)
	l.Title = "Notes"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	l.DisableQuitKeybindings()

	return ProjectNotesListModel{
		list:    l,
		loading: true,
		focused: true,
	}
}

func (m ProjectNotesListModel) Init() tea.Cmd {
	return nil
}

func (m ProjectNotesListModel) Update(msg tea.Msg) (ProjectNotesListModel, tea.Cmd) {
	prevIndex := m.list.Index()

	switch msg := msg.(type) {
	case notion.ProjectNoteMsg:
		if msg.Err != nil {
			m.error = msg.Err
			return m, nil
		}

		items := make([]list.Item, len(msg.Data))
		for i, page := range msg.Data {
			items[i] = NewNoteListItem(page)
		}

		m.list.SetItems(items)
		m.loading = false

		// emit selection for the first item
		if len(items) > 0 {
			if note, ok := items[0].(NoteListItem); ok {
				return m, func() tea.Msg {
					return notion.ProjectNoteSelectedMsg{ID: note.ID}
				}
			}
		}

		return m, nil

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	// emit selection msg if cursor moved
	newIndex := m.list.Index()
	if newIndex != prevIndex {
		if note, ok := m.list.SelectedItem().(NoteListItem); ok {
			return m, tea.Batch(cmd, func() tea.Msg {
				return notion.ProjectNoteSelectedMsg{ID: note.ID}
			})
		}
	}

	return m, cmd
}

func (m ProjectNotesListModel) View() string {
	containerStyle := lg.NewStyle().PaddingRight(1)
	return containerStyle.Render(m.list.View())
}

func (m ProjectNotesListModel) SelectedNoteID() string {
	if item, ok := m.list.SelectedItem().(NoteListItem); ok {
		return item.ID
	}
	return ""
}

func (m *ProjectNotesListModel) SetFocused(focused bool) {
	m.focused = focused
	m.list.SetDelegate(NewNoteListDelegate(focused))
}
