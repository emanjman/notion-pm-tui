package milestonelist

func mockMilestoneItems() []MilestoneListItem {
	return []MilestoneListItem{
		MilestoneListItem{
			ID:       "1",
			Name:     "Setup Project Structure",
			Status:   "🎉 complete",
			Progress: 1.0,
			Tags:     []string{"backend", "setup"},
		},
		MilestoneListItem{
			ID:       "2",
			Name:     "Implement Notion API",
			Status:   "🚧 under development",
			Progress: 0.75,
			Tags:     []string{"backend", "api"},
		},
		MilestoneListItem{
			ID:       "3",
			Name:     "Build TUI Dashboard",
			Status:   "🚧 under development",
			Progress: 0.4,
			Tags:     []string{"frontend", "tui"},
		},
		MilestoneListItem{
			ID:       "4",
			Name:     "Authentication System",
			Status:   "😴 idle",
			Progress: 0.0,
			Tags:     []string{"backend", "auth"},
		},
		MilestoneListItem{
			ID:       "5",
			Name:     "Data Persistence Layer",
			Status:   "😴 idle",
			Progress: 0.0,
			Tags:     []string{"backend", "database"},
		},
		MilestoneListItem{
			ID:       "6",
			Name:     "Testing & QA",
			Status:   "😴 idle",
			Progress: 0.0,
			Tags:     []string{"testing"},
		},
		MilestoneListItem{
			ID:       "7",
			Name:     "Documentation",
			Status:   "🚧 under development",
			Progress: 0.2,
			Tags:     []string{"docs"},
		},
	}
}
