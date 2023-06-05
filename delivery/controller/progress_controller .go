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

type ProgressController struct {
	router  *gin.Engine
	useCase usecase.ProgressUseCase
	api.BaseApi
}

func (pr *ProgressController) createHandler(c *gin.Context) {

	var payload model.Progress
	if err := pr.ParseRequestBody(c, &payload); err != nil {
		pr.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := pr.useCase.SaveData(&payload); err != nil {
		pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	pr.NewSuccessSingleResponse(c, "OK", payload)
}

func (pr *ProgressController) listHandler(c *gin.Context) {
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
		pr.NewFailedResponse(c, http.StatusBadRequest, "invalid page number")
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		pr.NewFailedResponse(c, http.StatusBadRequest, "invalid limit number")
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
	progress, paging, err := pr.useCase.Pagination(requestQueryParams)
	if err != nil {
		pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	var tagInterface []interface{}
	for _, pr := range progress {
		tagInterface = append(tagInterface, pr)
	}
	pr.NewSuccessPagedResponse(c, "OK", tagInterface, paging)
}

func (pr *ProgressController) getHandler(c *gin.Context) {
	id := c.Param("id")
	progress, err := pr.useCase.FindById(id)
	if err != nil {
		pr.NewFailedResponse(c, http.StatusNotFound, err.Error())
		return
	}

	pr.NewSuccessSingleResponse(c, "OK", progress)
}

func (pr *ProgressController) searchHandler(c *gin.Context) {
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
	progress, err := pr.useCase.SearchBy(filter)
	if err != nil {
		pr.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	pr.NewSuccessSingleResponse(c, "OK", progress)
}

func (pr *ProgressController) updateHandler(c *gin.Context) {

	var payload model.Progress
	if err := pr.ParseRequestBody(c, &payload); err != nil {
		pr.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := pr.useCase.UpdateData(&payload); err != nil {
		pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	pr.NewSuccessSingleResponse(c, "OK", payload)
}

func (pr *ProgressController) deleteHandler(c *gin.Context) {
	userTyped := utils.AccessInsideToken(pr.BaseApi, c)
	if userTyped.Role != "admin" && userTyped.Role != "super" {
		pr.NewFailedResponse(c, http.StatusForbidden, "access denied")
		return
	}
	id := c.Param("id")
	err := pr.useCase.DeleteData(id)
	if err != nil {
		pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "Success Delete")
}

func NewProgressController(r *gin.Engine, useCase usecase.ProgressUseCase, tokenMdw middleware.AuthTokenMiddlerware) *ProgressController {
	controller := &ProgressController{
		router:  r,
		useCase: useCase,
	}
	tagGroup := r.Group("/progress", tokenMdw.RequireToken())
	{
		tagGroup.GET("", controller.listHandler)
		tagGroup.GET("/:id", controller.getHandler)
		tagGroup.GET("/search", controller.searchHandler)
		tagGroup.POST("", controller.createHandler)
		tagGroup.PUT("", controller.updateHandler)
		tagGroup.DELETE("/:id", controller.deleteHandler)
	}
	return controller
}
