package model

type User struct {
	BaseModel
	Email    string `gorm:"varchar" json:"email"`
	Username string `gorm:"varchar" json:"username"`
	Password string `gorm:"varchar" json:"password"`
	Role     string `gorm:"varchar" json:"role"`
	Status   string `json:"status"`
}
