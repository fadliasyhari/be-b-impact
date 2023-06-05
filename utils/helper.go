package utils

import (
	"net/http"

	"be-b-impact.com/csr/delivery/api"
	"be-b-impact.com/csr/utils/authenticator"
	"github.com/gin-gonic/gin"
)

func AccessInsideToken(b api.BaseApi, c *gin.Context) authenticator.AccessDetail {
	user, exists := c.Get("user")
	if !exists {
		b.NewFailedResponse(c, http.StatusUnauthorized, "token invalid")
		return authenticator.AccessDetail{}
	}
	userTyped, ok := user.(authenticator.AccessDetail)
	if !ok {
		b.NewFailedResponse(c, http.StatusUnauthorized, "user invalid")
		return authenticator.AccessDetail{}
	}
	return userTyped
}
