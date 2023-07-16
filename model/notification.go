package model

import "gorm.io/gorm"

type Notification struct {
	BaseModel
	Body      string         `gorm:"varchar" json:"body"`
	UserID    string         `gorm:"varchar" json:"user_id"`
	Type      string         `gorm:"varchar" json:"type"`
	TypeID    string         `gorm:"varchat" json:"type_id"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
