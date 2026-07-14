package task

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
	SelectMode      // selecting field of a task
	EditMode        // editing field; reserve all keys
	DeleteMode
)

type SelectedField int

const (
	TaskType SelectedField = iota
	TaskTitle
	TaskPriority
	_SelectedFieldCount
)

// -- types --

type SelectModeCtx struct {
	field SelectedField
}

type EditModeCtx struct {
	taskID  string // page id to send update to
	taskIdx int    // list index to SetItem in client

	titleInput     textinput.Model
	titleBackup    string
	typeBackup     string
	priorityBackup int

	pendingDelete bool // delete confirmation
}

// todo: is this doable?
type DeleteModeCtx struct {
	pendingDelete bool            // delete confirmation
	taskBackup    notion.TaskPage // todo: mirrors mstone
}

// -- helpers --

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

// // ez switch
// func (m Model) switchMode(mode Mode) Model {
// 	// mutate THROUGH the shared ptr (not `m.Mode = &x`) so the list delegate,
// 	// which holds the same *Mode, sees the switch live
// 	switch mode {
// 	case NormalMode:
// 		*m.Mode = NormalMode
// 		m.ActiveKeyMap = NormalKeyMapper
// 	case EditMode:
// 		*m.Mode = EditMode
// 		m.ActiveKeyMap = EditKeyMapper
// 	case DeleteMode:
// 		*m.Mode = DeleteMode
// 		m.ActiveKeyMap = DeleteKeyMapper
// 	}
// 	return m
// }
