package model

type Tag struct {
	BaseModel
	Name        string        `gorm:"varchar" json:"name"`
	Content     []Content     `json:"content" gorm:"many2many:tags_contents;"`
	TagsContent []TagsContent `json:"tags_content"`
}
