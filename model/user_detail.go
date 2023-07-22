package model

import "gorm.io/gorm"

type UserDetail struct {
	BaseModel
	UserID    string         `json:"user_id"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	ImageURL  string         `gorm:"varchar" json:"image_url"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
