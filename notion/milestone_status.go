package notion

import "fmt"

// -- enum --

type MilestoneStatus int

const (
	MilestoneUnderDevelopment MilestoneStatus = iota
	MilestoneIdle
	MilestoneComplete
	_MilestoneStatusCount // sentinel for array
)

// -- msg --

// load more milestones of this milestone-status-group
type FetchMoreMilestonesMsg struct {
	Status MilestoneStatus
}

// -- helpers --

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
