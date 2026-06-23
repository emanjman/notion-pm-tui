package notion

import tea "github.com/charmbracelet/bubbletea"

func (c *Client) QueryProjectPages(cursor string) tea.Cmd {
	fprops := []string{
		projectPropTitle,
	}

	return func() tea.Msg {
		body := queryProjectBody()
		res, err := queryDatasource[ProjectPage](c, c.projectsDatasourceID, body, cursor, fprops)
		if err != nil {
			return QueryProjectPagesMsg{Err: err}
		}
		var nextCursor *string
		if res.HasMore {
			nextCursor = res.NextCursor
		}

		return QueryProjectPagesMsg{Pages: res.Results, NextCursor: nextCursor}
	}
}
