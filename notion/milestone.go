package notion

type MilestonePagesMsg struct {
	Pages []MilestonePage
	Err   error
}

type MilestoneIDsMsg struct {
	IDs []string
	Err error
}

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
