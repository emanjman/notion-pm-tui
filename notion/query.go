package notion

func taskQueryBody(milestoneID, status string, size int) map[string]any {
	return map[string]any{
		"filter": map[string]any{
			"and": []map[string]any{
				{
					"property": "@milestone",
					"relation": map[string]any{"contains": milestoneID},
				},
				{
					"property": "status",
					"status":   map[string]any{"equals": status},
				},
			},
		},
		"sorts": []map[string]any{
			{
				"property":  "priority",
				"direction": "descending",
			},
			{
				"property":  "created-at",
				"direction": "ascending",
			},
		},
		"page_size": size,
	}
}

func milestoneQueryBody(projID string, status MilestoneStatus, size int) map[string]any {
	return map[string]any{
		"filter": map[string]any{
			"and": []map[string]any{
				{
					"property": "@project",
					"relation": map[string]any{"contains": projID},
				},
				{
					"property": "$status",
					"formula": map[string]any{
						"string": map[string]any{
							"equals": status.String(),
						},
					},
				},
			},
		},
		"sorts": []map[string]any{
			{
				"property":  "$latest-activity-at",
				"direction": "descending",
			},
		},
		"page_size": size,
	}
}
