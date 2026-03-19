package notebook

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

type FetchNoteContentMsg struct {
	Idx  int
	Err  error
	Note *Item
}
type ItemStateMsg struct {
	Idx   int
	State ItemState
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

			reqCnt := 3
			fetchCmds, stateCmds := make([]tea.Cmd, reqCnt), make([]tea.Cmd, reqCnt)
			for i := range reqCnt {
				item := m.browser.Items()[i]
				if note, ok := item.(Item); ok {
					stateCmds[i] = m.emitItemState(i, Pending)
					fetchCmds[i] = m.fetchNoteContent(i, note)
				}
			}
			cmds := append(fetchCmds, stateCmds...)
			return m, tea.Batch(cmds...)
		}

	case FetchNoteContentMsg:
		m.browser.SetItem(msg.Idx, *msg.Note)
		m.reader.SetContent(m.getCurrContent()) // TODO: this isn't instant?

		var emitState tea.Cmd
		if msg.Err != nil {
			emitState = m.emitItemState(msg.Idx, Failed)
		} else {
			emitState = m.emitItemState(msg.Idx, Success)
		}
		return m, emitState

	case ItemStateMsg:
		temp := m.browser.Items()[msg.Idx]
		if note, ok := temp.(Item); ok {
			note.State = msg.State
			m.browser.SetItem(msg.Idx, note)
		}
		return m, nil

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

			case key.Matches(msg, m.keys.Down):
				m.browser.CursorDown()
				m.reader.SetContent(m.getCurrContent())
				return m, nil

			case key.Matches(msg, m.keys.Up):
				m.browser.CursorUp()
				m.reader.SetContent(m.getCurrContent())
				return m, nil

			case key.Matches(msg, m.keys.Enter):
				selected := m.browser.SelectedItem()
				if note, ok := selected.(Item); ok && note.State == Idle {
					idx := m.browser.Index()
					stateCmd := m.emitItemState(idx, Pending)
					fetchCmd := m.fetchNoteContent(idx, note)
					return m, tea.Batch(stateCmd, fetchCmd)
				}
				return m, nil
			}
		} else {
			switch {
			case key.Matches(msg, m.keys.LeftFocus):
				m.browsing = true
				m.browser.SetDelegate(NewItemDelegate(true))
				return m, nil
			case key.Matches(msg, m.keys.Up), key.Matches(msg, m.keys.Down):
				var cmd tea.Cmd
				m.reader, cmd = m.reader.Update(msg)
				return m, cmd
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

		return FetchNoteContentMsg{Idx: idx, Note: &note, Err: err}
	}
}

func (m Model) emitItemState(idx int, state ItemState) tea.Cmd {
	return func() tea.Msg {
		return ItemStateMsg{Idx: idx, State: state}
	}
}

func (m Model) getCurrContent() string {
	content := "Unable to render"
	if note, ok := m.browser.SelectedItem().(Item); ok {
		if note.State == Idle {
			content = "[Enter] to fetch content"
		} else if note.State == Pending {
			content = "Fetching..."
		} else {
			content = note.Content // may be blocks or the err msg
		}
	}
	return content
}
