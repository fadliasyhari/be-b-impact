package response

import (
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

func MapEventToSingleResponse(event *model.Event, total int64) dto.EventDTOResponse {
	res := dto.EventDTOResponse{
		ID:               event.ID,
		Title:            event.Title,
		Description:      event.Description,
		Location:         event.Location,
		StartDate:        event.StartDate,
		EndDate:          event.EndDate,
		Status:           event.Status,
		TotalParticipant: int(total),
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

func MapEventParticipantToResponse(eventParticipant *model.EventParticipant) dto.EventParticipantDto {
	res := dto.EventParticipantDto{
		ID:        eventParticipant.ID,
		Name:      eventParticipant.Name,
		Email:     eventParticipant.Email,
		Phone:     eventParticipant.Phone,
		CreatedAt: eventParticipant.CreatedAt,
	}
	return res
}
