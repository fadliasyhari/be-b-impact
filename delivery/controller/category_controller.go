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

type CategoryController struct {
	router  *gin.Engine
	useCase usecase.CategoryUseCase
	api.BaseApi
}

func (ca *CategoryController) createHandler(c *gin.Context) {
	// only admin and superadmin who can create category
	userTyped := utils.AccessInsideToken(ca.BaseApi, c)
	if userTyped.Role != "admin" && userTyped.Role != "super" {
		ca.NewFailedResponse(c, http.StatusForbidden, "access denied")
		return
	}

	var payload model.Category
	if err := ca.ParseRequestBody(c, &payload); err != nil {
		ca.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if payload.Parent == "" {
		payload.Parent = "0"
	}
	if payload.Status == "" {
		payload.Status = "1"
	}
	payload.CreatedBy = userTyped.UserId

	if err := ca.useCase.SaveData(&payload); err != nil {
		ca.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	ca.NewSuccessSingleResponse(c, "OK", payload)
}

func (ca *CategoryController) listHandler(c *gin.Context) {
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
		ca.NewFailedResponse(c, http.StatusBadRequest, "invalid page number")
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		ca.NewFailedResponse(c, http.StatusBadRequest, "invalid limit number")
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
	category, paging, err := ca.useCase.Pagination(requestQueryParams)
	if err != nil {
		ca.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	var categoryInterface []interface{}
	for _, ca := range category {
		categoryInterface = append(categoryInterface, ca)
	}
	ca.NewSuccessPagedResponse(c, "OK", categoryInterface, paging)
}

func (ca *CategoryController) getHandler(c *gin.Context) {
	id := c.Param("id")
	category, err := ca.useCase.FindById(id)
	if err != nil {
		ca.NewFailedResponse(c, http.StatusNotFound, err.Error())
		return
	}

	ca.NewSuccessSingleResponse(c, "OK", category)
}

func (ca *CategoryController) searchHandler(c *gin.Context) {
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
	category, err := ca.useCase.SearchBy(filter)
	if err != nil {
		ca.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	ca.NewSuccessSingleResponse(c, "OK", category)
}

func (ca *CategoryController) updateHandler(c *gin.Context) {
	userTyped := utils.AccessInsideToken(ca.BaseApi, c)
	if userTyped.Role != "admin" && userTyped.Role != "super" {
		ca.NewFailedResponse(c, http.StatusForbidden, "access denied")
		return
	}

	var payload model.Category
	if err := ca.ParseRequestBody(c, &payload); err != nil {
		ca.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := ca.useCase.UpdateData(&payload); err != nil {
		ca.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	ca.NewSuccessSingleResponse(c, "OK", payload)
}

func (ca *CategoryController) deleteHandler(c *gin.Context) {
	userTyped := utils.AccessInsideToken(ca.BaseApi, c)
	if userTyped.Role != "admin" && userTyped.Role != "super" {
		ca.NewFailedResponse(c, http.StatusForbidden, "access denied")
		return
	}
	id := c.Param("id")
	err := ca.useCase.DeleteData(id)
	if err != nil {
		ca.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "Success Delete")
}

func NewCategoryController(r *gin.Engine, useCase usecase.CategoryUseCase, tokenMdw middleware.AuthTokenMiddlerware) *CategoryController {
	controller := &CategoryController{
		router:  r,
		useCase: useCase,
	}
	categoryGroup := r.Group("/category")
	{
		categoryGroup.GET("", controller.listHandler)
		categoryGroup.GET("/:id", controller.getHandler)
		categoryGroup.GET("/search", controller.searchHandler)
		categoryGroup.Use(tokenMdw.RequireToken())
		categoryGroup.POST("", controller.createHandler)
		categoryGroup.PUT("", controller.updateHandler)
		categoryGroup.DELETE("/:id", controller.deleteHandler)
	}
	return controller
}
