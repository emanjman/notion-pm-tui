package projectnotes

import (
	"notion-project-tui/components/pagecontent"
	"notion-project-tui/components/projectnotelist"
	"notion-project-tui/notion"
	"notion-project-tui/styles"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type Panel int

const (
	NotesPanel Panel = iota
	ContentPanel
)

type ContentState int

const (
	ContentIdle ContentState = iota
	ContentPreview
	ContentFull
)

type ProjectNotesModel struct {
	focus        Panel
	notes        projectnotelist.ProjectNotesListModel
	content      pagecontent.PageContentModel
	keys         KeyMap
	previewCache map[string][]notion.Block
	fullCache    map[string][]notion.Block
	activeNoteID string
	contentState ContentState
	notion       *notion.Client
}

func NewProjectNotesModel(client *notion.Client) ProjectNotesModel {
	return ProjectNotesModel{
		focus:        NotesPanel,
		notes:        projectnotelist.NewProjectNotesListModel(),
		content:      pagecontent.NewPageContentModel(client),
		keys:         DefaultKeyMap,
		previewCache: map[string][]notion.Block{},
		fullCache:    map[string][]notion.Block{},
		contentState: ContentIdle,
		notion:       client,
	}
}

func (m ProjectNotesModel) Init() tea.Cmd {
	return nil
}

func (m ProjectNotesModel) Update(msg tea.Msg) (ProjectNotesModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.LeftFocus):
			m.focus = NotesPanel
			m.notes.SetFocused(true)
			return m, nil

		case key.Matches(msg, m.keys.RightFocus):
			// switch focus to content panel; if full content not cached, fetch it
			m.focus = ContentPanel
			m.notes.SetFocused(false)

			if m.activeNoteID != "" {
				if blocks, ok := m.fullCache[m.activeNoteID]; ok {
					m.content = m.content.WithContent(blocks)
					m.contentState = ContentFull
					return m, nil
				}
				return m, m.notion.FetchPageContent(m.activeNoteID)
			}
			return m, nil
		}

		// route keys to focused panel
		if m.focus == ContentPanel && m.contentState != ContentIdle {
			var cmd tea.Cmd
			m.content, cmd = m.content.Update(msg)
			return m, cmd
		}

		if m.focus == NotesPanel {
			var cmd tea.Cmd
			m.notes, cmd = m.notes.Update(msg)
			return m, cmd
		}

		return m, nil

	case tea.WindowSizeMsg:
		leftWidth := msg.Width * 30 / 100
		rightWidth := msg.Width - leftWidth - 1

		var notesCmd tea.Cmd
		m.notes, notesCmd = m.notes.Update(tea.WindowSizeMsg{
			Width:  leftWidth,
			Height: msg.Height,
		})

		var contentCmd tea.Cmd
		m.content, contentCmd = m.content.Update(tea.WindowSizeMsg{
			Width:  rightWidth,
			Height: msg.Height,
		})

		return m, tea.Batch(notesCmd, contentCmd)

	case notion.ProjectNoteMsg:
		var cmd tea.Cmd
		m.notes, cmd = m.notes.Update(msg)
		return m, cmd

	case notion.ProjectNoteSelectedMsg:
		if msg.ID == m.activeNoteID {
			return m, nil
		}

		m.activeNoteID = msg.ID

		// check fullCache first
		if blocks, ok := m.fullCache[msg.ID]; ok {
			m.content = m.content.WithContent(blocks)
			m.contentState = ContentFull
			return m, nil
		}

		// check previewCache
		if blocks, ok := m.previewCache[msg.ID]; ok {
			m.content = m.content.WithContent(blocks)
			m.contentState = ContentPreview
			return m, nil
		}

		// fetch preview
		return m, m.notion.FetchProjectNotePreview(msg.ID)

	case notion.ProjectNotePreviewMsg:
		if msg.PageID != m.activeNoteID {
			return m, nil
		}

		if msg.Err != nil {
			return m, nil
		}

		m.previewCache[msg.PageID] = msg.Blocks
		m.content = m.content.WithContent(msg.Blocks)
		m.contentState = ContentPreview
		return m, nil

	case notion.PageContentMsg:
		if msg.Err != nil {
			return m, nil
		}

		m.fullCache[m.activeNoteID] = msg.Data
		m.content = m.content.WithContent(msg.Data)
		m.contentState = ContentFull
		return m, nil
	}

	// forward other messages to both panels
	var notesCmd, contentCmd tea.Cmd
	m.notes, notesCmd = m.notes.Update(msg)
	m.content, contentCmd = m.content.Update(msg)
	return m, tea.Batch(notesCmd, contentCmd)
}

func (m ProjectNotesModel) View() string {
	left := lg.NewStyle().
		BorderRight(true).
		BorderStyle(lg.NormalBorder()).
		BorderForeground(styles.BorderForeground).
		Render(m.notes.View())
	right := m.content.View()
	return lg.JoinHorizontal(lg.Top, left, right)
}

func (m ProjectNotesModel) KeyMap() help.KeyMap {
	return m.keys
}

func (m ProjectNotesModel) InFocusMode() bool {
	return false
}
