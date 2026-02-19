package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	items    []string
	index    int
	selected map[int]struct{} // use `struct{}` as "exists" for memory efficiency
}

// init state
func initialModel() model {
	return model{
		items:    []string{"item1", "item2", "item3"},
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	// no i/o atm
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// handle based on the msg type
	switch msg := msg.(type) {

	case tea.KeyMsg:
		// handle the type of key
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		// navigation
		case "up", "k":
			if m.index > 0 {
				m.index--
			}
		case "down", "j":
			if m.index < len(m.items)-1 {
				m.index++
			}

		//
		case "enter":
			if _, ok := m.selected[m.index]; ok {
				delete(m.selected, m.index) // unselect
			} else {
				m.selected[m.index] = struct{}{} // create instance of empty struct
			}

		}
	}

	return m, nil
}

func (m model) View() string {
	view := "What should we buy?\n\n"

	// render each item
	for i, item := range m.items {
		// display cursor
		cursor := " "
		if i == m.index {
			cursor = ">"
		}

		// display check
		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		// render
		view += fmt.Sprintf("%s [%s] %s\n", cursor, checked, item)
	}

	view += "\nPress q to quit\n"

	return view
}

func main() {
	app := tea.NewProgram(initialModel())

	if _, err := app.Run(); err != nil {
		fmt.Printf("Program error: %v", err)
		os.Exit(1)
	}
}
