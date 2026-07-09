package notion

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (c *Client) QueryNotePages(projID, cursor string) tea.Cmd {
	fprops := []string{
		notePropTitle,
		notePropCreatedLabel,
	}

	return func() tea.Msg {
		body := queryNoteBody(projID, 100)
		res, err := queryDatasource[NotePage](c, c.notesDatasourceID, body, cursor, fprops)
		if err != nil {
			return QueryNotePagesMsg{Err: err}
		}
		var nextCursor *string
		if res.HasMore {
			nextCursor = res.NextCursor
		}
		return QueryNotePagesMsg{Pages: res.Results, NextCursor: nextCursor}
	}
}
