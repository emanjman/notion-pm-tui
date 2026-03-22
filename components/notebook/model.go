package notebook

import (
	"fmt"
	"notion-project-tui/notion"
	"notion-project-tui/styles"
	"os"
	"os/exec"
	"slices"

	"github.com/charmbracelet/bubbles/help"
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
type EditorFinishedMsg struct {
	Idx      int
	Note     *Item
	Markdown string
	Err      error
}

type NotebookState int

const (
	Browsing NotebookState = iota
	Reading
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
}

func New(notion *notion.Client, projID, notesPropID string) Model {
	l := list.New([]list.Item{}, NewItemDelegate(true), 0, 0)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	l.SetShowFilter(false)
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()

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

				if msg.Err != nil {
					note.ContentState = Failed
				} else if note.blocksReady {
					note.ContentState = Success
				}

				m.browser.SetItem(msg.Idx, note)
			}
		}
		return m, nil

	case EditorFinishedMsg:
		if msg.Err != nil {
			m.err = msg.Err
			m.State = Reading
			m.ActiveKeyMap = ReaderKeys
			return m, nil
		}
		if msg.Note != nil {
			note := *msg.Note
			idx := msg.Idx
			return m, func() tea.Msg {
				md, err := m.notion.ReplaceContentByMarkdown(note.ID, msg.Markdown)
				note.Markdown = md
				return ReplaceContentMsg{Idx: idx, Note: &note, Err: err}
			}
		}
		return m, nil

	case ReplaceContentMsg:
		if msg.Note != nil {
			m.browser.SetItem(msg.Idx, *msg.Note)
			if msg.Err == nil {
				return m, m.fetchNoteBlocks(msg.Idx, *msg.Note)
			}
		}
		m.State = Reading
		m.ActiveKeyMap = ReaderKeys
		m.browser.SetDelegate(NewItemDelegate(false))
		m.reader.SetContent(m.getCurrContent())
		return m, nil

	case tea.WindowSizeMsg:
		leftw := msg.Width * 30 / 100
		rightw := msg.Width - leftw - 1

		m.browser.SetSize(leftw, msg.Height)
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
				return m, nil

			case key.Matches(msg, m.browserKeyMap.Up):
				m.browser.CursorUp()
				m.reader.YPosition = 0
				m.reader.SetContent(m.getCurrContent())
				return m, nil

			case key.Matches(msg, m.browserKeyMap.Enter):
				curr := m.browser.SelectedItem()
				if note, ok := curr.(Item); ok {
					switch note.ContentState {
					case Idle, Failed:
						idx := m.browser.Index()
						note.ContentState = Pending
						m.browser.SetItem(idx, note)
						m.reader.SetContent(m.getCurrContent())
						return m, tea.Batch(m.fetchNoteBlocks(idx, note), m.fetchNoteMarkdown(idx, note))
					case Success:
						return m.launchEditor()
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
				return m.launchEditor()
			}
		}
	}

	var bcmd, rcmd tea.Cmd
	m.browser, bcmd = m.browser.Update(msg)
	m.reader, rcmd = m.reader.Update(msg)
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

	right := lg.NewStyle().
		Padding(0, 1).
		Render(m.reader.View())

	return lg.JoinHorizontal(lg.Top, left, right)
}

func (m Model) launchEditor() (Model, tea.Cmd) {
	note, ok := m.browser.SelectedItem().(Item)
	if !ok {
		return m, nil
	}
	idx := m.browser.Index()

	f, err := os.CreateTemp("", "notion-*.md")
	if err != nil {
		m.err = fmt.Errorf("could not create temp file: %w", err)
		return m, nil
	}
	if _, err := f.WriteString(note.Markdown); err != nil {
		f.Close()
		os.Remove(f.Name())
		m.err = fmt.Errorf("could not write temp file: %w", err)
		return m, nil
	}
	f.Close()

	tmpPath := f.Name()
	cmd := exec.Command(resolveEditor(), tmpPath)

	return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
		defer os.Remove(tmpPath)
		if err != nil {
			return EditorFinishedMsg{Idx: idx, Note: &note, Err: fmt.Errorf("editor exited with error: %w", err)}
		}
		content, readErr := os.ReadFile(tmpPath)
		if readErr != nil {
			return EditorFinishedMsg{Idx: idx, Note: &note, Err: fmt.Errorf("could not read temp file: %w", readErr)}
		}
		return EditorFinishedMsg{Idx: idx, Note: &note, Markdown: string(content)}
	})
}

func resolveEditor() string {
	if v := os.Getenv("VISUAL"); v != "" {
		return v
	}
	if e := os.Getenv("EDITOR"); e != "" {
		return e
	}
	return "vi"
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

	if note, ok := m.browser.Items()[i].(Item); ok {
		switch note.ContentState {
		case Idle:
			content = "[Enter] to fetch content"
		case Pending:
			content = "Fetching..."
		default:
			content = note.Content
		}
	}
	return content
}

func (m Model) buildNoteList(pages []notion.NotePage) []list.Item {
	items := make([]list.Item, len(pages))
	for i, page := range pages {
		items[i] = NewItem(page)
	}

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
