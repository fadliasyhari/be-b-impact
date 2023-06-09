package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type User struct {
	BaseModel
	Email    string `gorm:"varchar" json:"email"`
	Username string `gorm:"varchar" json:"username"`
	Password string `gorm:"varchar" json:"password"`
	Role     string `gorm:"varchar" json:"role"`
	Status   string `json:"status"`
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Email, validation.Required),
		validation.Field(&u.Username, validation.Required),
		validation.Field(&u.Password, validation.Required, validation.Length(6, 0)),
		validation.Field(&u.Role, validation.Required, validation.In("admin", "member")),
	)
}
