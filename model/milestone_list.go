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

type MilestoneListModel struct {
	list    list.Model
	loading bool
}

func NewMilestoneListModel() MilestoneListModel {
	l := list.New([]list.Item{}, MilestoneListDelegate{}, 0, 0)

	// custom configs
	l.Title = "Milestones"
	l.SetShowHelp(false)

	return MilestoneListModel{list: l, loading: true}
}

// just forward the list.Update(msg)
// and forward its returned response
func (m MilestoneListModel) Update(msg tea.Msg) (MilestoneListModel, tea.Cmd) {
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
func (m MilestoneListModel) View() string {
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
type MilestoneListDelegate struct {
	defaultStyle  lipgloss.Style
	selectedStyle lipgloss.Style
}

func (d MilestoneListDelegate) Height() int  { return 1 }
func (d MilestoneListDelegate) Spacing() int { return 0 }
func (d MilestoneListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

// required delegate function, where `index` holds the hovering item
func (d MilestoneListDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	milestone := item.(MilestoneListItem)

	// where `index` is the curr item's index, vs `m.Index()` is the item the user is hovering over
	style := d.defaultStyle
	if index == m.Index() {
		style = d.selectedStyle
	}

	// write to `w`
	fmt.Fprintf(w, style.Render(milestone.Name))
}

func NewMilestoneListDelegate() MilestoneListDelegate {
	return MilestoneListDelegate{
		defaultStyle:  lipgloss.NewStyle(),
		selectedStyle: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")),
	}
}
