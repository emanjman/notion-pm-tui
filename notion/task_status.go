package notion

import "fmt"

type TaskStatus int

const (
	TaskDev TaskStatus = iota
	TaskIdle
	TaskDone
	TaskArchive
	_TaskStatusCount // sentinel for array
)

// -- helpers --

// get task status as an ordered arr
func TaskStatusOrder() []TaskStatus {
	statuses := make([]TaskStatus, _TaskStatusCount)
	for i := range statuses {
		statuses[i] = TaskStatus(i)
	}
	return statuses
}

// map each status to a str idx'd by enum val
func (status TaskStatus) String() string {
	return [...]string{
		"dev",
		"idle",
		"done",
		"archive",
	}[status]
}

// map dedicated status string back to an enum
func TaskStatusFromString(s string) (TaskStatus, error) {
	if status, ok := taskStatusByString[s]; ok {
		return status, nil
	}
	return 0, fmt.Errorf("unknown task status: %q", s)
}

var taskStatusByString = func() map[string]TaskStatus {
	m := make(map[string]TaskStatus, _TaskStatusCount)
	for i := range _TaskStatusCount {
		s := TaskStatus(i)
		m[s.String()] = s
	}
	return m
}()
