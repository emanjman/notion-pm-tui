package milestone

import (
	"notion-project-tui/notion"
	"notion-project-tui/styles"

	"github.com/charmbracelet/bubbles/textinput"
	lg "github.com/charmbracelet/lipgloss"
)

// -- enum --

type Mode int

const (
	NormalMode Mode = iota
	EditMode        // editing title; reserve all keys
	DeleteMode
)

// -- types --

type EditModeCtx struct {
	milestoneID  string // page id to send update to
	milestoneIdx int    // list index to SetItem in client

	titleInput  textinput.Model
	titleBackup string // prev title in case notion-server req fails

	tempIDs int
}

type DeleteModeCtx struct {
	milestoneBackup notion.MilestonePage
}

// -- helpers --

// initialize `textinput` model
func (edit EditModeCtx) newTitleInput(item DefaultItem) textinput.Model {
	ti := textinput.New()
	ti.SetValue(item.Name)
	ti.CursorEnd()
	ti.Focus() // active
	ti.Placeholder = "Enter milestone name..."

	ti.TextStyle = lg.NewStyle().Foreground(styles.PrimaryForeground)
	ti.PlaceholderStyle = lg.NewStyle().Foreground(styles.MutedForeground)
	ti.Prompt = ""

	return ti
}

// ez switch
func (m Model) switchMode(mode Mode) Model {
	// mutate THROUGH the shared ptr (not `m.Mode = &x`) so the list delegate,
	// which holds the same *Mode, sees the switch live
	switch mode {
	case NormalMode:
		*m.Mode = NormalMode
		m.ActiveKeyMap = NormalKeyMapper
	case EditMode:
		*m.Mode = EditMode
		m.ActiveKeyMap = EditKeyMapper
	case DeleteMode:
		*m.Mode = DeleteMode
		m.ActiveKeyMap = DeleteKeyMapper
	}

	return m
}
