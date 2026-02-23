package model

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
)

type MilestonesModel struct {
	list    list.Model
	loading bool
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

func (m MilestoneListItem) Title() string       { return m.Name }
func (m MilestoneListItem) Description() string { return m.Status }
func (m MilestoneListItem) FilterValue() string { return m.Name }

// -------------------------------------------
