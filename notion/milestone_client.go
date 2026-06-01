package notion

import tea "github.com/charmbracelet/bubbletea"

func (c *Client) QueryMilestones(projID string, status MilestoneStatus, cursor string) tea.Cmd {
	return func() tea.Msg {
		body := milestoneQueryBody(projID, status, 5)
		fprops := []string{"name", "progress", "$status", "task-ct"}
		res, err := queryDatasource[MilestonePage](c, c.milestoneDsId, body, cursor, fprops)
		if err != nil {
			return MilestonePagesMsg{Err: err, Status: status}
		}
		var nextCursor *string
		if res.HasMore {
			nextCursor = res.NextCursor
		}
		return MilestonePagesMsg{Pages: res.Results, NextCursor: nextCursor, Status: status}
	}
}
