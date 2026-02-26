package tasklist

func mockTaskItems() []TaskListItem {
	return []TaskListItem{
		{ID: "1", Title: "Implement task list model", Status: "dev", Priority: 1, Type: "feat"},
		{ID: "2", Title: "Fix milestone progress formula", Status: "dev", Priority: 2, Type: "fix"},
		{ID: "3", Title: "Refactor notion client", Status: "dev", Priority: 3, Type: "refactor"},
		{ID: "4", Title: "Add status bar component", Status: "idle", Priority: 2, Type: "feat"},
		{ID: "5", Title: "Clean up delegate styles", Status: "idle", Priority: 4, Type: "style"},
		{ID: "6", Title: "Add CI pipeline", Status: "idle", Priority: 5, Type: "chore"},
		{ID: "7", Title: "Setup project structure", Status: "done", Priority: 1, Type: "chore"},
		{ID: "8", Title: "Fix key binding conflicts", Status: "done", Priority: 3, Type: "fix"},
		{ID: "9", Title: "Remove unused imports", Status: "archive", Priority: 5, Type: "refactor"},
		{ID: "10", Title: "Pad list items under headers", Status: "archive", Priority: 4, Type: "style"},
	}
}
