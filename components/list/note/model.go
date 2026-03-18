package note

import (
	"notion-project-tui/components/pagecontent"
	"notion-project-tui/notion"
	"notion-project-tui/styles"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type Model struct {
	projID      string
	notesPropID string
	loading     bool
	browsing    bool // focused on notes list
	err         error
	notion      *notion.Client
	keys        KeyMap

	browser list.Model
	reader  pagecontent.Model
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
		loading:     true,
		browsing:    true,
		err:         nil,
		notion:      notion,

		browser: l,
		reader:  pagecontent.New(notion),
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
			m.browser.SetItems(items)
		}
		m.loading = false
		return m, nil

	case tea.WindowSizeMsg:
		var readerCmd tea.Cmd
		leftw := msg.Width * 30 / 100
		rightw := msg.Width - leftw - 1 // div border

		m.browser.SetSize(leftw, msg.Height) // ? do we reorganize for greater sep of concerns?
		m.reader, readerCmd = m.reader.Update(tea.WindowSizeMsg{
			Width:  rightw,
			Height: msg.Height,
		})

		return m, readerCmd

	// toggle views
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.LeftFocus):
			if !m.browsing {
				m.browsing = true
				m.browser.SetDelegate(NewItemDelegate(true))
			}
		case key.Matches(msg, m.keys.RightFocus):
			if m.browsing {
				m.browsing = false
				m.browser.SetDelegate(NewItemDelegate(false))
			}
		}
	}

	// foward rest of commands to children
	var cmd tea.Cmd
	m.browser, cmd = m.browser.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	browserContent := m.browser.View()
	if m.loading {
		browserContent = "Loading..."
	}
	if m.err != nil {
		browserContent = m.err.Error()
	}

	left := lg.NewStyle().
		BorderRight(true).
		BorderStyle(lg.NormalBorder()).
		BorderForeground(styles.BorderForeground).
		Render(browserContent)

	right := m.reader.View()
	return lg.JoinHorizontal(lg.Top, left, right)
}
