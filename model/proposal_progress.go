package model

type ProposalProgress struct {
	BaseModel
	ProposalID string   `json:"proposal_id"`
	ProgressID string   `json:"progress_id"`
	Proposal   Proposal `json:"proposal" gorm:"foreignKey:ProposalID"`
	Progress   Progress `json:"progress" gorm:"foreignKey:ProgressID"`
	Note       string   `gorm:"text" json:"note"`
	Status     string   `json:"status"`
}
