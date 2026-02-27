package tasklist

import (
	"notion-project-tui/notion"
	"strconv"
)

type TaskListItem struct {
	ID          string
	Task        string
	Status      string
	Priority    int
	Type        string
	MilestoneID string
}

func (t TaskListItem) FilterValue() string { return t.Task + "_" + t.Type }
func (t TaskListItem) GroupKey() string    { return t.Status }

func NewTaskListItem(page notion.TaskPage) TaskListItem {
	t := TaskListItem{
		ID:          page.ID,
		Task:        notion.ExtractPlainText(page.Properties.Title.Title),
		Status:      page.Properties.Status.Status.Name,
		Type:        page.Properties.Type.Select.Name,
		MilestoneID: page.Properties.Milestone.Relation[0].ID,
	}

	// handle priority int conversion
	p := page.Properties.Priority.Select.Name
	priority, err := strconv.Atoi(p)
	if err != nil {
		t.Priority = -1
	}
	t.Priority = priority

	return t
}
