package milestonelist

import (
	"notion-project-tui/styles"

	"github.com/charmbracelet/bubbles/textinput"
	lg "github.com/charmbracelet/lipgloss"
)

type FocusMode int

const (
	NeutralMode   FocusMode = iota
	SelectingMode           // selecting field of a milestone
	WritingMode             // editing field; reserve all keys
)

type SelectedField int

const (
	MilestoneTitle SelectedField = iota
	MilestoneTag
)

const fieldCnt = 2

// ---
// main edit state

type FocusState struct {
	milestoneID  string // page id to send update to
	milestoneIdx int    // list index to SetItem in client
	Mode         FocusMode
	field        SelectedField

	// temp state of edits
	tempTitle textinput.Model
	tempTag   string
}

// ---
// helpers

var tagFieldOptions = []string{"backend", "frontend", "infrastructure", "design", "research", "docs", "testing"}

func cycleTagField(curr string, delta int) string {
	n := len(tagFieldOptions)

	for i, tag := range tagFieldOptions {
		if tag == curr {
			return tagFieldOptions[((i+delta)%n+n)%n]
		}
	}

	return tagFieldOptions[0]
}

func initTempTitle(item MilestoneListItem) textinput.Model {
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
