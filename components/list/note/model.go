package note

import (
	"notion-project-tui/notion"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	projID      string
	notesPropID string
	list        list.Model
	loading     bool
	err         error
	notion      *notion.Client
}

func New(notion *notion.Client, projID, notesPropID string) Model {
	// list configs
	l := list.New([]list.Item{}, NewItemDelegate(true), 0, 0)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	l.DisableQuitKeybindings()

	return Model{
		projID:      projID,
		notesPropID: notesPropID,
		list:        l,
		loading:     true,
		err:         nil,
		notion:      notion,
	}
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		ids, err := m.notion.FetchRelationIDs(m.projID, m.notesPropID)
		return notion.NoteIDsMsg{IDs: ids, Err: err}
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	// todo: on project selection msg, send out fetch-relation-ids req

	case notion.NoteIDsMsg:
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			return m, func() tea.Msg {
				pages, err := notion.FetchPages[notion.NotePage](m.notion, msg.IDs)
				return notion.NotePagesMsg{Pages: pages, Err: err}
			}
		}

		m.loading = false
		return m, nil

	case notion.NotePagesMsg:
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			items := make([]list.Item, len(msg.Pages))
			for i, page := range msg.Pages {
				items[i] = NewItem(page)
			}
			m.list.SetItems(items)
		}

		m.loading = false
		return m, nil

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
	}

	// foward rest of commands to children
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.loading {
		return "Loading..."
	}
	if m.err != nil {
		return m.err.Error()
	}
	return m.list.View()
}
