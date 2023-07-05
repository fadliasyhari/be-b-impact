package middleware

import (
	"net/http"

	"be-b-impact.com/csr/utils/authenticator"
	"github.com/gin-gonic/gin"
)

type AuthTokenMiddlerware interface {
	RequireToken() gin.HandlerFunc
}

type authTokenMiddlerware struct {
	acctToken authenticator.AccessToken
}

func (a *authTokenMiddlerware) RequireToken() gin.HandlerFunc {

	return func(c *gin.Context) {
		if c.Request.URL.Path == "/auth" {
			c.Next()
		} else {
			tokenString, err := authenticator.BindAuthHeader(c)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": err.Error(),
				})
				c.Abort()
				return
			}
			accessDetail, err := a.acctToken.VerifyAccessToken(tokenString)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": err.Error(),
				})
				c.Abort()
				return
			}

			err = a.acctToken.FetchAccessToken(accessDetail)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "Error Fetch to Redis :" + err.Error(),
				})
				c.Abort()
				return
			}
			c.Set("user", accessDetail)
			c.Next()
		}
	}
}

func NewTokenValidator(acctToken authenticator.AccessToken) AuthTokenMiddlerware {
	return &authTokenMiddlerware{
		acctToken: acctToken,
	}
}
