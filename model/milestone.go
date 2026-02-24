package model

import (
	"fmt"
	"io"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type MilestonesModel struct {
	list    list.Model
	loading bool
}

func NewMilestonesModel() MilestonesModel {
	delegate := MilestoneDelegate{}

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "Milestones"

	return MilestonesModel{
		list:    l,
		loading: true,
	}
}

// just forward the list.Update(msg)
// and forward its returned response
func (m MilestonesModel) Update(msg tea.Msg) (MilestonesModel, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// just forward the list.View()
func (m MilestonesModel) View() string {
	if m.loading {
		return "Loading milestones..."
	}
	return m.list.View()
}

// --------------------------------------------

// implementation for the `list.Item` interface
type MilestoneListItem struct {
	ID           string
	Name         string
	Status       string
	LastActivity time.Time
	Progress     float64
	Tags         []string
}

// func (m MilestoneListItem) Title() string       { return m.Name }
// func (m MilestoneListItem) Description() string { return m.Status }
func (m MilestoneListItem) FilterValue() string { return m.Name }

// -------------------------------------------

// implementation for delegate
type MilestoneDelegate struct{}

func (d MilestoneDelegate) Height() int  { return 1 }
func (d MilestoneDelegate) Spacing() int { return 0 }
func (d MilestoneDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}
func (d MilestoneDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	milestone := item.(MilestoneListItem)

	// write to `w`
	fmt.Fprintf(w, "%s\n", milestone.Name)
}
