package usecase

import (
	"fmt"
	"mime/multipart"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"
	"be-b-impact.com/csr/repository"
)

type EventUseCase interface {
	BaseUseCase[model.Event]
	BaseUseCasePaging[model.Event]
	SaveEvent(payload *model.Event, file multipart.File) error
	UpdateEvent(payload *model.Event, file multipart.File) error
	PaginationDto(requestQueryParams dto.RequestQueryParams) ([]dto.EventDTO, dto.Paging, error)
	FindByIdDto(id string) (*dto.EventDTO, error)
}

type eventUseCase struct {
	repo               repository.EventRepository
	eventParticipantUC EventParticipantUseCase
	eventImageUC       EventImageUseCase
	notificationUC     NotificationUseCase
}

// Pagination implements EventUseCase.
func (ev *eventUseCase) Pagination(requestQueryParams dto.RequestQueryParams) ([]model.Event, dto.Paging, error) {
	if !requestQueryParams.QueryParams.IsSortValid() {
		return nil, dto.Paging{}, fmt.Errorf("invalid sort by: %s", requestQueryParams.QueryParams.Sort)
	}
	return ev.repo.Paging(requestQueryParams)
}

// SaveData implements EventUseCase.
func (*eventUseCase) SaveData(payload *model.Event) error {
	panic("unimplemented")
}

func (ev *eventUseCase) DeleteData(id string) error {
	event, err := ev.FindById(id)
	if err != nil {
		return fmt.Errorf("event with ID %s not found", id)
	}

	ev.notificationUC.DeleteNotif("event", event.ID)

	return ev.repo.Delete(event.ID)
}

func (ev *eventUseCase) FindAll() ([]model.Event, error) {
	return ev.repo.List()
}

func (ev *eventUseCase) FindById(id string) (*model.Event, error) {
	event, err := ev.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("event with ID %s not found", id)
	}
	return event, nil
}

func (ev *eventUseCase) FindByIdDto(id string) (*dto.EventDTO, error) {
	event, err := ev.repo.GetDto(id)
	if err != nil {
		return nil, fmt.Errorf("event with ID %s not found", id)
	}
	return event, nil
}

func (ev *eventUseCase) SaveEvent(payload *model.Event, file multipart.File) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }

	if payload.Status == "1" && (payload.Title == "" || payload.Description == "" || payload.StartDate == "" || payload.EndDate == "" || payload.CategoryID == nil || payload.Location == "") {
		return fmt.Errorf("form is not completed")
	}

	tx := ev.repo.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := ev.repo.Save(payload); err != nil {
		tx.Rollback()
		return err
	}

	if file != nil {
		eventImageURL, err := ev.eventImageUC.FirebaseUpload(file)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Create the eventImage payload
		eventImagePayload := model.EventImage{
			ImageURL: eventImageURL,
			EventID:  payload.ID,
		}
		if err := ev.eventImageUC.SaveEventImage(&eventImagePayload, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if payload.Status == "1" {
		notifPayload := model.Notification{
			Body:   fmt.Sprintf("New event: <b>%s</b>", payload.Title),
			UserID: "0",
			Type:   "event",
			TypeID: payload.ID,
		}

		err := ev.notificationUC.SaveNotifDetail(&notifPayload, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (ev *eventUseCase) UpdateData(payload *model.Event) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }
	// cek jika data sudah ada -> count > 0
	if payload.ID != "" {
		_, err := ev.FindById(payload.ID)
		if err != nil {
			return fmt.Errorf("event with ID %s not found", payload.ID)
		}
	}

	return ev.repo.Update(payload)
}

func (ev *eventUseCase) UpdateEvent(payload *model.Event, file multipart.File) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }

	if payload.ID != "" {
		_, err := ev.FindById(payload.ID)
		if err != nil {
			return fmt.Errorf("event with ID %s not found", payload.ID)
		}
	}

	tx := ev.repo.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := ev.repo.Update(payload); err != nil {
		tx.Rollback()
		return err
	}

	if file != nil {

		if len(payload.EventImage) > 0 {
			for _, v := range payload.EventImage {
				if err := ev.eventImageUC.DeleteDataTrx(v.ID, tx); err != nil {
					tx.Rollback()
					return err
				}
			}
		}

		eventImageURL, err := ev.eventImageUC.FirebaseUpload(file)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Create the eventImage payload
		eventImagePayload := model.EventImage{
			ImageURL: eventImageURL,
			EventID:  payload.ID,
		}
		if err := ev.eventImageUC.SaveEventImage(&eventImagePayload, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if payload.Status == "1" && (payload.Title == "" || payload.Description == "" || payload.StartDate == "" || payload.EndDate == "" || payload.CategoryID == nil || payload.Location == "") {
		return fmt.Errorf("form is not completed")
	}

	if payload.Status == "1" {
		notifPayload := model.Notification{
			Body:   fmt.Sprintf("New event: <b>%s</b>", payload.Title),
			UserID: "0",
			Type:   "event",
			TypeID: payload.ID,
		}

		err := ev.notificationUC.SaveNotifDetail(&notifPayload, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (ev *eventUseCase) SearchBy(by map[string]interface{}) ([]model.Event, error) {
	events, err := ev.repo.Search(by)
	if err != nil {
		return nil, fmt.Errorf("request invalid")
	}
	return events, nil
}

func (ev *eventUseCase) PaginationDto(requestQueryParams dto.RequestQueryParams) ([]dto.EventDTO, dto.Paging, error) {
	if !requestQueryParams.QueryParams.IsSortValid() {
		return nil, dto.Paging{}, fmt.Errorf("invalid sort by: %s", requestQueryParams.QueryParams.Sort)
	}

	var eventIdList []string
	for k, v := range requestQueryParams.Filter {
		if k == "user_id" {
			filter := make(map[string]interface{})
			filter["user_id"] = v
			eventList, err := ev.eventParticipantUC.SearchBy(filter)
			if err != nil {
				return nil, dto.Paging{}, fmt.Errorf("event participants not found")
			}

			eventIDMap := make(map[string]bool) // Map to track unique event IDs
			for _, event := range eventList {
				eventIDMap[event.EventID] = true
			}

			for eventID := range eventIDMap {
				eventIdList = append(eventIdList, eventID)
			}
		}
	}

	return ev.repo.PagingDto(requestQueryParams, eventIdList)
}

func NewEventUseCase(repo repository.EventRepository, eventParticipantUC EventParticipantUseCase, eventImageUC EventImageUseCase, notificationUC NotificationUseCase) EventUseCase {
	return &eventUseCase{
		repo:               repo,
		eventParticipantUC: eventParticipantUC,
		eventImageUC:       eventImageUC,
		notificationUC:     notificationUC,
	}
}
