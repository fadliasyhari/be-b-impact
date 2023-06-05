package model

import (
	"gorm.io/gorm"
)

type Proposal struct {
	BaseModel
	OrgName            string             `gorm:"varchar" json:"org_name"`
	OrganizationTypeID string             `json:"organization_type_id"` // org type refer to category which useFor organization
	OrganizatonType    Category           `json:"organization_type,omitempty" gorm:"foreignKey:OrganizationTypeID"`
	Email              string             `gorm:"varchar" json:"email"`
	Phone              string             `gorm:"varchar" json:"phone"`
	PICName            string             `gorm:"varchar" json:"pic_name"`
	City               string             `gorm:"varchar" json:"city"`
	PostalCode         string             `gorm:"varchar" json:"postal_code"`
	Address            string             `gorm:"text" json:"address"`
	Description        string             `gorm:"text" json:"description"`
	Status             string             `gorm:"varchar" json:"status"`
	CreatedBy          string             `json:"created_by"`
	DeletedAt          gorm.DeletedAt     `gorm:"index" json:"-"`
	DeletedBy          string             `json:"deleted_by"`
	ReviewerID         string             `json:"reviewer_id,omitempty"`
	ProposalDetail     ProposalDetail     `json:"proposal_detail,omitempty"`
	File               []File             `json:"file,omitempty"`
	Progress           []Progress         `json:"progress,omitempty" gorm:"many2many:proposal_progresses;"`
	ProposalProgress   []ProposalProgress `json:"proposal_progress"`
}
