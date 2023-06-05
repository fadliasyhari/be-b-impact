package model

import "gorm.io/gorm"

type File struct {
	BaseModel
	Label      string         `gorm:"varchar" json:"label"` // propo_doc || org_profile
	FileURL    string         `gorm:"varchar" json:"file_url"`
	ProposalID string         `json:"proposal_id"`
	Proposal   Proposal       `json:"proposal,omitempty" gorm:"foreignKey:ProposalID"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
