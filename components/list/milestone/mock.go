package milestone

func mockMilestoneItems() []MilestoneListItem {
	return []MilestoneListItem{
		MilestoneListItem{
			ID:                  "1",
			Name:                "Setup Project Structure",
			Status:              "🎉 complete",
			Progress:            1.0,
			Tag:                 "backend",
			LatestActivityLabel: "12/8/25",
		},
		MilestoneListItem{
			ID:                  "2",
			Name:                "Implement Notion API",
			Status:              "🚧 under development",
			Progress:            0.75,
			Tag:                 "backend",
			LatestActivityLabel: "today",
		},
		MilestoneListItem{
			ID:                  "3",
			Name:                "Build TUI Dashboard",
			Status:              "🚧 under development",
			Progress:            0.4,
			Tag:                 "frontend",
			LatestActivityLabel: "3d ago",
		},
		MilestoneListItem{
			ID:                  "4",
			Name:                "Authentication System",
			Status:              "😴 idle",
			Progress:            0.0,
			Tag:                 "backend",
			LatestActivityLabel: "no activity",
		},
		MilestoneListItem{
			ID:                  "5",
			Name:                "Data Persistence Layer",
			Status:              "😴 idle",
			Progress:            0.0,
			Tag:                 "backend",
			LatestActivityLabel: "no activity",
		},
		MilestoneListItem{
			ID:                  "6",
			Name:                "Testing & QA",
			Status:              "😴 idle",
			Progress:            0.0,
			Tag:                 "testing",
			LatestActivityLabel: "no activity",
		},
		MilestoneListItem{
			ID:                  "7",
			Name:                "Documentation",
			Status:              "🚧 under development",
			Progress:            0.2,
			Tag:                 "docs",
			LatestActivityLabel: "3d ago",
		},
	}
}
