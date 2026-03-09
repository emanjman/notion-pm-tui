package notion

type BlockType string

const (
	// must handle
	BulletedListItem BlockType = "bulleted_list_item"
	Callout          BlockType = "callout"
	ChildPage        BlockType = "child_page"
	Divider          BlockType = "divider"
	Heading2         BlockType = "heading_2"
	Heading3         BlockType = "heading_3"
	NumberedListItem BlockType = "numbered_list_item"
	Paragraph        BlockType = "paragraph"
	Toggle           BlockType = "toggle"

	// could appear (handle)
	Column      BlockType = "column"
	ColumnList  BlockType = "column_list"
	Equation    BlockType = "equation"
	LinkPreview BlockType = "link_preview"
	Table       BlockType = "table"
	TableRow    BlockType = "table_row"

	// could appear (hide)
	Breadcrumb    BlockType = "breadcrumb"
	ChildDatabase BlockType = "child_database"
	Image         BlockType = "image"
	Unsupported   BlockType = "unsupported"

	// shouldn't appear
	Bookmark        BlockType = "bookmark"
	Embed           BlockType = "embed"
	File            BlockType = "file"
	Heading1        BlockType = "heading_1"
	PDF             BlockType = "pdf"
	Quote           BlockType = "quote"
	SyncedBlock     BlockType = "synced_block"
	TableOfContents BlockType = "table_of_contents"
	Template        BlockType = "template"
	ToDo            BlockType = "to_do"
	Transcription   BlockType = "transcription"
	Video           BlockType = "video"
)

// ---------------------------------

type Block struct {
	ID     string      `json:"id"`
	Parent ParentBlock `json:"parent"`
	Type   BlockType   `json:"type"`
}

type ParentBlock struct {
	BlockID string `json:"block_id"`
}

// ---------------------------------
