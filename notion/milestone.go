package notion

import "fmt"

type MilestonePagesMsg struct {
	Pages      []MilestonePage
	NextCursor *string         // bookmark for subsequent milestsone pages
	Status     MilestoneStatus // grouping key
	Err        error
}

type FetchMoreMilestonesMsg struct {
	Status MilestoneStatus
}

type MilestoneGroup struct {
	Milestones []MilestonePage
	NextCursor *string
	Hide       bool
	Loading    bool
}

type MilestoneGroups map[MilestoneStatus]MilestoneGroup

// ------

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

// ---

type MilestoneStatus int

const (
	MilestoneUnderDevelopment MilestoneStatus = iota
	MilestoneIdle
	MilestoneComplete
	_MilestoneStatusCount // sentinel for array
)

// get milestone status in an ordered arr
func MilestoneStatusOrder() []MilestoneStatus {
	statuses := make([]MilestoneStatus, _MilestoneStatusCount)
	for i := range statuses {
		statuses[i] = MilestoneStatus(i)
	}
	return statuses
}

// map each status to a neat string (indexed by enum val); use static/fixed arr
func (state MilestoneStatus) String() string {
	return [...]string{
		"🚧 under development",
		"😴 idle",
		"🎉 complete",
	}[state]
}

// map dedicated status string back to an enum
func MilestoneStatusFromString(s string) (MilestoneStatus, error) {
	if status, ok := milestoneStatusByString[s]; ok {
		return status, nil
	}
	return 0, fmt.Errorf("Unknown milestone status: %q", s)
}

var milestoneStatusByString = func() map[string]MilestoneStatus {
	m := make(map[string]MilestoneStatus, _MilestoneStatusCount)
	for i := range _MilestoneStatusCount {
		s := MilestoneStatus(i)
		m[s.String()] = s
	}
	return m
}()
