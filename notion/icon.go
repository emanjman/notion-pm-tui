package notion

type IconType string

const (
	IconEmoji      IconType = "emoji"
	IconFile       IconType = "file"
	IconFileUpload IconType = "file_upload"
)

type Icon struct {
	Type IconType `json:"type"`

	Emoji *string `json:"emoji"`
	File  *struct {
		URL        string `json:"url"`
		ExpiryTime string `json:"expiry_time"`
	} `json:"file"`
	FileUpload *struct {
		ID string `json:"id"`
	} `json:"file_upload"`
}
