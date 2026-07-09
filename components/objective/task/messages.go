package task

type UpdateTitleMsg struct{ Err error }

type UpdateSelectionsMsg struct{ Err error }

// result of a notion task status update. carries the task id + its prior status
// so the optimistic group-move can be reverted on failure.
type UpdateStatusMsg struct {
	TaskID     string
	PrevStatus string
	Err        error
}

// result of a notion task-page trash. carries the deleted task + its prior list
// index so the optimistic deletion can be reverted on failure.
type DeleteTaskMsg struct {
	Task Item
	Idx  int
	Err  error
}
