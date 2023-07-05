package model

import "gorm.io/gorm"

type EventParticipant struct {
	BaseModel
	UserID    string         `json:"user_id"`
	EventID   string         `json:"event_id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Phone     string         `json:"phone"`
	Status    string         `json:"status"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	Event     Event          `json:"event" gorm:"foreignKey:EventID"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
