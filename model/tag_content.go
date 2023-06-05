package model

import "gorm.io/gorm"

type TagsContent struct {
	BaseModel
	TagID     string         `json:"tag_id"`
	ContentID string         `json:"content_id"`
	Tag       Tag            `json:"tag" gorm:"foreignKey:TagID"`
	Content   Content        `json:"content" gorm:"foreignKey:ContentID"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
