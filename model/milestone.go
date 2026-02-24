package model

import (
	"fmt"
	"io"
	"notion-project-tui/notion"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)

	case notion.MilestoneMsg:
		if msg.Err != nil {
			return m, nil
		}

		// create the list items
		items := make([]list.Item, len(msg.Data))
		for i, page := range msg.Data {
			items[i] = NewMilestoneListItem(page)
		}

		m.list.SetItems(items)
		m.loading = false

	}

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

func NewMilestoneListItem(page notion.MilestonePage) MilestoneListItem {
	title := notion.ExtractPlainText(page.Properties.Title.Title)

	status := ""
	if page.Properties.Status.Formula.String != nil {
		status = *page.Properties.Status.Formula.String
	}

	progress := 0.0
	if page.Properties.Progress.Formula.Number != nil {
		progress = *page.Properties.Progress.Formula.Number
	}

	tags := make([]string, len(page.Properties.Tags.MultiSelect))
	for i, tag := range page.Properties.Tags.MultiSelect {
		tags[i] = tag.Name
	}

	return MilestoneListItem{
		ID:       page.ID,
		Name:     title,
		Status:   status,
		Progress: progress,
		Tags:     tags,
	}
}

// -------------------------------------------

// implementation for delegate
type MilestoneDelegate struct{}

func (d MilestoneDelegate) Height() int  { return 1 }
func (d MilestoneDelegate) Spacing() int { return 0 }
func (d MilestoneDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

// required delegate function, where `index` holds the hovering item
func (d MilestoneDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	milestone := item.(MilestoneListItem)
	line := milestone.Name

	// where `index` is the curr item's index, vs `m.Index()` is the item the user is hovering over
	if index == m.Index() {
		line = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Render("> " + line)
	} else {
		line = "  " + line
	}

	// write to `w`
	fmt.Fprintf(w, "%s", line)
}
