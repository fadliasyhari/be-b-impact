package authenticator

import (
	"errors"
	"net/http"
	"strings"

	"be-b-impact.com/csr/delivery/api/request"
	"github.com/gin-gonic/gin"
)

func BindAuthHeader(c *gin.Context) (string, error) {
	var h request.AuthHeader
	if err := c.ShouldBindHeader(&h); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		c.Abort()
		return "", err
	}
	tokenString := strings.Replace(h.AuthorizationHeader, "Bearer ", "", -1)
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		c.Abort()
		return "", errors.New("token is empty")
	}
	return tokenString, nil
}
