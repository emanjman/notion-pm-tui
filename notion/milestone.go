package notion

type MilestonePagesMsg struct {
	Pages      []MilestonePage
	NextCursor *string
	Status     string
	Err        error
}

type FetchMoreMilestonesMsg struct {
	Status string
}

type MilestoneGroup struct {
	Milestones []MilestonePage
	NextCursor *string
	Hide       bool
	Loading    bool
}

// keyed by status: "🚧 under development", "😴 idle", "🎉 complete"
type MilestoneGroups map[string]MilestoneGroup

// ------

type MilestonePage struct {
	ID         string              `json:"id"`
	Properties MilestoneProperties `json:"properties"`
	Icon       *Icon               `json:"icon"`
}

type MilestoneProperties struct {
	Title    TitleProperty   `json:"name"`
	Progress FormulaProperty `json:"progress"` // type:number
	Status   FormulaProperty `json:"$status"`  // type:string
}
