package model

import "gorm.io/gorm"

type Event struct {
	BaseModel
	Title            string             `gorm:"varchar" json:"title"`
	Description      string             `gorm:"text" json:"description"`
	Location         string             `gorm:"varchar" json:"location"`
	StartDate        string             `gorm:"varchar" json:"start_date"`
	EndDate          string             `gorm:"text" json:"end_date"`
	Status           string             `gorm:"varchar" json:"status"`
	Category         Category           `json:"category" gorm:"foreignKey:CategoryID"`
	CategoryID       *string            `json:"category_id"`
	CreatedBy        string             `json:"created_by"`
	DeletedAt        gorm.DeletedAt     `gorm:"index" json:"-"`
	DeletedBy        string             `json:"deleted_by"`
	EventImage       []EventImage       `json:"event_image,omitempty"`
	User             []User             `json:"user" gorm:"many2many:event_participants;"`
	EventParticipant []EventParticipant `json:"event_participant"`
}
