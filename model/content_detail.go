package model

type ContentDetail struct {
	BaseModel
	ContentID string `json:"content_id"`
	UserID    string `json:"user_id"`
}
