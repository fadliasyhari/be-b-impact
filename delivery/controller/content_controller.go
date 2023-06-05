package controller

import (
	"encoding/json"
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

type ContentController struct {
	router        *gin.Engine
	useCase       usecase.ContentUseCase
	imageUC       usecase.ImageUseCase
	tagsContentUC usecase.TagsContentUseCase
	api.BaseApi
}

func (co *ContentController) createHandler(c *gin.Context) {
	userTyped := utils.AccessInsideToken(co.BaseApi, c)
	if userTyped.Role != "admin" && userTyped.Role != "super" {
		co.NewFailedResponse(c, http.StatusForbidden, "access denied")
		return
	}

	// Parse the form data
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		co.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get the form values
	title := c.Request.FormValue("title")
	body := c.Request.FormValue("body")
	status := c.Request.FormValue("status")
	author := c.Request.FormValue("author")
	excerpt := c.Request.FormValue("excerpt")
	categoryID := c.Request.FormValue("category_id")
	tagsString := c.Request.FormValue("tags")

	// Create the content payload
	payload := model.Content{
		Title:      title,
		Body:       body,
		Status:     status,
		Author:     author,
		Excerpt:    excerpt,
		CategoryID: categoryID,
		CreatedBy:  userTyped.UserId,
	}

	if err := co.useCase.SaveData(&payload); err != nil {
		co.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var tags []string
	if err := json.Unmarshal([]byte(tagsString), &tags); err != nil {
		co.NewFailedResponse(c, http.StatusBadRequest, "invalid tags format")
		return
	}

	for _, v := range tags {
		tcPayload := model.TagsContent{
			TagID:     v,
			ContentID: payload.ID,
		}
		if err := co.tagsContentUC.SaveData(&tcPayload); err != nil {
			co.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	file, _, err := c.Request.FormFile("images")
	if err == nil {
		imageURL, err := co.imageUC.FirebaseUpload(file)
		if err != nil {
			co.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		// Create the image payload
		imagePayload := model.Image{
			ImageURL:  imageURL,
			ContentID: payload.ID,
		}
		if err := co.imageUC.SaveData(&imagePayload); err != nil {
			co.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	co.NewSuccessSingleResponse(c, "OK", payload)
}

func (co *ContentController) listHandler(c *gin.Context) {
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
		co.NewFailedResponse(c, http.StatusBadRequest, "invalid page number")
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		co.NewFailedResponse(c, http.StatusBadRequest, "invalid limit number")
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
	content, paging, err := co.useCase.Pagination(requestQueryParams)
	if err != nil {
		co.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	var contentInterface []interface{}
	for _, co := range content {
		contentInterface = append(contentInterface, co)
	}
	co.NewSuccessPagedResponse(c, "OK", contentInterface, paging)
}

func (co *ContentController) getHandler(c *gin.Context) {
	id := c.Param("id")
	content, err := co.useCase.FindById(id)
	if err != nil {
		co.NewFailedResponse(c, http.StatusNotFound, err.Error())
		return
	}

	co.NewSuccessSingleResponse(c, "OK", content)
}

func (co *ContentController) searchHandler(c *gin.Context) {
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
	content, err := co.useCase.SearchBy(filter)
	if err != nil {
		co.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	co.NewSuccessSingleResponse(c, "OK", content)
}

func (co *ContentController) updateHandler(c *gin.Context) {
	userTyped := utils.AccessInsideToken(co.BaseApi, c)
	if userTyped.Role != "admin" && userTyped.Role != "super" {
		co.NewFailedResponse(c, http.StatusForbidden, "access denied")
		return
	}

	// Parse the form data
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		co.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get the form values
	id := c.Request.FormValue("id")
	title := c.Request.FormValue("title")
	body := c.Request.FormValue("body")
	status := c.Request.FormValue("status")
	author := c.Request.FormValue("author")
	excerpt := c.Request.FormValue("excerpt")
	categoryID := c.Request.FormValue("category_id")

	existingContent, err := co.useCase.FindById(id)
	if err != nil {
		co.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	existingContent.Title = title
	existingContent.Body = body
	existingContent.Status = status
	existingContent.Author = author
	existingContent.Excerpt = excerpt
	existingContent.CategoryID = categoryID

	if err := co.useCase.UpdateData(existingContent); err != nil {
		co.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	tagsString := c.Request.FormValue("tags")
	var tags []string
	if err := json.Unmarshal([]byte(tagsString), &tags); err != nil {
		co.NewFailedResponse(c, http.StatusBadRequest, "invalid tags format")
		return
	}

	if len(tags) > 0 && len(existingContent.TagsContent) > 0 {
		for _, tag := range existingContent.TagsContent {
			if err := co.tagsContentUC.DeleteData(tag.ID); err != nil {
				co.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
				return
			}
		}
	}
	for _, v := range tags {
		tcPayload := model.TagsContent{
			TagID:     v,
			ContentID: existingContent.ID,
		}
		if err := co.tagsContentUC.SaveData(&tcPayload); err != nil {
			co.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	file, _, err := c.Request.FormFile("images")
	if err == nil {
		if len(existingContent.Image) > 0 {
			for _, v := range existingContent.Image {
				if err := co.imageUC.DeleteData(v.ID); err != nil {
					co.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
					return
				}
			}
		}
		imageURL, err := co.imageUC.FirebaseUpload(file)
		if err != nil {
			co.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		// Create the image payload
		imagePayload := model.Image{
			ImageURL:  imageURL,
			ContentID: existingContent.ID,
		}
		if err := co.imageUC.SaveData(&imagePayload); err != nil {
			co.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	co.NewSuccessSingleResponse(c, "OK", existingContent)
}

func (co *ContentController) deleteHandler(c *gin.Context) {
	userTyped := utils.AccessInsideToken(co.BaseApi, c)
	if userTyped.Role != "admin" && userTyped.Role != "super" {
		co.NewFailedResponse(c, http.StatusForbidden, "access denied")
		return
	}
	id := c.Param("id")
	err := co.useCase.DeleteData(id)
	if err != nil {
		co.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "Success Delete")
}

func NewContentController(r *gin.Engine, useCase usecase.ContentUseCase, imageUC usecase.ImageUseCase, tagsContentUC usecase.TagsContentUseCase, tokenMdw middleware.AuthTokenMiddlerware) *ContentController {
	controller := &ContentController{
		router:        r,
		useCase:       useCase,
		imageUC:       imageUC,
		tagsContentUC: tagsContentUC,
	}
	contentGroup := r.Group("/content")
	{
		contentGroup.GET("", controller.listHandler)
		contentGroup.GET("/:id", controller.getHandler)
		contentGroup.GET("/search", controller.searchHandler)
		contentGroup.Use(tokenMdw.RequireToken())
		contentGroup.POST("", controller.createHandler)
		contentGroup.PUT("", controller.updateHandler)
		contentGroup.DELETE("/:id", controller.deleteHandler)
	}
	return controller
}
