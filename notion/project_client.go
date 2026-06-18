package notion

import (
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// cmd func returns a tea.Msg
func (c *Client) FetchProject() tea.Cmd {
	return func() tea.Msg {
		start := time.Now()

		url := c.baseURL + "/pages/" + c.projID

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return ProjectMsg{Err: err, Duration: time.Since(start)}
		}

		// parse as json
		var proj ProjectPage
		if err := c.do(req, &proj); err != nil {
			return ProjectMsg{Err: err, Duration: time.Since(start)}
		}
		return ProjectMsg{Data: proj, Duration: time.Since(start)}
	}
}
