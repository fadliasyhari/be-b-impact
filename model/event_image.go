package model

import "gorm.io/gorm"

type EventImage struct {
	BaseModel
	EventID   string         `json:"event_id"`
	Event     Event          `json:"event" gorm:"foreignKey:EventID"`
	ImageURL  string         `gorm:"varchar" json:"image_url"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
