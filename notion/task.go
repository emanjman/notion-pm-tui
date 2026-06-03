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

// keyed by status: "dev", "idle", "done", "archive"
type TaskGroups map[string]TaskGroup

type ToggleTaskGroupMsg struct {
	Status string
}

type QueryMoreTaskPagesMsg struct {
	Status string
}

type QueryTaskPagesMsg struct {
	Pages        []TaskPage
	NextCursor   *string
	Status       string
	MilestoneIdx int
	Err          error
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
