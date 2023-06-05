package model

type Profile struct {
	BaseModel
	Name         string `gorm:"varchar" json:"name"`
	Phone        string `gorm:"varchar" json:"phone"`
	Organization string `gorm:"varchar" json:"organization"`
	UserID       int    `json:"user_id"`
}
