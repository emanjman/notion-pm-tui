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
}

type MilestoneProperties struct {
	Title               TitleProperty   `json:"name"`
	Tags                SelectProperty  `json:"tags"`
	Progress            FormulaProperty `json:"progress"`               // type:number
	Status              FormulaProperty `json:"$status"`                // type:string
	LatestActivityLabel FormulaProperty `json:"$latest-acitivty-label"` // type:string
	// LatestActivityAt    FormulaProperty     `json:"$latest-activity-at"`    // type:date
	Tasks RelationProperty `json:"@tasks"`
}
