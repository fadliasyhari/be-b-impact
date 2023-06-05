package dto

import "time"

type ProposalDTO struct {
	ID              string        `json:"id"`
	OrgName         string        `json:"org_name"`
	OrganizatonType string        `json:"organization_type,omitempty"`
	Email           string        `json:"email"`
	Phone           string        `json:"phone"`
	PICName         string        `json:"pic_name"`
	City            string        `json:"city"`
	PostalCode      string        `json:"postal_code"`
	Address         string        `json:"address"`
	Description     string        `json:"description"`
	Status          string        `json:"status"`
	CurrentProgress ProgressDTO   `json:"current_progress"`
	ProjectName     string        `json:"project_name"`
	PartnershipType string        `json:"partnership_type,omitempty"`
	StartDate       time.Time     `json:"start_date"`
	EndDate         time.Time     `json:"end_date"`
	Objective       string        `json:"objective"`
	Alignment       string        `json:"alignment"`
	CreatedBy       string        `json:"created_by"`
	Reviewer        string        `json:"reviewer,omitempty"`
	Files           []FileDTO     `json:"files,omitempty"`
	Progresses      []ProgressDTO `json:"progresses,omitempty"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

type ProgressDTO struct {
	Name      string    `json:"name"`
	Label     string    `json:"label"`
	Status    string    `json:"status"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FileDTO struct {
	Label     string    `json:"label"`
	FileURL   string    `json:"file_url"`
	CreatedAt time.Time `json:"created_at"`
}
