package notebook

import (
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
type FetchNoteMarkdownMsg struct {
	Idx  int
	Err  error
	Note *Item
}
type ReplaceContentMsg struct {
	Idx  int
	Err  error
	Note *Item
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
			blxCmds, mdCmds := make([]tea.Cmd, reqCnt), make([]tea.Cmd, reqCnt)
			for i := range reqCnt {
				item := m.browser.Items()[i]
				if note, ok := item.(Item); ok {
					note.ContentState = Pending
					m.browser.SetItem(i, note)

					blxCmds[i] = m.fetchNoteBlocks(i, note)
					mdCmds[i] = m.fetchNoteMarkdown(i, note)
				}
			}
			cmds := append(blxCmds, mdCmds...)
			return m, tea.Batch(cmds...)
		}

	case FetchNoteContentMsg:
		if curr := m.browser.Items()[msg.Idx]; curr != nil {
			if note, ok := curr.(Item); ok {
				note.Content = msg.Note.Content
				note.blocksReady = true

				// Update combined state
				if msg.Err != nil {
					note.ContentState = Failed
				} else if note.markdownReady {
					note.ContentState = Success
				}

				m.browser.SetItem(msg.Idx, note)
			}
		}
		m.reader.SetContent(m.getCurrContent())
		return m, nil

	case FetchNoteMarkdownMsg:
		if curr := m.browser.Items()[msg.Idx]; curr != nil {
			if note, ok := curr.(Item); ok {
				note.Markdown = msg.Note.Markdown
				note.markdownReady = true

				// update combined state
				if msg.Err != nil {
					note.ContentState = Failed
				} else if note.blocksReady {
					note.ContentState = Success
				}

				m.browser.SetItem(msg.Idx, note)
			}
		}
		m.editor.SetValue(m.getCurrMarkdown())

	case ReplaceContentMsg:
		if msg.Note != nil {
			// Update the item with new markdown from API
			m.browser.SetItem(msg.Idx, *msg.Note)

			// Re-fetch blocks to reflect the new content
			if msg.Err == nil {
				return m, m.fetchNoteBlocks(msg.Idx, *msg.Note)
			}
		}
		m.reader.SetContent(m.getCurrContent())
		m.editor.SetValue(m.getCurrMarkdown())
		return m, nil

	case tea.WindowSizeMsg:
		leftw := msg.Width * 30 / 100
		rightw := msg.Width - leftw - 1 // mind the border

		m.browser.SetSize(leftw, msg.Height)
		m.editor.SetWidth(rightw)
		m.editor.SetHeight(msg.Height)
		m.reader.Width, m.reader.Height = rightw, msg.Height
		return m, nil

	case tea.KeyMsg:
		switch m.State {
		case Browsing:
			switch {
			case key.Matches(msg, m.browserKeyMap.Right):
				if note, ok := m.browser.SelectedItem().(Item); ok && note.ContentState == Success {
					m.State = Reading
					m.ActiveKeyMap = ReaderKeys
					m.browser.SetDelegate(NewItemDelegate(false))
				}
				return m, nil

			case key.Matches(msg, m.browserKeyMap.Down):
				m.browser.CursorDown()
				m.reader.YPosition = 0
				m.reader.SetContent(m.getCurrContent())
				m.editor.SetValue(m.getCurrMarkdown())
				return m, nil

			case key.Matches(msg, m.browserKeyMap.Up):
				m.browser.CursorUp()
				m.reader.YPosition = 0
				m.reader.SetContent(m.getCurrContent())
				m.editor.SetValue(m.getCurrMarkdown())
				return m, nil

			case key.Matches(msg, m.browserKeyMap.Enter):
				curr := m.browser.SelectedItem()
				if note, ok := curr.(Item); ok {
					switch note.ContentState {
					case Idle, Failed:
						idx := m.browser.Index()
						note.ContentState = Pending
						m.browser.SetItem(idx, note)
						m.reader.SetContent(m.getCurrContent()) // show pending state
						return m, tea.Batch(m.fetchNoteBlocks(idx, note), m.fetchNoteMarkdown(idx, note))
					case Success:
						m.State = Editing
						m.ActiveKeyMap = EditorKeys
						m.browser.SetDelegate(NewItemDelegate(false))
					}
				}
				return m, nil
			}
		case Reading:
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
				m.browser.SetDelegate(NewItemDelegate(false))
			}
		case Editing:
			switch {
			case key.Matches(msg, m.editorKeyMap.Esc):
				m.State = Browsing
				m.ActiveKeyMap = BrowserKeys
				m.browser.SetDelegate(NewItemDelegate(true))

				// submit changes to notion
				if item, ok := m.browser.SelectedItem().(Item); ok {
					idx := m.browser.Index()

					return m, func() tea.Msg {
						md, err := m.notion.ReplaceContentByMarkdown(item.ID, m.editor.Value())
						item.Markdown = md
						if err != nil {
							return ReplaceContentMsg{Idx: idx, Note: &item, Err: err}
						}
						return ReplaceContentMsg{Idx: idx, Note: &item, Err: nil}
					}
				}

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

func (m Model) fetchNoteBlocks(idx int, note Item) tea.Cmd {
	return func() tea.Msg {
		blocks, err := m.notion.FetchPageBlocks(note.ID)
		if err != nil {
			note.Content = err.Error()
		} else {
			note.Content = notion.BlocksToContent(blocks, m.reader.Width, 0)
		}
		return FetchNoteContentMsg{Idx: idx, Note: &note, Err: err}
	}
}

func (m Model) fetchNoteMarkdown(idx int, note Item) tea.Cmd {
	return func() tea.Msg {
		md, err := m.notion.FetchPageMarkdown(note.ID)
		if err != nil {
			note.Markdown = err.Error()
		} else {
			note.Markdown = md
		}
		return FetchNoteMarkdownMsg{Idx: idx, Note: &note, Err: err}
	}
}

func (m Model) getCurrContent() string {
	content := "Unable to render"
	i := m.browser.Index()

	// pull from actual items, not filteredItems
	if note, ok := m.browser.Items()[i].(Item); ok {
		switch note.ContentState {
		case Idle:
			content = "[Enter] to fetch content"
		case Pending:
			content = "Fetching..."
		default:
			content = note.Content // may be blocks or the err msg
		}
	}
	return content
}

func (m Model) getCurrMarkdown() string {
	if note, ok := m.browser.Items()[m.browser.Index()].(Item); ok {
		return note.Markdown
	}
	return ""
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
