package milestonelist

import (
	"notion-project-tui/notion"
	"time"

	"github.com/charmbracelet/bubbles/list"
)

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

type MilestoneListItemHeader struct {
	Label  string
	Hidden bool
	Count  int
}

// exclude header in filter-search (required item function)
func (g MilestoneListItemHeader) FilterValue() string { return "" }

// helper func to group items into a map (keyed by their status)
func groupByStatus(items []MilestoneListItem) map[string][]MilestoneListItem {

	groups := map[string][]MilestoneListItem{}
	for _, item := range items {
		groups[item.Status] = append(groups[item.Status], item)
	}
	return groups
}

var statusOrder = []string{"🚧 under development", "😴 idle", "🎉 complete"}

// conforms the group-map into a []list.Item + add headers
func (m MilestoneListModel) buildGroupedList() []list.Item {
	var items []list.Item

	// build list in this group order
	for _, status := range statusOrder {
		group, ok := m.grouped[status]
		if !ok {
			continue
		}

		// add group header
		items = append(items, MilestoneListItemHeader{
			Label:  status,
			Hidden: m.hidden[status],
			Count:  len(group),
		})

		// add the group's items (if not hidden)
		if !m.hidden[status] {
			for _, milestone := range group {
				items = append(items, milestone)
			}
		}
	}

	return items
}
