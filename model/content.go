package model

import "gorm.io/gorm"

type Content struct {
	BaseModel
	Title       string         `gorm:"varchar" json:"title"`
	Body        string         `gorm:"text" json:"body"`
	Author      string         `gorm:"varchar" json:"author"`
	Excerpt     string         `gorm:"text" json:"excerpt"`
	Status      string         `gorm:"varchar" json:"status"`
	Category    Category       `json:"category" gorm:"foreignKey:CategoryID"`
	CategoryID  *string        `json:"category_id"`
	CreatedBy   string         `json:"created_by"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	DeletedBy   string         `json:"deleted_by"`
	ProposalID  string         `json:"proposal_id"`
	Image       []Image        `json:"image,omitempty"`
	Tag         []Tag          `json:"tag" gorm:"many2many:tags_contents;"`
	TagsContent []TagsContent  `json:"tags_content"`
}
