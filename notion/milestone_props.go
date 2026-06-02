package notion

type MilestoneProperties struct {
	Title     TitleProperty   `json:"name"`
	Progress  FormulaProperty `json:"$progress"` // type:number
	Status    FormulaProperty `json:"$status"`   // type:string
	TaskCount FormulaProperty `json:"$task-ct"`  // type:number
}

var (
	milestonePropTitle                 = "name"
	milestonePropProgress              = "$progress"
	milestonePropVersionRelation       = "@version"
	milestonePropProjectRollupRelation = "r/@project"
	milestonePropStatusFormula         = "$status"
	milestonePropLatestActivityAt      = "$latest-activity-at"
	milestonePropTaskCount             = "$task-ct"
)
