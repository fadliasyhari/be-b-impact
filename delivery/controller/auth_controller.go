package controller

import (
	"net/http"

	"be-b-impact.com/csr/delivery/api"
	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/usecase"
	"be-b-impact.com/csr/utils/authenticator"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	router *gin.Engine
	userUC usecase.UsersUseCase
	authUC usecase.AuthUseCase
	api.BaseApi
}

func (au *AuthController) registerHandler(c *gin.Context) {
	var payload model.User

	if err := au.ParseRequestBody(c, &payload); err != nil {
		au.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := au.userUC.SaveData(&payload); err != nil {
		au.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := au.authUC.TokenRegister(&payload)
	if err != nil {
		au.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	data := map[string]interface{}{
		"user":  payload,
		"token": token,
	}

	au.NewSuccessMultiResponse(c, "OK", data)
}

func (au *AuthController) login(c *gin.Context) {
	var payload model.User

	if err := au.ParseRequestBody(c, &payload); err != nil {
		au.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	token, err := au.authUC.Login(payload.Username, payload.Password)
	if err != nil {
		au.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	au.NewSuccessSingleResponse(c, "OK", token)
}

func (au *AuthController) logout(c *gin.Context) {
	token, err := authenticator.BindAuthHeader(c)
	if err != nil {
		c.AbortWithStatus(401)
	}
	err = au.authUC.Logout(token)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Success Logout",
	})
}

func NewAuthController(r *gin.Engine, userUC usecase.UsersUseCase, authUC usecase.AuthUseCase) *AuthController {
	controller := &AuthController{
		router: r,
		userUC: userUC,
		authUC: authUC,
	}
	r.POST("/auth/register", controller.registerHandler)
	r.POST("/auth/login", controller.login)
	r.GET("/auth/logout", controller.logout)
	return controller
}
