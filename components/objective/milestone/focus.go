package milestone

import (
	"notion-project-tui/styles"

	"github.com/charmbracelet/bubbles/textinput"
	lg "github.com/charmbracelet/lipgloss"
)

type FocusMode int

const (
	NeutralMode FocusMode = iota
	WritingMode             // editing title; reserve all keys
)

// ---
// main edit state

type FocusState struct {
	milestoneID  string // page id to send update to
	milestoneIdx int    // list index to SetItem in client
	Mode         FocusMode

	// temp state of edits
	tempTitle textinput.Model
}

// ---
// helpers

func initTempTitle(item Item) textinput.Model {
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
