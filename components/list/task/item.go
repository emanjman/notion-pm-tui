package task

import (
	"notion-project-tui/notion"
	"strconv"
)

type TaskItem struct {
	ID          string
	Task        string
	Status      string
	Priority    int
	Type        string
	MilestoneID string
}

func (t TaskItem) FilterValue() string { return t.Task + "_" + t.Type }
func (t TaskItem) GroupKey() string    { return t.Status }

func NewTaskItem(page notion.TaskPage) TaskItem {
	t := TaskItem{
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
