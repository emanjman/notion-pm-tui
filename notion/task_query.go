package notion

import "strconv"

func addTaskBody(tasksDatasourceID, title, milestoneID, status, taskType string, priority int) map[string]any {
	return map[string]any{
		"parent": map[string]any{
			"data_source_id": tasksDatasourceID,
		},
		"properties": map[string]any{
			taskPropTitle: TitleProperty{
				Title: []RichText{{Text: TextContent{Content: title}}},
			},
			taskPropMilestoneRelation: RelationProperty{
				Relation: []RelationItem{{ID: milestoneID}},
			},
			taskPropStatus: map[string]any{
				"status": map[string]any{"name": status},
			},
			taskPropTypeSelect: map[string]any{
				"select": map[string]any{"name": taskType},
			},
			taskPropPriority: map[string]any{
				"select": map[string]any{"name": strconv.Itoa(priority)},
			},
		},
	}
}

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
