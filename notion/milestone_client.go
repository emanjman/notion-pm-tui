package notion

import (
	"bytes"
	"encoding/json"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

// todo: need to wire in getting milestone pages by versionID, not projID
func (c *Client) QueryMilestonePages(projID string, status MilestoneStatus, cursor string) tea.Cmd {
	fprops := []string{
		milestonePropTitle,
		milestonePropProgress,
		milestonePropStatusFormula,
		milestonePropTaskCount,
	}

	return func() tea.Msg {
		body := queryMilestoneBody(projID, status, 5)
		res, err := queryDatasource[MilestonePage](c, c.milestonesDatasourceID, body, cursor, fprops)
		if err != nil {
			return QueryMilestonePagesMsg{Err: err, Status: status}
		}
		var nextCursor *string
		if res.HasMore {
			nextCursor = res.NextCursor
		}
		return QueryMilestonePagesMsg{Pages: res.Results, NextCursor: nextCursor, Status: status}
	}
}

// AddMilestone creates a milestone page on notion. TempID is echoed back on the
// msg so the caller can swap the optimistic local item's id for the real one.
func (c *Client) AddMilestone(tempID, title string) tea.Cmd {
	// todo: wire the @version datasource so the version can be selected per-project.
	// for now the demo project ("Hoop Archives") has a single version, so we hardcode
	// its page id. a milestone must hang off a @version (the @project is a rollup
	// through it), otherwise it won't roll up to any project.
	const demoVersionPageID = "346b7273-944b-80ee-bc8d-e9ead7e1e623"

	return func() tea.Msg {
		endpt := c.baseURL + "/pages"
		body := addMilestoneBody(c.milestonesDatasourceID, title, demoVersionPageID)
		b, err := json.Marshal(body)
		if err != nil {
			return AddMilestonePageMsg{Err: err, TempID: tempID}
		}

		req, err := http.NewRequest("POST", endpt, bytes.NewReader(b))
		if err != nil {
			return AddMilestonePageMsg{Err: err, TempID: tempID}
		}
		req.Header.Add("Content-Type", "application/json")

		// POST /pages returns a single page object, not a paginated list
		var res MilestonePage
		if err := c.do(req, &res); err != nil {
			return AddMilestonePageMsg{Err: err, TempID: tempID}
		}
		return AddMilestonePageMsg{TempID: tempID, Page: &res}
	}
}
