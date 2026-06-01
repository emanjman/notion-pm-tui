package notion

func taskQueryBody(milestoneID, status string, size int) map[string]any {
	var (
		milestonePropName = "@milestone"
		statusPropName    = "status"
		priorityPropName  = "priority"
		createdPropName   = "created-at"
	)

	return map[string]any{
		"filter": map[string]any{
			"and": []map[string]any{
				{
					"property": milestonePropName,
					"relation": map[string]any{"contains": milestoneID},
				},
				{
					"property": statusPropName,
					"status":   map[string]any{"equals": status},
				},
			},
		},
		"sorts": []map[string]any{
			{
				"property":  priorityPropName,
				"direction": "descending",
			},
			{
				"property":  createdPropName,
				"direction": "ascending",
			},
		},
		"page_size": size,
	}
}

func milestoneQueryBody(projID string, status MilestoneStatus, size int) map[string]any {
	var (
		projectPropName  = "r/@project"
		statusPropName   = "$status"
		activityPropName = "$latest-activity-at"
	)

	return map[string]any{
		"filter": map[string]any{
			"and": []map[string]any{
				{
					"property": projectPropName,
					"rollup": map[string]any{
						"any": map[string]any{
							"relation": map[string]any{"contains": projID},
						},
					},
				},
				{
					"property": statusPropName,
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
				"property":  activityPropName,
				"direction": "descending",
			},
		},
		"page_size": size,
	}
}
