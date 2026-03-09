package listutil

import "github.com/charmbracelet/bubbles/list"

type ListItemGroupHeader struct {
	Label  string
	Hidden bool
	Count  int
}

func (h ListItemGroupHeader) FilterValue() string { return "" }

// -------

// to conform, implement list.Item + GroupKey()
type Groupable interface {
	list.Item
	GroupKey() string // value to group by (e.g. status)
}

// maps each item by their group-key
func GroupByKey[T Groupable](items []T) map[string][]T {
	groups := map[string][]T{}

	for _, item := range items {
		key := item.GroupKey()
		groups[key] = append(groups[key], item)
	}

	return groups
}

func BuildGroupList[T Groupable](
	groups map[string][]T,
	hidden map[string]bool,
	order []string, // groups ordered
) []list.Item {
	var items []list.Item

	// build list using the group order
	for _, key := range order {
		group, ok := groups[key]
		if !ok || len(group) == 0 {
			continue
		}

		// add group header
		items = append(items, ListItemGroupHeader{
			Label:  key,
			Hidden: hidden[key],
			Count:  len(group),
		})

		// add the group's children if not hidden
		if !hidden[key] {
			for _, item := range group {
				items = append(items, item)
			}
		}
	}

	return items
}
