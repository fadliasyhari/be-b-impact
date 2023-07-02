package controller

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"be-b-impact.com/csr/delivery/api"
	"be-b-impact.com/csr/delivery/api/middleware"
	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/usecase"
	"be-b-impact.com/csr/utils"
	"be-b-impact.com/csr/utils/authenticator"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	router *gin.Engine
	userUC usecase.UsersUseCase
	authUC usecase.AuthUseCase
	api.BaseApi
}

func generateOTP() string {
	min := 100000
	max := 999999
	src := rand.NewSource(time.Now().UnixNano()) // Create a new random source
	r := rand.New(src)                           // Create a new random generator
	otp := r.Intn(max-min+1) + min
	return fmt.Sprintf("%06d", otp)
}

func (au *AuthController) registerHandler(c *gin.Context) {
	var payload model.User

	if err := au.ParseRequestBody(c, &payload); err != nil {
		au.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if payload.Role == "" {
		payload.Role = "member"
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

func (au *AuthController) registerAdminHandler(c *gin.Context) {
	userTyped := utils.AccessInsideToken(au.BaseApi, c)
	if userTyped.Role != "super" {
		au.NewFailedResponse(c, http.StatusForbidden, "access denied")
		return
	}

	var payload model.User

	if err := au.ParseRequestBody(c, &payload); err != nil {
		au.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if payload.Role == "" || payload.Role == "admin" {
		payload.Role = "admin"
	} else {
		au.NewFailedResponse(c, http.StatusBadRequest, "invalid role")
		return
	}

	if err := au.userUC.SaveData(&payload); err != nil {
		au.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	au.NewSuccessSingleResponse(c, "OK", payload)
}

func (au *AuthController) sendOTP(c *gin.Context) {
	var payload struct {
		Email string `json:"email"`
	}

	if err := au.ParseRequestBody(c, &payload); err != nil {
		au.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	ctx := context.Background()

	_, err := au.authUC.GetOTP(ctx, payload.Email)
	if err != nil {
		// Delete the previous OTP from Redis
		err := au.authUC.DeleteOTP(ctx, payload.Email)
		if err != nil {
			au.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// Generate a new OTP
	newOTP := generateOTP()

	// Store the new OTP in Redis
	err = au.authUC.StoreOTP(ctx, payload.Email, newOTP)
	if err != nil {
		au.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Send the new OTP email
	err = utils.SendOTP(payload.Email, newOTP)
	if err != nil {
		au.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	au.NewSuccessSingleResponse(c, "OTP sent successfully", http.StatusOK)
}

func (au *AuthController) verifyOTP(c *gin.Context) {
	var payload struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := au.ParseRequestBody(c, &payload); err != nil {
		au.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	ctx := context.Background()

	// Retrieve the OTP from Redis
	savedOTP, err := au.authUC.GetOTP(ctx, payload.Email)
	if err != nil {
		au.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Compare the OTP
	if savedOTP != payload.OTP {
		au.NewFailedResponse(c, http.StatusBadRequest, "Invalid OTP")
		return
	}

	// OTP verification successful
	au.NewSuccessSingleResponse(c, "OTP verification successful", http.StatusOK)
}

func (au *AuthController) login(c *gin.Context) {
	var payload model.User

	if err := au.ParseRequestBody(c, &payload); err != nil {
		au.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	currentUser, token, err := au.authUC.Login(payload.Email, payload.Password)
	if err != nil {
		au.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	data := map[string]interface{}{
		"user_id":  currentUser.ID,
		"name":     currentUser.Name,
		"email":    currentUser.Email,
		"username": currentUser.Username,
		"role":     currentUser.Role,
		"token":    token,
	}

	au.NewSuccessMultiResponse(c, "OK", data)
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

func (au *AuthController) forgetPassword(c *gin.Context) {
	var payload model.User

	if err := au.ParseRequestBody(c, &payload); err != nil {
		au.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Perform necessary logic for forget password, e.g., generate reset token, store in database, etc.

	// Generate a password reset link
	resetLink := "https://example.com/reset-password?token=YOUR_RESET_TOKEN"

	// Send password reset email
	err := utils.SendResetEmail(payload.Email, resetLink)
	if err != nil {
		au.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	au.NewSuccessSingleResponse(c, "Password reset email has been sent successfully", http.StatusOK)
}

func NewAuthController(r *gin.Engine, userUC usecase.UsersUseCase, authUC usecase.AuthUseCase, tokenMdw middleware.AuthTokenMiddlerware) *AuthController {
	controller := &AuthController{
		router: r,
		userUC: userUC,
		authUC: authUC,
	}
	r.POST("/auth/register-admin", tokenMdw.RequireToken(), controller.registerAdminHandler)
	r.POST("/auth/register", controller.registerHandler)
	r.POST("/auth/verify-otp", controller.verifyOTP)
	r.POST("/auth/send-otp", controller.sendOTP)
	r.POST("/auth/login", controller.login)
	r.POST("/auth/forget-password", controller.forgetPassword)
	r.GET("/auth/logout", controller.logout)
	return controller
}
