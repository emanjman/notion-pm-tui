package note

import (
	"log"
	"notion-project-tui/notion"
	"notion-project-tui/styles"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type SwitchContentMsg struct{ Content string }
type FetchNoteContentMsg struct {
	Idx  int
	Note *Item
}

type Model struct {
	projID      string
	notesPropID string
	loading     bool
	browsing    bool // focused on notes list
	err         error
	notion      *notion.Client
	keys        KeyMap

	browser list.Model
	reader  viewport.Model
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
		keys:        DefaultKeyMap,

		browser: l,
		reader:  viewport.New(0, 0),
	}
}

func (m Model) Init() tea.Cmd {
	browserInit := func() tea.Msg {
		ids, err := m.notion.FetchRelationIDs(m.projID, m.notesPropID)
		return notion.NoteIDsMsg{IDs: ids, Err: err}
	}
	return tea.Batch(browserInit, m.reader.Init())
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

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
			m.loading = false
			return m, nil
		} else {
			items := make([]list.Item, len(msg.Pages))
			for i, page := range msg.Pages {
				items[i] = NewItem(page)
			}
			m.browser.SetItems(items)
			m.loading = false

			log.Printf("here we have received note pages")
			// return m, m.fetchAllNoteContent()

			cmds := make([]tea.Cmd, len(m.browser.Items()))
			for i, item := range m.browser.Items() {
				if note, ok := item.(Item); ok {
					cmds[i] = m.fetchNoteContent(i, note)
				}
			}
			return m, tea.Batch(cmds...)
		}

	case FetchNoteContentMsg:
		m.browser.SetItem(msg.Idx, *msg.Note)
		return m, m.displayCurrentContent()

	case SwitchContentMsg:
		m.reader.SetContent(msg.Content)

	case tea.WindowSizeMsg:
		var readerCmd tea.Cmd
		leftw := msg.Width * 30 / 100
		rightw := msg.Width - leftw - 1 // div border

		m.browser.SetSize(leftw, msg.Height)
		m.reader.Width, m.reader.Height = rightw, msg.Height

		return m, readerCmd

	case tea.KeyMsg:
		if m.browsing {
			switch {
			case key.Matches(msg, m.keys.RightFocus):
				m.browsing = false
				m.browser.SetDelegate(NewItemDelegate(false))
				return m, nil
			case key.Matches(msg, m.keys.Up), key.Matches(msg, m.keys.Down):
				var cmd tea.Cmd
				m.browser, cmd = m.browser.Update(msg)
				return m, tea.Batch(m.displayCurrentContent(), cmd)
			}
		} else {
			switch {
			case key.Matches(msg, m.keys.LeftFocus):
				m.browsing = true
				m.browser.SetDelegate(NewItemDelegate(true))
				return m, nil
			}
		}
	}

	// foward rest of commands to children
	var bcmd, rcmd tea.Cmd
	m.browser, bcmd = m.browser.Update(msg)
	m.reader, rcmd = m.reader.Update(msg)
	return m, tea.Batch(bcmd, rcmd)
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
	right := lg.NewStyle().
		Padding(0, 1).
		Render(m.reader.View())
	return lg.JoinHorizontal(lg.Top, left, right)
}

func (m Model) fetchNoteContent(idx int, note Item) tea.Cmd {
	return func() tea.Msg {
		blocks, err := m.notion.FetchPageContent(note.ID)
		log.Printf("fetching blocks for note: " + note.Title)

		if err != nil {
			note.Content = err.Error()
		} else {
			note.Content = notion.RenderBlocks(blocks, m.reader.Width, 0)
		}

		note.Loading = false
		return FetchNoteContentMsg{Idx: idx, Note: &note}
	}
}

func (m Model) displayCurrentContent() tea.Cmd {
	return func() tea.Msg {
		content := "Unable to render"
		if note, ok := m.browser.SelectedItem().(Item); ok {
			if note.Loading {
				content = "Loading..."
			} else {
				content = note.Content // may be blocks or the err msg
			}
		}
		return SwitchContentMsg{Content: content}
	}
}
