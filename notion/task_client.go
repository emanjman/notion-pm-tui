package notion

import (
	"bytes"
	"encoding/json"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

func (c *Client) QueryTasks(milestoneID, status, cursor string, milestoneIdx int) tea.Cmd {
	fprops := []string{
		taskPropTitle,
		taskPropTypeSelect,
		taskPropPriority,
		taskPropStatus,
	}

	return func() tea.Msg {
		body := taskQueryBody(milestoneID, status, 5)
		res, err := queryDatasource[TaskPage](c, c.tasksDatasourceID, body, cursor, fprops)
		if err != nil {
			return QueryTaskPagesMsg{Err: err, Status: status, MilestoneIdx: milestoneIdx}
		}
		var nextCursor *string
		if res.HasMore {
			nextCursor = res.NextCursor
		}
		return QueryTaskPagesMsg{
			Pages:        res.Results,
			NextCursor:   nextCursor,
			Status:       status,
			MilestoneIdx: milestoneIdx,
		}
	}
}

// AddTask creates a task page on notion under the given milestone. TempID is
// echoed back on the msg so the caller can swap the optimistic local item's id
// for the real one.
func (c *Client) AddTask(tempID, title, milestoneID, status, taskType string, priority int) tea.Cmd {
	return func() tea.Msg {
		endpt := c.baseURL + "/pages"
		body := addTaskBody(c.tasksDatasourceID, title, milestoneID, status, taskType, priority)
		b, err := json.Marshal(body)
		if err != nil {
			return AddTaskPageMsg{Err: err, TempID: tempID}
		}

		req, err := http.NewRequest("POST", endpt, bytes.NewReader(b))
		if err != nil {
			return AddTaskPageMsg{Err: err, TempID: tempID}
		}
		req.Header.Add("Content-Type", "application/json")

		// POST /pages returns a single page object, not a paginated list
		var res TaskPage
		if err := c.do(req, &res); err != nil {
			return AddTaskPageMsg{Err: err, TempID: tempID}
		}
		return AddTaskPageMsg{TempID: tempID, Page: &res}
	}
}

// options for `type` property, e.g. "style" "feat" "refactor"
func (c *Client) FetchTaskTypeOptions() tea.Cmd {
	return func() tea.Msg {
		url := c.baseURL + "/data_sources/" + c.tasksDatasourceID
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return TaskTypeOptionsMsg{Err: err}
		}
		var res TaskDatasourceResponse
		if err := c.do(req, &res); err != nil {
			return TaskTypeOptionsMsg{Err: err}
		}
		return TaskTypeOptionsMsg{Options: res.Properties.Type.Select.Options}
	}
}
