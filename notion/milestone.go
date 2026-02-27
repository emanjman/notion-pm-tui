package notion

import "time"

type MilestoneMsg struct {
	Data     []MilestonePage
	Err      error
	Duration time.Duration
}

// ------

type MilestonePage struct {
	ID         string              `json:"id"`
	Properties MilestoneProperties `json:"properties"`
}

type MilestoneProperties struct {
	Title            TitleProperty       `json:"name"`
	Tags             MultiSelectProperty `json:"tags"`
	Progress         FormulaProperty     `json:"progress"`            // type:number
	Status           FormulaProperty     `json:"$status"`             // type:string
	LatestActivityAt FormulaProperty     `json:"$latest-activity-at"` // type:date
	Tasks            RelationProperty    `json:"@tasks"`
}

type SelectedMilestone struct {
	ID          string
	TasksPropID string
}

type MilestoneSelectedMsg struct {
	Milestone SelectedMilestone
}
