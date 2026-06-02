package notion

func queryMilestoneBody(projID string, status MilestoneStatus, size int) map[string]any {
	return map[string]any{
		"filter": map[string]any{
			"and": []map[string]any{
				{
					"property": milestonePropProjectRollupRelation,
					"rollup": map[string]any{
						"any": map[string]any{
							"relation": map[string]any{"contains": projID},
						},
					},
				},
				{
					"property": milestonePropStatusFormula,
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
				"property":  milestonePropLatestActivityAt,
				"direction": "descending",
			},
		},
		"page_size": size,
	}
}

func addMilestoneBody(milestonesDatasourceID, title, versionID string) map[string]any {
	return map[string]any{
		"parent": map[string]any{
			"data_source_id": milestonesDatasourceID,
		},
		"properties": map[string]any{
			milestonePropTitle: TitleProperty{
				Title: []RichText{{Text: TextContent{Content: title}}},
			},
			milestonePropVersionRelation: RelationProperty{
				Relation: []RelationItem{{ID: versionID}},
			},
		},
	}
}
