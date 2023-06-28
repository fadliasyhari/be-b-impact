package model

import (
	"errors"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type User struct {
	BaseModel
	Name     string `gorm:"varchar" json:"name"`
	Phone    string `gorm:"varchar" json:"phone"`
	Email    string `gorm:"varchar" json:"email"`
	Username string `gorm:"varchar" json:"username"`
	Password string `gorm:"varchar" json:"password"`
	Role     string `gorm:"varchar" json:"role"`
	Status   string `json:"status"`
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Username, validation.By(checkNoSpaces)),
		validation.Field(&u.Password, validation.Required, validation.Length(6, 0)),
		validation.Field(&u.Role, validation.Required, validation.In("admin", "member")),
	)
}

func checkNoSpaces(value interface{}) error {
	username, ok := value.(string)
	if !ok {
		return errors.New("username should be string")
	}

	if strings.Contains(username, " ") {
		return errors.New("username cannot contain spaces")
	}

	return nil
}
