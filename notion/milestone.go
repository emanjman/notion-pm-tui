package notion

type MilestonePage struct {
	ID         string              `json:"id"`
	Properties MilestoneProperties `json:"properties"`
	Icon       *Icon               `json:"icon"`
}

type MilestoneProperties struct {
	Title     TitleProperty   `json:"name"`
	Progress  FormulaProperty `json:"progress"` // type:number
	Status    FormulaProperty `json:"$status"`  // type:string
	TaskCount FormulaProperty `json:"task-ct"`  // type:number
}

type MilestonePagesMsg struct {
	Pages      []MilestonePage
	NextCursor *string         // bookmark for subsequent milestsone pages
	Status     MilestoneStatus // grouping key
	Err        error
}

type MilestoneGroup struct {
	Milestones []MilestonePage
	NextCursor *string
	Hide       bool
	Loading    bool
}

type MilestoneGroups map[MilestoneStatus]MilestoneGroup
