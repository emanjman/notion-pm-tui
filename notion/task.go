package notion

type TaskTypeOptionsMsg struct {
	Options []SelectItem
	Err     error
}

type TaskDatasourceResponse struct {
	Properties struct {
		Type struct {
			Select struct {
				Options []SelectItem `json:"options"`
			} `json:"select"`
		} `json:"type"`
	} `json:"properties"`
}

type TaskGroup struct {
	Tasks      []TaskPage
	NextCursor *string
	Hide       bool
	Loading    bool
}

type TaskGroups map[TaskStatus]TaskGroup

type ToggleTaskGroupMsg struct {
	Status TaskStatus
}

type QueryMoreTaskPagesMsg struct {
	Status TaskStatus
}

type QueryTaskPagesMsg struct {
	Pages        []TaskPage
	NextCursor   *string
	Status       TaskStatus
	MilestoneIdx int
	Err          error
}

// result of notion task-page creation. TempID is echoed back so the caller can
// swap the optimistic local item's id for the real one.
type AddTaskPageMsg struct {
	TempID string
	Page   *TaskPage
	Err    error
}

type TaskPage struct {
	ID         string         `json:"id"`
	Properties TaskProperties `json:"properties"`
}

type TaskProperties struct {
	Title     TitleProperty    `json:"task"`
	Status    StatusProperty   `json:"status"`
	Priority  SelectProperty   `json:"priority"`
	Type      SelectProperty   `json:"type"`
	Milestone RelationProperty `json:"@milestone"`
}
