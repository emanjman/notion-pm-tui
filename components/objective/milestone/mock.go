package milestone

func mockItems() []Item {
	return []Item{
		Item{
			ID:                  "1",
			Name:                "Setup Project Structure",
			Status:              "🎉 complete",
			Progress:            1.0,
			Tag:                 "backend",
			LatestActivityLabel: "12/8/25",
		},
		Item{
			ID:                  "2",
			Name:                "Implement Notion API",
			Status:              "🚧 under development",
			Progress:            0.75,
			Tag:                 "backend",
			LatestActivityLabel: "today",
		},
		Item{
			ID:                  "3",
			Name:                "Build TUI Dashboard",
			Status:              "🚧 under development",
			Progress:            0.4,
			Tag:                 "frontend",
			LatestActivityLabel: "3d ago",
		},
		Item{
			ID:                  "4",
			Name:                "Authentication System",
			Status:              "😴 idle",
			Progress:            0.0,
			Tag:                 "backend",
			LatestActivityLabel: "no activity",
		},
		Item{
			ID:                  "5",
			Name:                "Data Persistence Layer",
			Status:              "😴 idle",
			Progress:            0.0,
			Tag:                 "backend",
			LatestActivityLabel: "no activity",
		},
		Item{
			ID:                  "6",
			Name:                "Testing & QA",
			Status:              "😴 idle",
			Progress:            0.0,
			Tag:                 "testing",
			LatestActivityLabel: "no activity",
		},
		Item{
			ID:                  "7",
			Name:                "Documentation",
			Status:              "🚧 under development",
			Progress:            0.2,
			Tag:                 "docs",
			LatestActivityLabel: "3d ago",
		},
	}
}
