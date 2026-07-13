package task

import (
	"notion-project-tui/notion"
	"notion-project-tui/styles"

	"github.com/charmbracelet/bubbles/textinput"
	lg "github.com/charmbracelet/lipgloss"
)

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
	prevTitle    string
	prevType     string
	prevPriority int

	// deletion confirmation
	pendingDelete bool
}

// ---
// helpers

func cycleTypeField(curr string, delta int, options []notion.SelectItem) string {
	n := len(options)
	if n == 0 {
		return curr
	}
	for i, opt := range options {
		if opt.Name == curr {
			return options[((i+delta)%n+n)%n].Name
		}
	}
	return options[0].Name
}

func cyclePriorityField(curr, delta int) int {
	const n = 6
	return ((curr+delta)%n + n) % n // todo: vet logic
}

func cycleStatus(curr notion.TaskStatus, delta int) notion.TaskStatus {
	order := notion.TaskStatusOrder()

	for i, status := range order {
		if status == curr {
			newIdx := i + delta

			// clamp at boundaries (don't wrap around)
			if newIdx < 0 {
				return curr // stay at idle
			}
			if newIdx >= len(order) {
				return curr // stay at done
			}

			return order[newIdx]
		}
	}
	return notion.TaskIdle // default
}

func initTempTitle(item Item) textinput.Model {
	ti := textinput.New()
	ti.SetValue(item.Task)
	ti.CursorEnd()
	ti.Focus() // active
	ti.Placeholder = "Enter task..."

	ti.TextStyle = lg.NewStyle().Foreground(styles.PrimaryForeground)
	ti.PlaceholderStyle = lg.NewStyle().Foreground(styles.MutedForeground)
	ti.Prompt = ""

	return ti

}
