package usecase

import (
	"fmt"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"
	"be-b-impact.com/csr/repository"
	"gorm.io/gorm"
)

type NotificationUseCase interface {
	BaseUseCase[model.Notification]
	BaseUseCasePaging[model.Notification]
	SaveNotifDetail(payload *model.Notification, tx *gorm.DB) error
	UpdateNotifDetail(payload *model.Notification, tx *gorm.DB) error
	DeleteNotif(types string, typeId string) error
}

type notificationUseCase struct {
	repo repository.NotificationRepository
}

// DeleteNotif implements NotificationUseCase.
func (no *notificationUseCase) DeleteNotif(types string, typeId string) error {
	return no.repo.DeleteNotif(types, typeId)
}

func (no *notificationUseCase) DeleteData(id string) error {
	notification, err := no.FindById(id)
	if err != nil {
		return fmt.Errorf("notification with ID %s not found", id)
	}
	return no.repo.Delete(notification.ID)
}

func (no *notificationUseCase) FindAll() ([]model.Notification, error) {
	return no.repo.List()
}

func (no *notificationUseCase) FindById(id string) (*model.Notification, error) {
	notification, err := no.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("notification with ID %s not found", id)
	}
	return notification, nil
}

func (no *notificationUseCase) SaveData(payload *model.Notification) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }
	// cek jika data sudah ada -> count > 0

	if payload.ID != "" {
		_, err := no.FindById(payload.ID)
		if err != nil {
			return fmt.Errorf("notification with ID %s not found", payload.ID)
		}
	}
	return no.repo.Save(payload)
}

func (no *notificationUseCase) SaveNotifDetail(payload *model.Notification, tx *gorm.DB) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }

	return no.repo.SaveTrx(payload, tx)
}

func (no *notificationUseCase) UpdateNotifDetail(payload *model.Notification, tx *gorm.DB) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }

	return no.repo.UpdateTrx(payload, tx)
}

func (no *notificationUseCase) UpdateData(payload *model.Notification) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }
	// cek jika data sudah ada -> count > 0

	if payload.ID != "" {
		_, err := no.FindById(payload.ID)
		if err != nil {
			return fmt.Errorf("notification with ID %s not found", payload.ID)
		}
	}
	return no.repo.Update(payload)
}

func (no *notificationUseCase) SearchBy(by map[string]interface{}) ([]model.Notification, error) {
	notifications, err := no.repo.Search(by)
	if err != nil {
		return nil, fmt.Errorf("request invalid")
	}
	return notifications, nil
}

func (no *notificationUseCase) Pagination(requestQueryParams dto.RequestQueryParams) ([]model.Notification, dto.Paging, error) {
	if !requestQueryParams.QueryParams.IsSortValid() {
		return nil, dto.Paging{}, fmt.Errorf("invalid sort by: %s", requestQueryParams.QueryParams.Sort)
	}
	return no.repo.Paging(requestQueryParams)
}

func NewNotificationUseCase(repo repository.NotificationRepository) NotificationUseCase {
	return &notificationUseCase{repo: repo}
}
