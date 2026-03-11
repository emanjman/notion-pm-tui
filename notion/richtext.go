package notion

type RichText struct {
	PlainText   string `json:"plain_text"`
	Annotations struct {
		Bold          bool  `json:"bold"`
		Italic        bool  `json:"italic"`
		Strikethrough bool  `json:"strikethrough"`
		Underline     bool  `json:"underline"`
		InlineCode    bool  `json:"code"`
		Color         Color `json:"color"`
	} `json:"annotations"`
}
