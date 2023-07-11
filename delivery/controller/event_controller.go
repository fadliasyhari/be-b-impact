package controller

import (
	"net/http"
	"strconv"

	"be-b-impact.com/csr/delivery/api"
	"be-b-impact.com/csr/delivery/api/middleware"
	"be-b-impact.com/csr/delivery/api/response"
	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"
	"be-b-impact.com/csr/usecase"
	"be-b-impact.com/csr/utils"
	"github.com/gin-gonic/gin"
)

type EventController struct {
	router             *gin.Engine
	useCase            usecase.EventUseCase
	eventParticipantUC usecase.EventParticipantUseCase
	api.BaseApi
}

func (ev *EventController) createHandler(c *gin.Context) {
	userTyped := utils.AccessInsideToken(ev.BaseApi, c)
	if userTyped.Role != "admin" && userTyped.Role != "super" {
		ev.NewFailedResponse(c, http.StatusForbidden, "access denied")
		return
	}

	// Parse the form data
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		ev.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get the form values
	title := c.Request.FormValue("title")
	description := c.Request.FormValue("description")
	location := c.Request.FormValue("location")
	startDate := c.Request.FormValue("start_date")
	endDate := c.Request.FormValue("end_date")
	categoryID := c.Request.FormValue("category_id")
	status := c.Request.FormValue("status")
	file, _, err := c.Request.FormFile("images")
	if err != nil {
		ev.NewFailedResponse(c, http.StatusBadRequest, "image not valid")
	}

	// Create the event payload
	payload := model.Event{
		Title:       title,
		Description: description,
		StartDate:   startDate,
		EndDate:     endDate,
		Location:    location,
		Status:      status,
		CategoryID:  categoryID,
	}

	if err := ev.useCase.SaveEvent(&payload, file); err != nil {
		ev.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	ev.NewSuccessSingleResponse(c, "OK", payload)
}

func (ev *EventController) listHandler(c *gin.Context) {
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
		ev.NewFailedResponse(c, http.StatusBadRequest, "invalid page number")
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		ev.NewFailedResponse(c, http.StatusBadRequest, "invalid limit number")
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
	event, paging, err := ev.useCase.PaginationDto(requestQueryParams)
	if err != nil {
		ev.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	var eventInterface []interface{}
	for _, ev := range event {
		res := response.MapEventToResponse(&ev)
		eventInterface = append(eventInterface, res)
	}
	ev.NewSuccessPagedResponse(c, "OK", eventInterface, paging)
}

func (ev *EventController) getHandler(c *gin.Context) {
	id := c.Param("id")
	event, err := ev.useCase.FindById(id)
	if err != nil {
		ev.NewFailedResponse(c, http.StatusNotFound, err.Error())
		return
	}

	total_participant, _ := ev.eventParticipantUC.CountParticipant(id)

	res := response.MapEventToSingleResponse(event, total_participant)

	ev.NewSuccessSingleResponse(c, "OK", res)
}

func (ev *EventController) searchHandler(c *gin.Context) {
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
	event, err := ev.useCase.SearchBy(filter)
	if err != nil {
		ev.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	ev.NewSuccessSingleResponse(c, "OK", event)
}

func (ev *EventController) updateHandler(c *gin.Context) {
	userTyped := utils.AccessInsideToken(ev.BaseApi, c)
	if userTyped.Role != "admin" && userTyped.Role != "super" {
		ev.NewFailedResponse(c, http.StatusForbidden, "access denied")
		return
	}

	// Parse the form data
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		ev.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get the form values ok
	id := c.Request.FormValue("id")
	title := c.Request.FormValue("title")
	description := c.Request.FormValue("description")
	location := c.Request.FormValue("location")
	startDate := c.Request.FormValue("start_date")
	endDate := c.Request.FormValue("end_date")
	categoryID := c.Request.FormValue("category_id")
	status := c.Request.FormValue("status")

	existingEvent, err := ev.useCase.FindById(id)
	if err != nil {
		ev.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	existingEvent.Title = title
	existingEvent.Description = description
	existingEvent.Status = status
	existingEvent.StartDate = startDate
	existingEvent.EndDate = endDate
	existingEvent.Location = location
	existingEvent.CategoryID = categoryID

	file, _, err := c.Request.FormFile("images")
	if err != nil {
		ev.NewFailedResponse(c, http.StatusBadRequest, "image not valid")
	}

	if err := ev.useCase.UpdateEvent(existingEvent, file); err != nil {
		ev.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	ev.NewSuccessSingleResponse(c, "OK", existingEvent)
}

func (ev *EventController) deleteHandler(c *gin.Context) {
	userTyped := utils.AccessInsideToken(ev.BaseApi, c)
	if userTyped.Role != "admin" && userTyped.Role != "super" {
		ev.NewFailedResponse(c, http.StatusForbidden, "access denied")
		return
	}
	id := c.Param("id")
	err := ev.useCase.DeleteData(id)
	if err != nil {
		ev.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "Success Delete")
}

func (ev *EventController) joinHandler(c *gin.Context) {
	userTyped := utils.AccessInsideToken(ev.BaseApi, c)

	var payload model.EventParticipant
	if err := ev.ParseRequestBody(c, &payload); err != nil {
		ev.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	payload.UserID = userTyped.UserId

	if err := ev.eventParticipantUC.SaveData(&payload); err != nil {
		ev.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	ev.NewSuccessSingleResponse(c, "OK", payload)
}

func (ev *EventController) cancelJoinHandler(c *gin.Context) {
	id := c.Param("id")
	err := ev.eventParticipantUC.DeleteData(id)
	if err != nil {
		ev.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	var data interface{}
	ev.NewSuccessSingleResponse(c, "Successfully canceled your event participation", data)
}

func (ev *EventController) getParticipantHandler(c *gin.Context) {
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
		ev.NewFailedResponse(c, http.StatusBadRequest, "invalid page number")
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		ev.NewFailedResponse(c, http.StatusBadRequest, "invalid limit number")
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
	eventParticipants, paging, err := ev.eventParticipantUC.Pagination(requestQueryParams)
	if err != nil {
		ev.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	var eventParticipantInterface []interface{}
	for _, ep := range eventParticipants {
		res := response.MapEventParticipantToResponse(&ep)
		eventParticipantInterface = append(eventParticipantInterface, res)
	}
	ev.NewSuccessPagedResponse(c, "OK", eventParticipantInterface, paging)
}

func NewEventController(r *gin.Engine, useCase usecase.EventUseCase, eventParticipantUC usecase.EventParticipantUseCase, tokenMdw middleware.AuthTokenMiddlerware) *EventController {
	controller := &EventController{
		router:             r,
		useCase:            useCase,
		eventParticipantUC: eventParticipantUC,
	}
	eventGroup := r.Group("/event")
	{
		eventGroup.GET("", controller.listHandler)
		eventGroup.GET("/:id", controller.getHandler)
		eventGroup.GET("/participants", controller.getParticipantHandler)
		eventGroup.GET("/search", controller.searchHandler)
		eventGroup.Use(tokenMdw.RequireToken())
		eventGroup.POST("/join", controller.joinHandler)
		eventGroup.POST("/cancel-join/:id", controller.cancelJoinHandler)
		eventGroup.POST("", controller.createHandler)
		eventGroup.PUT("", controller.updateHandler)
		eventGroup.DELETE("/:id", controller.deleteHandler)
	}

	return controller
}
