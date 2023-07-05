package usecase

import (
	"fmt"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"
	"be-b-impact.com/csr/repository"
	"gorm.io/gorm"
)

type EventParticipantUseCase interface {
	BaseUseCase[model.EventParticipant]
	BaseUseCasePaging[model.EventParticipant]
	SaveEventParticipant(payload *model.EventParticipant, tx *gorm.DB) error
	DeleteDataTrx(id string, tx *gorm.DB) error
	CountParticipant(id string) (int64, error)
}

type eventParticipantUseCase struct {
	repo repository.EventParticipantRepository
}

// CountParticipant implements EventParticipantUseCase.
func (ep *eventParticipantUseCase) CountParticipant(id string) (int64, error) {
	return ep.repo.CountParticipant(id)
}

func (ep *eventParticipantUseCase) DeleteData(id string) error {
	eventParticipant, err := ep.FindById(id)
	if err != nil {
		return fmt.Errorf("eventParticipant with ID %s not found", id)
	}
	return ep.repo.Delete(eventParticipant.ID)
}

func (ep *eventParticipantUseCase) DeleteDataTrx(id string, tx *gorm.DB) error {
	eventParticipant, err := ep.FindById(id)
	if err != nil {
		return fmt.Errorf("eventParticipant with ID %s not found", id)
	}
	return ep.repo.DeleteTrx(eventParticipant.ID, tx)
}

func (ep *eventParticipantUseCase) FindAll() ([]model.EventParticipant, error) {
	return ep.repo.List()
}

func (ep *eventParticipantUseCase) FindById(id string) (*model.EventParticipant, error) {
	eventParticipant, err := ep.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("eventParticipant with ID %s not found", id)
	}
	return eventParticipant, nil
}

func (ep *eventParticipantUseCase) SaveData(payload *model.EventParticipant) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }
	// cek jika data sudah ada -> count > 0
	filter := make(map[string]interface{})
	filter["event_id"] = payload.EventID
	filter["user_id"] = payload.UserID
	requestQueryParams := dto.RequestQueryParams{
		QueryParams: dto.QueryParams{
			Sort:  "ASC",
			Order: "id",
		},
		PaginationParam: dto.PaginationParam{
			Page:  1,
			Limit: 50,
		},
		Filter: filter,
	}

	res, _, _ := ep.repo.Paging(requestQueryParams)
	if len(res) > 0 {
		return fmt.Errorf("cannot join event multiple times")
	}

	if payload.ID != "" {
		_, err := ep.FindById(payload.ID)
		if err != nil {
			return fmt.Errorf("eventParticipant with ID %s not found", payload.ID)
		}
	}
	return ep.repo.Save(payload)
}

func (ep *eventParticipantUseCase) SaveEventParticipant(payload *model.EventParticipant, tx *gorm.DB) error {
	// Save the EventParticipant using the provided transaction
	if err := ep.repo.SaveTrx(payload, tx); err != nil {
		return err
	}
	return nil
}

func (ep *eventParticipantUseCase) UpdateData(payload *model.EventParticipant) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }
	// cek jika data sudah ada -> count > 0

	if payload.ID != "" {
		_, err := ep.FindById(payload.ID)
		if err != nil {
			return fmt.Errorf("eventParticipant with ID %s not found", payload.ID)
		}
	}
	return ep.repo.Update(payload)
}

func (ep *eventParticipantUseCase) SearchBy(by map[string]interface{}) ([]model.EventParticipant, error) {
	eventParticipants, err := ep.repo.Search(by)
	if err != nil {
		return nil, fmt.Errorf("request invalid")
	}
	return eventParticipants, nil
}

func (ep *eventParticipantUseCase) Pagination(requestQueryParams dto.RequestQueryParams) ([]model.EventParticipant, dto.Paging, error) {
	if !requestQueryParams.QueryParams.IsSortValid() {
		return nil, dto.Paging{}, fmt.Errorf("invalid sort by: %s", requestQueryParams.QueryParams.Sort)
	}
	return ep.repo.Paging(requestQueryParams)
}

func NewEventParticipantUseCase(repo repository.EventParticipantRepository) EventParticipantUseCase {
	return &eventParticipantUseCase{repo: repo}
}
