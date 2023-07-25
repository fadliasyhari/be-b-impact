package dto

import (
	"time"

	"be-b-impact.com/csr/model"
)

type EventDTO struct {
	model.Event
	TotalParticipant int `json:"total_participant"`
}

type EventDTOResponse struct {
	ID               string      `json:"id"`
	Title            string      `json:"title"`
	Description      string      `json:"description"`
	Location         string      `json:"location"`
	StartDate        string      `json:"start_date"`
	EndDate          string      `json:"end_date"`
	Status           string      `json:"status"`
	IsJoined         bool        `json:"is_joined"`
	ParticipantID    string      `json:"participant_id"`
	Category         string      `json:"category"`
	CategoryDetail   CategoryDTO `json:"category_detail"`
	ImageURLs        []ImageDTO  `json:"image_urls"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
	TotalParticipant int         `json:"total_participant"`
}

type EventParticipantDto struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}
