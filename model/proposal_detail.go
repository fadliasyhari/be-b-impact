package model

import "time"

type ProposalDetail struct {
	BaseModel
	ProposalID        string    `json:"proposal_id"`
	Proposal          string    `json:"proposal,omitempty" gorm:"foreignKey:ProposalID"`
	ProjectName       string    `json:"project_name" gorm:"varchar"`
	PartnershipTypeID *string   `json:"partnership_type_id"` // part type refer to category which useFor partnership
	PartnershipType   Category  `json:"partnership_type,omitempty" gorm:"foreignKey:PartnershipTypeID"`
	StartDate         time.Time `json:"start_date"`
	EndDate           time.Time `json:"end_date"`
	Objective         string    `gorm:"text" json:"objective"`
	Alignment         string    `gorm:"text" json:"alignment"`
	AccountableReport string    `gorm:"text" json:"accountable_report"`
}
