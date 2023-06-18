package controller

import (
	"net/http"
	"strconv"

	"be-b-impact.com/csr/delivery/api"
	"be-b-impact.com/csr/delivery/api/middleware"
	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"
	"be-b-impact.com/csr/usecase"
	"be-b-impact.com/csr/utils"
	"github.com/gin-gonic/gin"
)

type UsersController struct {
	router  *gin.Engine
	useCase usecase.UsersUseCase
	api.BaseApi
}

func (us *UsersController) createHandler(c *gin.Context) {
	var payload model.User
	if err := us.ParseRequestBody(c, &payload); err != nil {
		us.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// only user with role super who can create user admin
	if payload.Role == "admin" {
		userTyped := utils.AccessInsideToken(us.BaseApi, c)
		if userTyped.Role != "super" {
			us.NewFailedResponse(c, http.StatusForbidden, "access denied")
			return
		}
	}

	if err := us.useCase.SaveData(&payload); err != nil {
		us.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	us.NewSuccessSingleResponse(c, "OK", payload)
}

func (us *UsersController) listHandler(c *gin.Context) {
	// only super admin can get list of users
	userTyped := utils.AccessInsideToken(us.BaseApi, c)
	if userTyped.Role != "super" {
		us.NewFailedResponse(c, http.StatusForbidden, "access denied")
		return
	}
	filter := make(map[string]interface{})

	// Iterate over the query parameters
	for key, values := range c.Request.URL.Query() {
		// Skip if the key is empty or has multiple values
		if key == "page" || key == "limit" || key == "order" || key == "sort" || key == "" || len(values) != 1 {
			continue
		}

		// Add key-value pair to the filter
		filter[key] = values[0]
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		us.NewFailedResponse(c, http.StatusBadRequest, "invalid page number")
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		us.NewFailedResponse(c, http.StatusBadRequest, "invalid limit number")
		return
	}
	order := c.DefaultQuery("order", "id")
	sort := c.DefaultQuery("sort", "ASC")
	requestQueryParams := dto.RequestQueryParams{
		QueryParams: dto.QueryParams{
			Sort:  sort,
			Order: order,
		},
		PaginationParam: dto.PaginationParam{
			Page:  page,
			Limit: limit,
		},
		Filter: filter,
	}
	users, paging, err := us.useCase.Pagination(requestQueryParams)
	if err != nil {
		us.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	var usersInterface []interface{}
	for _, us := range users {
		usersInterface = append(usersInterface, us)
	}
	us.NewSuccessPagedResponse(c, "OK", usersInterface, paging)
}

func (us *UsersController) getHandler(c *gin.Context) {
	userTyped := utils.AccessInsideToken(us.BaseApi, c)
	id := c.Param("id")
	users, err := us.useCase.FindById(id)
	if err != nil {
		us.NewFailedResponse(c, http.StatusNotFound, err.Error())
		return
	}
	// if it's not super admin, user only can get their user detail
	if users.ID != userTyped.UserId && userTyped.Role != "super" {
		us.NewFailedResponse(c, http.StatusForbidden, "access denied")
		return
	}
	us.NewSuccessSingleResponse(c, "OK", users)
}

func (us *UsersController) searchHandler(c *gin.Context) {
	filter := make(map[string]interface{})

	// Iterate over the query parameters
	for key, values := range c.Request.URL.Query() {
		// Skip if the key is empty or has multiple values
		if key == "" || len(values) != 1 {
			continue
		}

		// Add key-value pair to the filter
		filter[key] = values[0]
	}
	users, err := us.useCase.SearchBy(filter)
	if err != nil {
		us.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	us.NewSuccessSingleResponse(c, "OK", users)
}

func (us *UsersController) updateHandler(c *gin.Context) {
	var payload model.User
	if err := us.ParseRequestBody(c, &payload); err != nil {
		us.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// if it's not super admin, user only can update their user detail
	userTyped := utils.AccessInsideToken(us.BaseApi, c)
	if payload.ID != userTyped.UserId && userTyped.Role != "super" {
		us.NewFailedResponse(c, http.StatusForbidden, "access denied")
		return
	}

	if err := us.useCase.UpdateData(&payload); err != nil {
		us.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	us.NewSuccessSingleResponse(c, "OK", payload)
}

func (us *UsersController) deleteHandler(c *gin.Context) {
	id := c.Param("id")
	err := us.useCase.DeleteData(id)
	if err != nil {
		us.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "Success Delete")
}

func NewUsersController(r *gin.Engine, useCase usecase.UsersUseCase, tokenMdw middleware.AuthTokenMiddlerware) *UsersController {
	controller := &UsersController{
		router:  r,
		useCase: useCase,
	}
	usersGroup := r.Group("/users", tokenMdw.RequireToken())
	{
		usersGroup.GET("", controller.listHandler)
		usersGroup.GET("/:id", controller.getHandler)
		usersGroup.GET("/search", controller.searchHandler)
		usersGroup.POST("", controller.createHandler)
		usersGroup.PUT("", controller.updateHandler)
		usersGroup.DELETE("/:id", controller.deleteHandler)
	}
	return controller
}
