package notion

func queryVersionBody(projID string) map[string]any {
	return map[string]any{
		"filter": map[string]any{
			"property": versionPropProjectRelation,
			"relation": map[string]any{
				"contains": projID,
			},
		},
		"sorts": []map[string]any{
			{
				"property":  versionPropCreatedAt,
				"direction": "descending",
			},
		},
		"page_size": 10, // hardcoded sane size
	}
}
