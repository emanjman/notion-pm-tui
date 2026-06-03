package notion

import (
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
		res, err := queryDatasource[TaskPage](c, c.tasksDsId, body, cursor, fprops)
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

// options for `type` property, e.g. "style" "feat" "refactor"
func (c *Client) FetchTaskTypeOptions() tea.Cmd {
	return func() tea.Msg {
		url := c.baseURL + "/data_sources/" + c.tasksDsId
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
