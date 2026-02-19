package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const url = "https://charm.sh/"

// -----------------------------------------

type errMsg struct {
	err error
}

// get the internal error message (so that we don't have to reach in)
func (e errMsg) Error() string {
	return e.err.Error()
}

type statusMsg int

// -----------------------------------------

type model struct {
	status int
	err    error
}

// this is some tea.Cmd
func checkServer() tea.Msg {
	c := &http.Client{Timeout: 10 * time.Second}
	res, err := c.Get(url)

	if err != nil {
		return errMsg{err}
	}

	// server response
	return statusMsg(res.StatusCode)
}

// return the function definition
func (m model) Init() tea.Cmd {
	return checkServer
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// `msg` is assigned to narrow it to represent the case-by-case type
	// i.e. `msg` narrows to `statusMsg`, `errMsg`, `tea.KeyMsg` per case
	switch msg := msg.(type) {

	case statusMsg:
		m.status = int(msg)
		return m, tea.Quit

	case errMsg:
		m.err = msg
		return m, tea.Quit

	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	}

	// tell user we're doing something
	view := fmt.Sprintf("Checking %s ... ", url)

	if m.status > 0 {
		view += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
	}

	return "\n" + view + "\n\n"
}

func main() {
	app := tea.NewProgram(model{})
	if _, err := app.Run(); err != nil {
		fmt.Printf("Progrma error: %v\n", err)
		os.Exit(1)
	}
}
