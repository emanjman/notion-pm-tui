package task

import (
	"notion-project-tui/notion"
	"strconv"
)

type LoadMoreItem struct {
	Status  notion.TaskStatus
	Loading bool
}

func (l LoadMoreItem) FilterValue() string { return "" }

type GroupHeader struct {
	Status  notion.TaskStatus
	Hidden  bool
	Count   int
	HasMore bool
}

func (h GroupHeader) FilterValue() string { return "" }

type Item struct {
	ID       string
	Task     string
	Status   notion.TaskStatus
	Priority int
	Type     string
}

func (t Item) FilterValue() string { return t.Task + "_" + t.Type }

func NewItem(page notion.TaskPage) Item {
	statusString := page.Properties.Status.Status.Name
	status, _ := notion.TaskStatusFromString(statusString)

	titleProp := page.Properties.Title.Title
	title := notion.ExtractPlainText(titleProp)

	t := Item{
		ID:     page.ID,
		Task:   title,
		Status: status,
		Type:   page.Properties.Type.Select.Name,
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
