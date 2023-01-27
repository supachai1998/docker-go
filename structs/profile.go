package structs

type Profile struct {
	UserID        string `json:"user_id"`
	DisplayName   string `json:"display_name"`
	PictureURL    string `json:"picture_url"`
	StatusMessage string `json:"status_message"`
}
