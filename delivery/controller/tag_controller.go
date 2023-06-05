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

type TagController struct {
	router  *gin.Engine
	useCase usecase.TagUseCase
	api.BaseApi
}

func (ta *TagController) createHandler(c *gin.Context) {

	var payload model.Tag
	if err := ta.ParseRequestBody(c, &payload); err != nil {
		ta.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := ta.useCase.SaveData(&payload); err != nil {
		ta.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	ta.NewSuccessSingleResponse(c, "OK", payload)
}

func (ta *TagController) listHandler(c *gin.Context) {
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
		ta.NewFailedResponse(c, http.StatusBadRequest, "invalid page number")
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		ta.NewFailedResponse(c, http.StatusBadRequest, "invalid limit number")
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
	tag, paging, err := ta.useCase.Pagination(requestQueryParams)
	if err != nil {
		ta.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	var tagInterface []interface{}
	for _, ta := range tag {
		tagInterface = append(tagInterface, ta)
	}
	ta.NewSuccessPagedResponse(c, "OK", tagInterface, paging)
}

func (ta *TagController) getHandler(c *gin.Context) {
	id := c.Param("id")
	tag, err := ta.useCase.FindById(id)
	if err != nil {
		ta.NewFailedResponse(c, http.StatusNotFound, err.Error())
		return
	}

	ta.NewSuccessSingleResponse(c, "OK", tag)
}

func (ta *TagController) searchHandler(c *gin.Context) {
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
	tag, err := ta.useCase.SearchBy(filter)
	if err != nil {
		ta.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	ta.NewSuccessSingleResponse(c, "OK", tag)
}

func (ta *TagController) updateHandler(c *gin.Context) {

	var payload model.Tag
	if err := ta.ParseRequestBody(c, &payload); err != nil {
		ta.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := ta.useCase.UpdateData(&payload); err != nil {
		ta.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	ta.NewSuccessSingleResponse(c, "OK", payload)
}

func (ta *TagController) deleteHandler(c *gin.Context) {
	userTyped := utils.AccessInsideToken(ta.BaseApi, c)
	if userTyped.Role != "admin" && userTyped.Role != "super" {
		ta.NewFailedResponse(c, http.StatusForbidden, "access denied")
		return
	}
	id := c.Param("id")
	err := ta.useCase.DeleteData(id)
	if err != nil {
		ta.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "Success Delete")
}

func NewTagController(r *gin.Engine, useCase usecase.TagUseCase, tokenMdw middleware.AuthTokenMiddlerware) *TagController {
	controller := &TagController{
		router:  r,
		useCase: useCase,
	}
	tagGroup := r.Group("/tag")
	{
		tagGroup.GET("", controller.listHandler)
		tagGroup.GET("/:id", controller.getHandler)
		tagGroup.GET("/search", controller.searchHandler)
		tagGroup.Use(tokenMdw.RequireToken())
		tagGroup.POST("", controller.createHandler)
		tagGroup.PUT("", controller.updateHandler)
		tagGroup.DELETE("/:id", controller.deleteHandler)
	}
	return controller
}
