package notion

func taskQueryBody(milestoneID, status string, size int) map[string]any {
	return map[string]any{
		"filter": map[string]any{
			"and": []map[string]any{
				{
					"property": taskPropMilestoneRelation,
					"relation": map[string]any{"contains": milestoneID},
				},
				{
					"property": taskPropStatus,
					"status":   map[string]any{"equals": status},
				},
			},
		},
		"sorts": []map[string]any{
			{
				"property":  taskPropPriority,
				"direction": "descending",
			},
			{
				"property":  taskPropCreatedAt,
				"direction": "ascending",
			},
		},
		"page_size": size,
	}
}
