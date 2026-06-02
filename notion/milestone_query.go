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

// todo: wire the @version datasource so the version can be selected per-project.
// for now the demo project ("Hoop Archives") has a single version, so we hardcode
// its page id. a milestone must hang off a @version (the @project is a rollup
// through it), otherwise it won't roll up to any project.
const demoVersionPageID = "346b7273-944b-80ee-bc8d-e9ead7e1e623"

func addMilestoneBody(milestonesDatasourceId, title string) map[string]any {
	return map[string]any{
		"parent": map[string]any{
			"data_source_id": milestonesDatasourceId,
		},
		// only writable props — the formula/rollup props (progress, $status,
		// task-ct, r/@project) are read-only and rejected by the create endpoint.
		"properties": map[string]any{
			milestonePropTitle: TitleProperty{
				Title: []RichText{{Text: TextContent{Content: title}}},
			},
			milestonePropVersionRelation: RelationProperty{
				Relation: []RelationItem{{ID: demoVersionPageID}},
			},
		},
	}
}
