package model

import "time"

type ProposalProgress struct {
	BaseModel
	ProposalID     string    `json:"proposal_id"`
	ProgressID     string    `json:"progress_id"`
	Proposal       Proposal  `json:"proposal" gorm:"foreignKey:ProposalID"`
	Progress       Progress  `json:"progress" gorm:"foreignKey:ProgressID"`
	Note           string    `gorm:"text" json:"note"`
	ReviewLocation string    `json:"review_location"`
	ReviewDate     time.Time `json:"review_date"`
	ReviewCP       string    `json:"review_cp"`
	ReviewFeedback string    `json:"review_feedback"`
	Status         string    `json:"status"`
}
