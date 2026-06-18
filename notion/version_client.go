package notion

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (c *Client) QueryVersionPages(projID, cursor string) tea.Cmd {
	fprops := []string{
		versionPropTitle,
		versionPropCreatedAt,
		versionPropMilestonesRelation,
		versionPropProjectRelation,
	}

	return func() tea.Msg {
		body := queryVersionBody(projID)
		res, err := queryDatasource[VersionPage](c, c.versionsDatasourceID, body, cursor, fprops)
		if err != nil {
			return QueryVersionPagesMsg{Err: err}
		}

		var nextCursor *string
		if res.HasMore {
			nextCursor = res.NextCursor
		}
		return QueryVersionPagesMsg{Pages: res.Results, NextCursor: nextCursor}
	}
}
