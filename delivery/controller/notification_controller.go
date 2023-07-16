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

type NotificationController struct {
	router  *gin.Engine
	useCase usecase.NotificationUseCase
	api.BaseApi
}

func (no *NotificationController) createHandler(c *gin.Context) {
	// only admin and superadmin who can create notification
	userTyped := utils.AccessInsideToken(no.BaseApi, c)
	if userTyped.Role != "admin" && userTyped.Role != "super" {
		no.NewFailedResponse(c, http.StatusForbidden, "access denied")
		return
	}

	var payload model.Notification
	if err := no.ParseRequestBody(c, &payload); err != nil {
		no.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := no.useCase.SaveData(&payload); err != nil {
		no.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	no.NewSuccessSingleResponse(c, "OK", payload)
}

func (no *NotificationController) listHandler(c *gin.Context) {
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
		no.NewFailedResponse(c, http.StatusBadRequest, "invalid page number")
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		no.NewFailedResponse(c, http.StatusBadRequest, "invalid limit number")
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
	notification, paging, err := no.useCase.Pagination(requestQueryParams)
	if err != nil {
		no.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	var notificationInterface []interface{}
	for _, no := range notification {
		notificationInterface = append(notificationInterface, no)
	}
	no.NewSuccessPagedResponse(c, "OK", notificationInterface, paging)
}

func (no *NotificationController) getHandler(c *gin.Context) {
	id := c.Param("id")
	notification, err := no.useCase.FindById(id)
	if err != nil {
		no.NewFailedResponse(c, http.StatusNotFound, err.Error())
		return
	}

	no.NewSuccessSingleResponse(c, "OK", notification)
}

func (no *NotificationController) searchHandler(c *gin.Context) {
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
	notification, err := no.useCase.SearchBy(filter)
	if err != nil {
		no.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	no.NewSuccessSingleResponse(c, "OK", notification)
}

func (no *NotificationController) updateHandler(c *gin.Context) {
	userTyped := utils.AccessInsideToken(no.BaseApi, c)
	if userTyped.Role != "admin" && userTyped.Role != "super" {
		no.NewFailedResponse(c, http.StatusForbidden, "access denied")
		return
	}

	var payload model.Notification
	if err := no.ParseRequestBody(c, &payload); err != nil {
		no.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := no.useCase.UpdateData(&payload); err != nil {
		no.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	no.NewSuccessSingleResponse(c, "OK", payload)
}

func (no *NotificationController) deleteHandler(c *gin.Context) {
	userTyped := utils.AccessInsideToken(no.BaseApi, c)
	if userTyped.Role != "admin" && userTyped.Role != "super" {
		no.NewFailedResponse(c, http.StatusForbidden, "access denied")
		return
	}
	id := c.Param("id")
	err := no.useCase.DeleteData(id)
	if err != nil {
		no.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "Success Delete")
}

func NewNotificationController(r *gin.Engine, useCase usecase.NotificationUseCase, tokenMdw middleware.AuthTokenMiddlerware) *NotificationController {
	controller := &NotificationController{
		router:  r,
		useCase: useCase,
	}
	notificationGroup := r.Group("/notification")
	{
		notificationGroup.GET("", controller.listHandler)
		notificationGroup.GET("/:id", controller.getHandler)
		notificationGroup.GET("/search", controller.searchHandler)
		notificationGroup.Use(tokenMdw.RequireToken())
		notificationGroup.POST("", controller.createHandler)
		notificationGroup.PUT("", controller.updateHandler)
		notificationGroup.DELETE("/:id", controller.deleteHandler)
	}
	return controller
}
