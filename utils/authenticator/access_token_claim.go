package authenticator

import "github.com/golang-jwt/jwt"

type MyClaims struct {
	jwt.StandardClaims
	UserId     string `json:"UserId"`
	Username   string `json:"Username"`
	Name       string `json:"Name"`
	Email      string `json:"Email"`
	Role       string `json:"Role"`
	Status     string `json:"Status"`
	AccessUUID string
}
