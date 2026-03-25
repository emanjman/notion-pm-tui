package task

func mockItems() []Item {
	return []Item{
		{ID: "1", MilestoneID: "1", Task: "Implement task list model", Status: "dev", Priority: 1, Type: "feat"},
		{ID: "2", MilestoneID: "1", Task: "Fix milestone progress formula", Status: "dev", Priority: 2, Type: "fix"},
		{ID: "3", MilestoneID: "1", Task: "Refactor notion client", Status: "dev", Priority: 3, Type: "refactor"},
		{ID: "4", MilestoneID: "1", Task: "Add status bar component", Status: "idle", Priority: 2, Type: "feat"},
		{ID: "5", MilestoneID: "1", Task: "Clean up delegate styles", Status: "idle", Priority: 4, Type: "style"},
		{ID: "6", MilestoneID: "2", Task: "Add CI pipeline", Status: "idle", Priority: 5, Type: "chore"},
		{ID: "7", MilestoneID: "2", Task: "Setup project structure", Status: "done", Priority: 1, Type: "chore"},
		{ID: "8", MilestoneID: "3", Task: "Fix key binding conflicts", Status: "done", Priority: 3, Type: "fix"},
		{ID: "9", MilestoneID: "4", Task: "Remove unused imports", Status: "archive", Priority: 5, Type: "refactor"},
		{ID: "10", MilestoneID: "4", Task: "Pad list items under headers", Status: "archive", Priority: 4, Type: "style"},
	}
}
