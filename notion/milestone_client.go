package notion

import (
	"bytes"
	"encoding/json"
	"log"
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
	Err  error
	Page *MilestonePage
}

// todo: needs fixing
func (c *Client) AddMilestone(tempPage MilestonePage) tea.Cmd {
	endpt := c.baseURL + "/pages"
	body := addMilestoneBody(c.milestoneDsId, tempPage)
	b, err := json.Marshal(body)
	if err != nil {
		return AddMilestonePageMsg{Err: err, Page: nil}
	}

	req, err := http.NewRequest("POST", endpt, bytes.NewReader(b))
	if err != nil {
		return AddMilestonePageMsg{Err: err, Page: nil}
	}

	req.Header.Add("Content-Type", "application/json")
	var res PaginationResponse[T]
	if err := c.do(req, &res); err != nil {
		return AddMilestonePageMsg{Err: err, Page: nil}
	}

	log.Printf("res:", res)
	return AddMilestonePageMsg{Err: nil, Page: res}
}
