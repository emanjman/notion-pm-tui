package tasklist

import "github.com/charmbracelet/bubbles/textinput"

type FocusMode int

const (
	NeutralMode   FocusMode = iota
	SelectingMode           // selecting field of a task
	WritingMode             // editing field; reserve all keys
)

type SelectedField int

const (
	TaskType SelectedField = iota
	TaskTitle
	TaskPriority
)
const fieldCnt = 3

// ---
// main edit state

type FocusState struct {
	taskID  string // page id to send update to
	taskIdx int    // list index to SetItem in client
	Mode    FocusMode
	field   SelectedField

	// temp state of edits
	tempType     string
	tempTitle    textinput.Model
	tempPriority int
}

// ---
// helpers

var typeFieldOptions = []string{"feat", "fix", "chore", "refactor", "style"}

func cycleTypeField(curr string, delta int) string {
	n := len(typeFieldOptions)

	for i, typ := range typeFieldOptions {
		if typ == curr {
			return typeFieldOptions[((i+delta)%n+n)%n]
		}
	}

	return typeFieldOptions[0]
}

func cyclePriorityField(curr, delta int) int {
	const n = 6
	return ((curr+delta)%n + n) % n
}
