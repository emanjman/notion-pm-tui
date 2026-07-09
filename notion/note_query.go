package notion

func queryNoteBody(projID string, size int) map[string]any {
	return map[string]any{
		"filter": map[string]any{
			"property": notePropProjectRelation,
			"relation": map[string]any{
				"contains": projID,
			},
		},
		"sorts": []map[string]any{
			{
				"timestamp": "created_time",
				"direction": "descending",
			},
		},
		"page_size": size,
	}
}
