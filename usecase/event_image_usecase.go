package usecase

import (
	"fmt"
	"mime/multipart"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"
	"be-b-impact.com/csr/repository"
	"gorm.io/gorm"
)

type EventImageUseCase interface {
	BaseUseCase[model.EventImage]
	BaseUseCasePaging[model.EventImage]
	FirebaseUpload(file multipart.File) (string, error)
	SaveEventImage(payload *model.EventImage, tx *gorm.DB) error
	DeleteDataTrx(id string, tx *gorm.DB) error
}

type eventImageUseCase struct {
	repo repository.EventImageRepository
}

func (im *eventImageUseCase) DeleteData(id string) error {
	eventImage, err := im.FindById(id)
	if err != nil {
		return fmt.Errorf("eventImage with ID %s not found", id)
	}
	return im.repo.Delete(eventImage.ID)
}

func (im *eventImageUseCase) DeleteDataTrx(id string, tx *gorm.DB) error {
	eventImage, err := im.FindById(id)
	if err != nil {
		return fmt.Errorf("eventImage with ID %s not found", id)
	}
	return im.repo.DeleteTrx(eventImage.ID, tx)
}

func (im *eventImageUseCase) FindAll() ([]model.EventImage, error) {
	return im.repo.List()
}

func (im *eventImageUseCase) FindById(id string) (*model.EventImage, error) {
	eventImage, err := im.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("eventImage with ID %s not found", id)
	}
	return eventImage, nil
}

func (im *eventImageUseCase) FirebaseUpload(file multipart.File) (string, error) {
	return im.repo.FirebaseSave(file)
}

func (im *eventImageUseCase) SaveData(payload *model.EventImage) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }
	return im.repo.Save(payload)
}

func (im *eventImageUseCase) SaveEventImage(payload *model.EventImage, tx *gorm.DB) error {
	// Save the EventImage using the provided transaction
	if err := im.repo.SaveTrx(payload, tx); err != nil {
		return err
	}
	return nil
}

func (im *eventImageUseCase) UpdateData(payload *model.EventImage) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }
	// cek jika data sudah ada -> count > 0

	if payload.ID != "" {
		_, err := im.FindById(payload.ID)
		if err != nil {
			return fmt.Errorf("eventImage with ID %s not found", payload.ID)
		}
	}
	return im.repo.Update(payload)
}

func (im *eventImageUseCase) SearchBy(by map[string]interface{}) ([]model.EventImage, error) {
	eventImages, err := im.repo.Search(by)
	if err != nil {
		return nil, fmt.Errorf("request invalid")
	}
	return eventImages, nil
}

func (im *eventImageUseCase) Pagination(requestQueryParams dto.RequestQueryParams) ([]model.EventImage, dto.Paging, error) {
	if !requestQueryParams.QueryParams.IsSortValid() {
		return nil, dto.Paging{}, fmt.Errorf("invalid sort by: %s", requestQueryParams.QueryParams.Sort)
	}
	return im.repo.Paging(requestQueryParams)
}

func NewEventImageUseCase(repo repository.EventImageRepository) EventImageUseCase {
	return &eventImageUseCase{repo: repo}
}
