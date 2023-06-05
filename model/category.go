package model

type Category struct {
	BaseModel
	Parent    string `gorm:"varchar" json:"parent"`
	Name      string `gorm:"varchar" json:"name"`
	UseFor    string `gorm:"varchar" json:"use_for"`
	Status    string `json:"status"`
	CreatedBy string `json:"created_by"`
}
