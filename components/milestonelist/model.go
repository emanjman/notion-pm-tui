package milestonelist

import (
	"notion-project-tui/notion"
	listutil "notion-project-tui/util/list"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type MilestoneListModel struct {
	list    list.Model
	loading bool
	groups  map[string][]MilestoneListItem // header to items
	hidden  map[string]bool                // hidden group
	Keys    KeyMap
}

var statusOrder = []string{"🚧 under development", "😴 idle", "🎉 complete"}

func NewMilestoneListModel() MilestoneListModel {
	l := list.New([]list.Item{}, NewMilestoneListDelegate(), 0, 0)
	l.Title = "Milestones"
	l.SetShowHelp(false)

	m := MilestoneListModel{
		list:    l,
		loading: true,
		groups:  listutil.GroupByKey(mockMilestoneItems()),
		hidden:  map[string]bool{},
		Keys:    DefaultKeyMap,
	}
	m.list.SetItems(listutil.BuildGroupList(m.groups, m.hidden, statusOrder))

	return m
}

type MilestoneSelectedMsg struct {
	ID string
}

// just forward the list.Update(msg)
// and forward its returned response
func (m MilestoneListModel) Update(msg tea.Msg) (MilestoneListModel, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Select):
			selected := m.list.SelectedItem()

			if header, ok := selected.(listutil.ListItemGroupHeader); ok {
				m.hidden[header.Label] = !m.hidden[header.Label]                          // toggle
				m.list.SetItems(listutil.BuildGroupList(m.groups, m.hidden, statusOrder)) // rebuild list
				return m, nil
			} else {
				return m, func() tea.Msg {
					return MilestoneSelectedMsg{ID: m.SelectedMilestoneId()}
				}
			}
		}

	case notion.MilestoneMsg:
		if msg.Err != nil {
			return m, nil
		}

		// create the list items
		tempItems := make([]MilestoneListItem, len(msg.Data))
		for i, page := range msg.Data {
			tempItems[i] = NewMilestoneListItem(page)
		}

		m.groups = listutil.GroupByKey(tempItems)
		items := listutil.BuildGroupList(m.groups, m.hidden, statusOrder)

		m.list.SetItems(items)
		m.loading = false

	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg) // handles up/down nav
	return m, cmd
}

// just forward the list.View()
func (m MilestoneListModel) View() string {
	// ! temp, styling
	// if m.loading {
	// 	return "Loading milestones..."
	// }
	return m.list.View()
}

func (m MilestoneListModel) SelectedMilestoneId() string {
	item := m.list.SelectedItem()

	switch item := item.(type) {
	// get first milestone id of this group
	case listutil.ListItemGroupHeader:
		return m.groups[item.Label][0].ID
	// otherwise, on milestone, return its id
	case MilestoneListItem:
		return item.ID
	}

	return ""
}
