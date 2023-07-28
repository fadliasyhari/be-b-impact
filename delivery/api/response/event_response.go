package response

import (
	"time"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"
)

func MapEventToResponse(event *dto.EventDTO) dto.EventDTOResponse {
	res := dto.EventDTOResponse{
		ID:               event.ID,
		Title:            event.Title,
		Description:      event.Description,
		Location:         event.Location,
		StartDate:        event.StartDate,
		EndDate:          event.EndDate,
		Status:           event.Status,
		TotalParticipant: event.TotalParticipant,
		Category:         event.Category.Name,
		CategoryDetail: dto.CategoryDTO{
			ID:   event.Category.ID,
			Name: event.Category.Name,
		},
		CreatedAt: event.CreatedAt,
	}
	for _, v := range event.EventImage {
		res.ImageURLs = append(res.ImageURLs, dto.ImageDTO{
			ID:       v.ID,
			ImageURL: v.ImageURL,
		})
	}
	return res
}

func MapEventToSingleResponse(event *model.Event, total int64, eventParticipant []model.EventParticipant) dto.EventDTOResponse {
	regisStatus := false
	eventParticipantId := ""
	var participantDetail dto.EventParticipantDto
	if len(eventParticipant) > 0 {
		regisStatus = true
		eventParticipantId = eventParticipant[0].ID
		participantDetail.ID = eventParticipant[0].ID
		participantDetail.Name = eventParticipant[0].Name
		participantDetail.Email = eventParticipant[0].Email
		participantDetail.Phone = eventParticipant[0].Phone
		participantDetail.UserID = eventParticipant[0].UserID
		participantDetail.CreatedAt = eventParticipant[0].CreatedAt
	}
	res := dto.EventDTOResponse{
		ID:                event.ID,
		Title:             event.Title,
		Description:       event.Description,
		Location:          event.Location,
		StartDate:         event.StartDate,
		EndDate:           event.EndDate,
		Status:            event.Status,
		IsJoined:          regisStatus,
		ParticipantID:     eventParticipantId,
		ParticipantDetail: participantDetail,
		Category:          event.Category.Name,
		CategoryDetail:    dto.CategoryDTO{ID: event.Category.ID, Name: event.Category.Name},
		ImageURLs:         []dto.ImageDTO{},
		CreatedAt:         event.CreatedAt,
		UpdatedAt:         time.Time{},
		TotalParticipant:  int(total),
	}
	for _, v := range event.EventImage {
		res.ImageURLs = append(res.ImageURLs, dto.ImageDTO{
			ID:       v.ID,
			ImageURL: v.ImageURL,
		})
	}
	return res
}

func MapEventParticipantToResponse(eventParticipant *model.EventParticipant) dto.EventParticipantDto {
	res := dto.EventParticipantDto{
		ID:        eventParticipant.ID,
		UserID:    eventParticipant.UserID,
		Name:      eventParticipant.Name,
		Email:     eventParticipant.Email,
		Phone:     eventParticipant.Phone,
		CreatedAt: eventParticipant.CreatedAt,
	}
	return res
}
