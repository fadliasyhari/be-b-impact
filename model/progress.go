package model

type Progress struct {
	BaseModel
	Name     string     `gorm:"varchar" json:"name"`
	Label    string     `gorm:"varchar" json:"label"`
	Proposal []Proposal `json:"proposal" gorm:"many2many:proposal_progresses;"`
}
