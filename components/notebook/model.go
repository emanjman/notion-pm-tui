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

	State            NotebookState
	browser          list.Model
	browserKeyMap    BrowserKeyMap
	reader           viewport.Model
	readerKeyMap     ReaderKeyMap
	editor           textarea.Model
	editorKeyMap     EditorKeyMap
	vimMode          VimMode
	pendingVimKey    string
	vimCount         string
	pendingInsertEsc bool
	ogMarkdown       string
}

func New(notion *notion.Client, projID, propID string) Model {
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
	ta.FocusedStyle.LineNumber = ta.FocusedStyle.LineNumber.
		Foreground(styles.MutedForeground)
	ta.Focus()

	return Model{
		projID:       projID,
		notesPropID:  propID,
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
			m.loading = false
			return m, nil
		}
		return m, func() tea.Msg {
			pages, err := notion.FetchPages[notion.NotePage](m.notion, msg.IDs)
			return notion.NotePagesMsg{Pages: pages, Err: err}
		}

	case notion.NotePagesMsg:
		if msg.Err != nil {
			m.err = msg.Err
			m.loading = false
			return m, nil
		}

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

	case EditorFinishedMsg:
		return m, func() tea.Msg {
			md, err := m.notion.ReplaceContentByMarkdown(msg.Note.ID, msg.Content)
			msg.Note.Markdown = md
			return ReplaceContentMsg{Note: msg.Note, Idx: msg.Idx, Err: err}
		}

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

	case insertEscTimeoutMsg:
		if m.pendingInsertEsc {
			m.pendingInsertEsc = false
			if m.State == Editing && m.vimMode == InsertMode {
				var cmd tea.Cmd
				m.editor, cmd = m.editor.Update(runeK('j'))
				return m, cmd
			}
		}
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
				m = m.updateContentOnNav(true)
				return m, nil
			case key.Matches(msg, m.browserKeyMap.Up):
				m = m.updateContentOnNav(false)
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
						m = m.enterEditMode()
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
				m = m.enterEditMode()

			case key.Matches(msg, m.readerKeyMap.OpenEditor):
				idx := m.browser.Index()
				note := m.browser.Items()[idx]
				if note, ok := note.(Item); ok {
					return m, openMarkdownInEditor(m.getCurrMarkdown(), idx, note)
				}
				return m, nil
			}
		case Editing:
			switch {
			case msg.String() == "ctrl+j":
				var cmd tea.Cmd
				m.editor, cmd = sendKeys(m.editor,
					k(tea.KeyDown), k(tea.KeyDown), k(tea.KeyDown), k(tea.KeyDown), k(tea.KeyDown),
				)
				return m, cmd

			case msg.String() == "ctrl+k":
				var cmd tea.Cmd
				m.editor, cmd = sendKeys(m.editor,
					k(tea.KeyUp), k(tea.KeyUp), k(tea.KeyUp), k(tea.KeyUp), k(tea.KeyUp),
				)
				return m, cmd

			case key.Matches(msg, m.editorKeyMap.Esc):
				if m.vimMode == InsertMode {
					// Esc in insert mode → back to normal mode
					m.vimMode = NormalMode
					return m, nil
				}
				// Esc in normal mode → save and exit to browsing
				m.State = Browsing
				m.ActiveKeyMap = BrowserKeys
				m.browser.SetDelegate(NewItemDelegate(true))

				if item, ok := m.browser.SelectedItem().(Item); ok && m.editor.Value() != m.ogMarkdown {
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

			default:
				if m.vimMode == NormalMode {
					return handleNormalMode(m, msg)
				}
				// InsertMode: check for jk escape sequence
				if m.pendingInsertEsc {
					m.pendingInsertEsc = false
					if msg.String() == "k" {
						m.vimMode = NormalMode
						return m, nil
					}
					// not jk — flush the held j, then fall through to handle current key
					m.editor, _ = m.editor.Update(runeK('j'))
				}
				if msg.String() == "j" {
					m.pendingInsertEsc = true
					return m, insertEscTimeout()
				}
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
		modeLabel := " NORMAL "
		modeStyle := lg.NewStyle().
			Background(styles.MutedForeground).
			Foreground(lg.Color("#000000")).
			Bold(true)
		if m.vimMode == InsertMode {
			modeLabel = " INSERT "
			modeStyle = modeStyle.Background(lg.Color("#6fb7b7"))
		}
		rightContent = lg.JoinVertical(lg.Left,
			m.editor.View(),
			modeStyle.Render(modeLabel),
		)
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

func (m Model) enterEditMode() Model {
	m.State = Editing
	m.ActiveKeyMap = EditorKeys
	m.browser.SetDelegate(NewItemDelegate(false))
	m.ogMarkdown = m.editor.Value()
	m.vimMode = NormalMode
	return m
}

func (m Model) updateContentOnNav(scrollDown bool) Model {
	if scrollDown {
		m.browser.CursorDown()
	} else {
		m.browser.CursorUp()
	}
	m.reader.YPosition = 0
	m.reader.SetContent(m.getCurrContent())
	m.editor.SetValue(m.getCurrMarkdown())
	return m
}
