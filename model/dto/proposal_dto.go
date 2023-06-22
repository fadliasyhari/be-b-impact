package dto

import (
	"time"

	"be-b-impact.com/csr/model"
	"gorm.io/gorm"
)

type ProposalDTO struct {
	ID                string        `json:"id"`
	OrgName           string        `json:"org_name"`
	OrganizatonType   CategoryDTO   `json:"organization_type,omitempty"`
	Email             string        `json:"email"`
	Phone             string        `json:"phone"`
	PICName           string        `json:"pic_name"`
	City              string        `json:"city"`
	PostalCode        string        `json:"postal_code"`
	Address           string        `json:"address"`
	Description       string        `json:"description"`
	Status            string        `json:"status"`
	CurrentProgress   string        `json:"current_progress"`
	ProposalDetailID  string        `json:"proposal_detail_id"`
	ProjectName       string        `json:"project_name"`
	PartnershipType   CategoryDTO   `json:"partnership_type,omitempty"`
	StartDate         time.Time     `json:"start_date"`
	EndDate           time.Time     `json:"end_date"`
	Objective         string        `json:"objective"`
	Alignment         string        `json:"alignment"`
	AccountableReport string        `json:"accountable_report"`
	CreatedBy         string        `json:"created_by"`
	Reviewer          string        `json:"reviewer,omitempty"`
	Files             []FileDTO     `json:"files,omitempty"`
	Progresses        []ProgressDTO `json:"progresses,omitempty"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
}

type ProgressDTO struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Label          string    `json:"label"`
	Status         string    `json:"status"`
	Note           string    `json:"note"`
	ReviewLocation string    `json:"review_location,omitempty"`
	ReviewDate     time.Time `json:"review_date,omitempty"`
	ReviewCP       string    `json:"review_cp,omitempty"`
	ReviewFeedback string    `json:"review_feedback,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type FileDTO struct {
	ID        string    `json:"id"`
	Label     string    `json:"label"`
	FileURL   string    `json:"file_url"`
	CreatedAt time.Time `json:"created_at"`
}

type Proposal struct {
	ID                 string                   `json:"id"`
	CreatedAt          time.Time                `json:"created_at"`
	UpdatedAt          time.Time                `json:"updated_at"`
	OrgName            string                   `gorm:"varchar" json:"org_name"`
	OrganizationTypeID string                   `json:"organization_type_id"` // org type refer to category which useFor organization
	OrganizatonType    model.Category           `json:"organization_type,omitempty" gorm:"foreignKey:OrganizationTypeID"`
	Email              string                   `gorm:"varchar" json:"email"`
	Phone              string                   `gorm:"varchar" json:"phone"`
	PICName            string                   `gorm:"varchar" json:"pic_name"`
	City               string                   `gorm:"varchar" json:"city"`
	CurrentProgress    string                   `json:"current_progress"`
	PostalCode         string                   `gorm:"varchar" json:"postal_code"`
	Address            string                   `gorm:"text" json:"address"`
	Description        string                   `gorm:"text" json:"description"`
	Status             string                   `gorm:"varchar" json:"status"`
	CreatedBy          string                   `json:"created_by"`
	DeletedAt          gorm.DeletedAt           `gorm:"index" json:"-"`
	DeletedBy          string                   `json:"deleted_by"`
	ReviewerID         string                   `json:"reviewer_id,omitempty"`
	ProposalDetail     model.ProposalDetail     `json:"proposal_detail,omitempty"`
	File               []model.File             `json:"file,omitempty"`
	Progress           []model.Progress         `json:"progress,omitempty" gorm:"many2many:proposal_progresses;"`
	ProposalProgress   []model.ProposalProgress `json:"proposal_progress"`
}
