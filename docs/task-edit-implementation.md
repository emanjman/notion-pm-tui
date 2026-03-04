# Task Inline Editing — Implementation Guide

## Overview

This guide walks through implementing 2-deep inline task editing. The task list already has a skeleton: Level 1 (field selection) works, but the `EnableEdit` keypress in that mode is a TODO. This fills in the rest.

### State Machine

```
Normal    →  [Enter on task]  →  Level 1  (h/l navigates between fields)
Level 1   →  [Enter]          →  Level 2  (locked onto a field, editing it)
Level 1   →  [Esc]            →  Normal
Level 2   →  [Enter / Esc]    →  Level 1  (both commit the current value)
```

### Field Behavior in Level 2

| Field    | Interaction          | Visual                   |
|----------|----------------------|--------------------------|
| Type     | h/l cycles options   | `‹ feat ›`               |
| Priority | h/l cycles 0–5       | `‹ [3] ›`                |
| Task     | full text input      | inline textinput w/cursor |

---

## Files to Change

- `components/tasklist/model.go`
- `components/tasklist/delegate.go`
- `components/objective/model.go`

---

## Step 1 — Extend `EditState` (`model.go`)

Add the fields needed for Level 2 state. The existing struct is:

```go
type EditState struct {
    active bool
    taskID string
    field  FieldIndex
}
```

Replace it with:

```go
type EditState struct {
    active    bool
    taskID    string
    taskIndex int   // list index of the item being edited (for SetItem)
    field     FieldIndex
    subActive bool  // true = Level 2 engaged

    // Level 2 temp state (only one is in use at a time, depending on field)
    TextInput    textinput.Model
    TempType     string
    TempPriority int
}
```

Add the import at the top of the file:

```go
"github.com/charmbracelet/bubbles/textinput"
```

(`textinput` is part of `github.com/charmbracelet/bubbles` which is already in `go.mod`)

Also add the `lipgloss` import alias if it's not already there:

```go
lg "github.com/charmbracelet/lipgloss"
```

---

## Step 2 — Cycling Helpers (`model.go`)

Add these package-level vars and helpers below the `const fieldCnt = 3` line:

```go
var typeOptions = []string{"feat", "fix", "chore", "refactor", "style"}

func cycleType(current string, delta int) string {
    for i, t := range typeOptions {
        if t == current {
            n := len(typeOptions)
            return typeOptions[((i+delta)%n+n)%n]
        }
    }
    return typeOptions[0]
}

func cyclePriority(current, delta int) int {
    const n = 6 // 0–5
    return ((current+delta)%n+n) % n
}
```

---

## Step 3 — `commitSubEdit` Helper (`model.go`)

Add this function (not a method, just a package-level helper):

```go
func commitSubEdit(m TaskListModel) TaskListModel {
    items := m.list.Items()
    item, ok := items[m.EditState.taskIndex].(TaskListItem)
    if !ok {
        return m
    }
    switch m.EditState.field {
    case TypeField:
        item.Type = m.EditState.TempType
    case PriorityField:
        item.Priority = m.EditState.TempPriority
    case TaskField:
        item.Task = m.EditState.TextInput.Value()
    }
    m.list.SetItem(m.EditState.taskIndex, item)
    return m
}
```

> Note: `m.groups` is intentionally not updated here. It only matters when rebuilding the list (e.g. on milestone re-selection). Since API persistence comes later, this is fine for now.

---

## Step 4 — `IsCapturingTextInput` Method (`model.go`)

Add this method to `TaskListModel`:

```go
func (m TaskListModel) IsCapturingTextInput() bool {
    return m.EditState.active && m.EditState.subActive && m.EditState.field == TaskField
}
```

This is used by `objective` to know when to let h/l pass through to the text input instead of switching panels.

---

## Step 5 — Key Handling in `Update` (`model.go`)

The current edit-mode key handler block looks like:

```go
if m.EditState.active {
    switch {
    case key.Matches(msg, m.EditKeys.Exit):
        ...
    case key.Matches(msg, m.EditKeys.PrevField):
        ...
    case key.Matches(msg, m.EditKeys.NextField):
        ...
    case key.Matches(msg, m.EditKeys.EnableEdit):
        // todo: enter 2-deep edit mode
    }
    return m, nil
}
```

Replace the entire `if m.EditState.active { ... }` block with:

```go
if m.EditState.active {

    // --- Level 2: sub-edit is engaged ---
    if m.EditState.subActive {
        switch m.EditState.field {

        case TypeField:
            switch {
            case key.Matches(msg, m.EditKeys.PrevField):
                m.EditState.TempType = cycleType(m.EditState.TempType, -1)
            case key.Matches(msg, m.EditKeys.NextField):
                m.EditState.TempType = cycleType(m.EditState.TempType, +1)
            case key.Matches(msg, m.EditKeys.EnableEdit), key.Matches(msg, m.EditKeys.Exit):
                m = commitSubEdit(m)
                m.EditState.subActive = false
            }

        case PriorityField:
            switch {
            case key.Matches(msg, m.EditKeys.PrevField):
                m.EditState.TempPriority = cyclePriority(m.EditState.TempPriority, -1)
            case key.Matches(msg, m.EditKeys.NextField):
                m.EditState.TempPriority = cyclePriority(m.EditState.TempPriority, +1)
            case key.Matches(msg, m.EditKeys.EnableEdit), key.Matches(msg, m.EditKeys.Exit):
                m = commitSubEdit(m)
                m.EditState.subActive = false
            }

        case TaskField:
            if key.Matches(msg, m.EditKeys.EnableEdit) || key.Matches(msg, m.EditKeys.Exit) {
                m = commitSubEdit(m)
                m.EditState.subActive = false
            } else {
                var cmd tea.Cmd
                m.EditState.TextInput, cmd = m.EditState.TextInput.Update(msg)
                return m, cmd
            }
        }

        return m, nil
    }

    // --- Level 1: field selection ---
    switch {

    case key.Matches(msg, m.EditKeys.Exit):
        m.EditState.active = false
        // todo: send some command to save notion changes
        return m, nil

    case key.Matches(msg, m.EditKeys.PrevField):
        if m.EditState.field == TypeField {
            m.EditState.field = fieldCnt - 1
        } else {
            m.EditState.field = (m.EditState.field - 1) % fieldCnt
        }
        return m, nil

    case key.Matches(msg, m.EditKeys.NextField):
        m.EditState.field = (m.EditState.field + 1) % fieldCnt
        return m, nil

    case key.Matches(msg, m.EditKeys.EnableEdit):
        item, ok := m.list.Items()[m.list.Index()].(TaskListItem)
        if !ok {
            return m, nil
        }
        m.EditState.taskIndex = m.list.Index()

        switch m.EditState.field {
        case TypeField:
            m.EditState.TempType = item.Type
        case PriorityField:
            m.EditState.TempPriority = item.Priority
        case TaskField:
            ti := textinput.New()
            ti.SetValue(item.Task)
            ti.CursorEnd()
            ti.Width = m.list.Width() - lg.Width(item.Type) - 1 - 3 - 7
            ti.Focus()
            m.EditState.TextInput = ti
        }

        m.EditState.subActive = true
        return m, nil
    }

    // consume all keys in edit mode, don't forward to list navigation
    return m, nil
}
```

---

## Step 6 — Delegate Visual Rendering (`delegate.go`)

The current `case TaskListItem` rendering block in `Render()` applies `editStyle` to whichever field is active in Level 1. We need to extend it to handle Level 2 differently.

Find the section that begins:

```go
if d.editState.active {
    editStyle := ...
    switch d.editState.field {
    case TypeField:
        typStyle = editStyle
    case TaskField:
        taskStyle = editStyle
    case PriorityField:
        priorityStyle = ...
    }
}
```

Replace it with:

```go
if d.editState.active {
    editStyle := segStyle.
        Foreground(styles.PrimaryForeground).
        Background(styles.SelectedBackground)

    if d.editState.subActive {
        // Level 2: show cycling indicators or live text input
        cycleStyle := segStyle.
            Foreground(lg.Color("#ffffff")).
            Background(styles.SelectedBackground)

        switch d.editState.field {
        case TypeField:
            typ = cycleStyle.Render("‹ " + d.editState.TempType + " ›")
        case TaskField:
            task = d.editState.TextInput.View()
        case PriorityField:
            p = d.editState.TempPriority
            if p < 0 || p >= len(priorityColors) {
                p = 0
            }
            priorityStyle = cycleStyle
            priority = priorityStyle.Render(fmt.Sprintf("‹ [%d] ›", p))
        }
    } else {
        // Level 1: highlight the focused field
        switch d.editState.field {
        case TypeField:
            typStyle = editStyle
        case TaskField:
            taskStyle = editStyle
        case PriorityField:
            priorityStyle = priorityStyle.Background(styles.SelectedBackground)
        }
    }
}
```

> **Important:** The `typ`, `task`, and `priority` variables need to be declared before the `if d.editState.active` block so they can be reassigned inside it. Currently the code declares them as individual renders right before building `left`/`right`. Restructure it so the string variables are set first, then the layout is assembled:

```go
// declare segment strings (will be overridden by sub-edit if active)
typ := typStyle.Render(item.Type)
space := segStyle.Render(" ")
task := taskStyle.Render(item.Task)
p := item.Priority
if p < 0 || p >= len(priorityColors) {
    p = 0
}
priority := priorityStyle.Render(fmt.Sprintf("[%d]", p))

// --- apply edit state overrides above ---

left := typ + space + task
right := priority

px := styles.GetPaddingBetween(typ+space+task, priority, m.Width(), contStyle)
content := left + styles.RenderPadding(segStyle, px) + right
```

> The `GetPaddingBetween` call uses the raw string widths. When `TextInput.View()` is the task, its width is already fixed by `ti.Width`, so padding should stay stable.

---

## Step 7 — Objective Panel-Switch Guard (`objective/model.go`)

In `ObjectiveModel.Update`, the `tea.KeyMsg` case checks `LeftFocus` / `RightFocus` at the top of the switch before forwarding to panels. These h/l bindings will intercept the text input if not guarded.

Find:

```go
case tea.KeyMsg:
    switch {
    case key.Matches(msg, m.keys.LeftFocus):
        ...
    case key.Matches(msg, m.keys.RightFocus):
        ...
    }
```

Wrap the switch body:

```go
case tea.KeyMsg:
    if !m.tasks.IsCapturingTextInput() {
        switch {
        case key.Matches(msg, m.keys.LeftFocus):
            ...
        case key.Matches(msg, m.keys.RightFocus):
            ...
        }
    }
```

---

## Verification Checklist

Run with `go run .`, select a milestone, and focus the task panel (`l`):

- [ ] **Level 1 entry:** Press Enter on a task → field highlight appears, h/l moves between Type / Task / Priority
- [ ] **TypeField Level 2:** Navigate to Type, press Enter → `‹ feat ›` renders; h/l cycles through options
- [ ] **TypeField commit:** Press Enter or Esc → type is saved on the row, back to Level 1 field select
- [ ] **PriorityField Level 2:** Navigate to Priority, press Enter → `‹ [3] ›` renders; h/l cycles 0–5
- [ ] **PriorityField commit:** Enter or Esc saves, back to Level 1
- [ ] **TaskField Level 2:** Navigate to Task, press Enter → inline text input appears with existing title pre-filled
- [ ] **Text input isolation:** Press h or l → characters appear in the input, panels do NOT switch
- [ ] **TaskField commit:** Press Enter or Esc → updated title shown on row, back to Level 1
- [ ] **Level 1 exit:** Press Esc from Level 1 → back to normal list navigation
