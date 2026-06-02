package notion

import (
	"bytes"
	"encoding/json"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

func (c *Client) QueryMilestones(projID string, status MilestoneStatus, cursor string) tea.Cmd {
	fprops := []string{
		milestonePropTitle,
		milestonePropProgress,
		milestonePropStatusFormula,
		milestonePropTaskCount,
	}

	return func() tea.Msg {
		body := queryMilestoneBody(projID, status, 5)
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

type AddMilestonePageMsg struct {
	Err    error
	TempID string         // optimistic temp id to reconcile against the created page
	Page   *MilestonePage // created page (carries the real notion id)
}

// AddMilestone creates a milestone page on notion. TempID is echoed back on the
// msg so the caller can swap the optimistic local item's id for the real one.
func (c *Client) AddMilestone(tempID, title string) tea.Cmd {
	return func() tea.Msg {
		endpt := c.baseURL + "/pages"
		body := addMilestoneBody(c.milestoneDsId, title)
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
