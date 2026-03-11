package notion

type BlockType string

const (
	// must handle
	BulletedListItem BlockType = "bulleted_list_item"
	NumberedListItem BlockType = "numbered_list_item"
	Callout          BlockType = "callout"
	Divider          BlockType = "divider" // block unnecessary
	Heading2         BlockType = "heading_2"
	Heading3         BlockType = "heading_3"
	Paragraph        BlockType = "paragraph"
	Toggle           BlockType = "toggle"
	ChildPage        BlockType = "child_page"
	Code             BlockType = "code"

	// could appear (handle)
	Equation   BlockType = "equation"
	ColumnList BlockType = "column_list"
	Column     BlockType = "column"
	Table      BlockType = "table"
	TableRow   BlockType = "table_row"

	// could appear (hide)
	Breadcrumb    BlockType = "breadcrumb"
	ChildDatabase BlockType = "child_database"
	Image         BlockType = "image"

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
	LinkPreview     BlockType = "link_preview"
	Unsupported     BlockType = "unsupported"
)

// ---------------------------------

type Block struct {
	ID          string      `json:"id"`
	Parent      ParentBlock `json:"parent"`
	InTrash     bool        `json:"in_trash"`
	HasChildren bool        `json:"has_children"`
	Type        BlockType   `json:"type"`

	// manually populated, not from json
	Children []Block `json:"-"`

	// dynamic type object
	BulletedListItem *BulletedListItemBlock `json:"bulleted_list_item,omitempty"`
	NumberedListItem *NumberedListItemBlock `json:"numbered_list_item,omitempty"`
	Callout          *CalloutBlock          `json:"callout,omitempty"`
	Divider          *struct{}              `json:"divider,omitempty"`
	Heading2         *HeadingBlock          `json:"heading_2,omitempty"`
	Heading3         *HeadingBlock          `json:"heading_3,omitempty"`
	Paragraph        *ParagraphBlock        `json:"paragraph,omitempty"`
	Toggle           *ToggleBlock           `json:"toggle,omitempty"`
	ChildPage        *ChildPageBlock        `json:"child_page,omitempty"`
	Code             *CodeBlock             `json:"code,omitempty"`

	// uncommon, handle anyways
	Equation   *EquationBlock `json:"equation,omitempty"`
	ColumnList *struct{}      `json:"column_list,omitempty"`
	Column     *ColumnBlock   `json:"column,omitempty"`
	Table      *TableBlock    `json:"table,omitempty"`
	TableRow   *TableRowBlock `json:"table_row,omitempty"`

	// expect to exist, should ignore
	Breadcrumb    *struct{}           `json:"breadcrumb,omitempty"`
	ChildDatabase *ChildDatabaseBlock `json:"child_database,omitempty"`
	Image         *ImageBlock         `json:"image,omitempty"`
}

type ParentBlock struct {
	BlockID string `json:"block_id"`
}

// ---------------------------------

type BulletedListItemBlock struct {
	RichText []RichText `json:"rich_text"`
	Color    Color      `json:"color"`
	Children []Block    `json:"children"`
}

type NumberedListItemBlock struct {
	RichText       []RichText      `json:"rich_text"`
	Color          Color           `json:"color"`
	ListStartIndex *int            `json:"list_start_index,omitempty"`
	ListFormat     *ListFormatType `json:"list_format"`
	Children       []Block         `json:"children"`
}

type CalloutBlock struct {
	RichText []RichText `json:"rich_text"`
	Color    Color      `json:"color"`
}

type HeadingBlock struct {
	RichText []RichText `json:"rich_text"`
	Color    Color      `json:"color"`
}

type ListFormatType string

const (
	Numbers ListFormatType = "numbers"
	Letters ListFormatType = "letters"
	Roman   ListFormatType = "roman"
)

type ParagraphBlock struct {
	RichText []RichText `json:"rich_text"`
	Color    Color      `json:"color"`
	Children []Block    `json:"children"`
}

type ToggleBlock struct {
	RichText []RichText `json:"rich_text"`
	Color    Color      `json:"color"`
	Children []Block    `json:"children"`
}

type ChildPageBlock struct {
	Title string `json:"title"`
}

type CodeBlock struct {
	Caption  []Block    `json:"caption"`
	RichText []RichText `json:"rich_text"`
	Language string     `json:"language"`
}

type EquationBlock struct {
	Expression string `json:"expression"`
}

type ColumnBlock struct {
	WidthRatio int `json:"width_ratio"` // 0-1
}

type TableBlock struct {
	ColumnCount     int  `json:"table_width"`
	HasColumnHeader bool `json:"has_column_header"`
	HasRowHeader    bool `json:"has_row_header"`
}

type TableRowBlock struct {
	Cells []RichText `json:"cells"`
}

type ChildDatabaseBlock struct {
	Title string `json:"title"`
}

type FileType string

const (
	External   FileType = "external"
	HostedFile FileType = "file"
	FileUpload FileType = "file_upload"
)

type ImageBlock struct {
	Type FileType `json:"type"`

	FileUpload *struct {
		ID string `json:"id"`
	} `json:"file_upload"`

	File *struct {
		URL        string `json:"url"`
		ExpiryTime string `json:"expiry_time"`
	} `json:"file"`

	External *struct {
		URL string `json:"url"`
	} `json:"external"`
}
