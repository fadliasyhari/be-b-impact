package model

import "gorm.io/gorm"

type Image struct {
	BaseModel
	ContentID string         `json:"content_id"`
	Content   Content        `json:"content" gorm:"foreignKey:ContentID"`
	ImageURL  string         `gorm:"varchar" json:"image_url"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
