package notebook

import (
	"log"
	"notion-project-tui/notion"
	"notion-project-tui/styles"
	"slices"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
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

type NotebookState int

const (
	Browsing NotebookState = iota
	Reading
	Editing
)

type Model struct {
	projID       string
	notesPropID  string
	loading      bool
	err          error
	notion       *notion.Client
	ActiveKeyMap help.KeyMap

	State         NotebookState
	browser       list.Model
	browserKeyMap BrowserKeyMap
	reader        viewport.Model
	readerKeyMap  ReaderKeyMap
	editor        textarea.Model
	editorKeyMap  EditorKeyMap
}

func New(notion *notion.Client, projID, notesPropID string) Model {
	// list config
	l := list.New([]list.Item{}, NewItemDelegate(true), 0, 0)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	l.SetShowFilter(false)
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()

	// text area config
	ta := textarea.New()
	ta.Focus()
	ta.SetValue("init text")
	ta.ShowLineNumbers = false

	return Model{
		projID:       projID,
		notesPropID:  notesPropID,
		loading:      true,
		err:          nil,
		notion:       notion,
		ActiveKeyMap: BrowserKeys,

		State:         Browsing,
		browser:       l,
		browserKeyMap: BrowserKeys,
		reader:        viewport.New(0, 0),
		readerKeyMap:  ReaderKeys,
		editor:        ta,
		editorKeyMap:  EditorKeys,
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
			items := m.buildNoteList(msg.Pages)
			m.browser.SetItems(items)
			m.loading = false

			reqCnt := 5
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
		m.reader.SetContent(m.getCurrContent())
		return m, nil

	case ItemStateMsg:
		temp := m.browser.Items()[msg.Idx]
		if note, ok := temp.(Item); ok {
			note.State = msg.State
			m.browser.SetItem(msg.Idx, note)
		}
		return m, nil

	case tea.WindowSizeMsg:
		leftw := msg.Width * 30 / 100
		rightw := msg.Width - leftw - 1 // mind the border
		m.browser.SetSize(leftw, msg.Height)
		m.reader.Width, m.reader.Height = rightw, msg.Height
		return m, nil

	case tea.KeyMsg:
		if m.State == Browsing {
			switch {
			case key.Matches(msg, m.browserKeyMap.Right):
				if note, ok := m.browser.SelectedItem().(Item); ok && note.State == Success {
					m.State = Reading
					m.ActiveKeyMap = ReaderKeys
					m.browser.SetDelegate(NewItemDelegate(false))
				}
				return m, nil

			case key.Matches(msg, m.browserKeyMap.Down):
				m.browser.CursorDown()
				m.reader.YPosition = 0
				m.reader.SetContent(m.getCurrContent())
				return m, nil

			case key.Matches(msg, m.browserKeyMap.Up):
				m.browser.CursorUp()
				m.reader.YPosition = 0
				m.reader.SetContent(m.getCurrContent())
				return m, nil

			case key.Matches(msg, m.browserKeyMap.Enter):
				selected := m.browser.SelectedItem()
				if note, ok := selected.(Item); ok && note.State == Idle {
					idx := m.browser.Index()
					note.State = Pending
					m.browser.SetItem(idx, note)
					m.reader.SetContent(m.getCurrContent()) // show pending state
					return m, m.fetchNoteContent(idx, note)
				}
				return m, nil
			}
		} else if m.State == Reading {
			switch {
			case key.Matches(msg, m.readerKeyMap.Left):
				m.State = Browsing
				m.ActiveKeyMap = BrowserKeys
				m.browser.SetDelegate(NewItemDelegate(true))
				return m, nil

			case key.Matches(msg, m.readerKeyMap.Up5):
				m.reader.ScrollUp(5)
			case key.Matches(msg, m.readerKeyMap.Down5):
				m.reader.ScrollDown(5)

			case key.Matches(msg, m.readerKeyMap.Up), key.Matches(msg, m.readerKeyMap.Down):
				var cmd tea.Cmd
				m.reader, cmd = m.reader.Update(msg)
				return m, cmd

			case key.Matches(msg, m.readerKeyMap.Enter):
				m.State = Editing
				m.ActiveKeyMap = EditorKeys
			}
		} else if m.State == Editing {
			switch {
			case key.Matches(msg, m.editorKeyMap.Esc):
				m.State = Reading
				m.ActiveKeyMap = ReaderKeys
				// todo: submit changes to notion

			// forward all keys into textarea model
			default:
				var cmd tea.Cmd
				m.editor, cmd = m.editor.Update(msg)
				return m, cmd
			}
		}
	}

	// foward rest of commands to children
	var bcmd, rcmd tea.Cmd
	m.browser, bcmd = m.browser.Update(msg)
	m.reader, rcmd = m.reader.Update(msg)
	// m.editor, ecmd = m.editor.Update(msg)
	return m, tea.Batch(bcmd, rcmd)
}

func (m Model) View() string {
	leftContent := m.browser.View()
	if m.loading {
		leftContent = "Loading..."
	}
	if m.err != nil {
		leftContent = m.err.Error()
	}

	left := lg.NewStyle().
		BorderRight(true).
		BorderStyle(lg.NormalBorder()).
		BorderForeground(styles.BorderForeground).
		Render(leftContent)

	rightContent := m.reader.View()
	if m.State == Editing {
		rightContent = m.editor.View()
	}

	right := lg.NewStyle().
		Padding(0, 1).
		Render(rightContent)

	return lg.JoinHorizontal(lg.Top, left, right)
}

func (m Model) fetchNoteContent(idx int, note Item) tea.Cmd {
	return func() tea.Msg {
		blocks, err := m.notion.FetchPageContent(note.ID)
		log.Printf("fetching blocks for note: " + note.Title)

		if err != nil {
			note.Content = err.Error()
			note.State = Failed
		} else {
			note.Content = notion.RenderBlocks(blocks, m.reader.Width, 0)
			note.State = Success
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
	i := m.browser.Index()

	// pull from actual items, not filteredItems
	if note, ok := m.browser.Items()[i].(Item); ok {
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

func (m Model) buildNoteList(pages []notion.NotePage) []list.Item {
	items := make([]list.Item, len(pages))
	for i, page := range pages {
		items[i] = NewItem(page)
	}

	// sort by desc (most recent, first)
	slices.SortFunc(items, func(a, b list.Item) int {
		noteA, okA := a.(Item)
		noteB, okB := b.(Item)
		if !okA || !okB {
			return 0
		}
		return noteB.CreatedDate.Compare(noteA.CreatedDate)
	})
	return items
}
