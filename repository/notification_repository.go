package repository

import (
	"fmt"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"

	"be-b-impact.com/csr/utils/common"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	BaseRepository[model.Notification]
	BaseRepositoryCount[model.Notification]
	BaseRepositoryPaging[model.Notification]
	SaveTrx(payload *model.Notification, tx *gorm.DB) error
	UpdateTrx(payload *model.Notification, tx *gorm.DB) error
}
type notificationRepository struct {
	db *gorm.DB
}

func (no *notificationRepository) Delete(id string) error {
	return no.db.Delete(&model.Notification{}, "id=?", id).Error
}

func (no *notificationRepository) Get(id string) (*model.Notification, error) {
	var notification model.Notification
	result := no.db.First(&notification, "id=?", id).Error
	if result != nil {
		return nil, result
	}
	return &notification, nil
}

func (no *notificationRepository) List() ([]model.Notification, error) {
	var notification []model.Notification
	result := no.db.Find(&notification).Error
	if result != nil {
		return nil, result
	}
	return notification, nil
}

func (no *notificationRepository) Save(payload *model.Notification) error {
	return no.db.Save(payload).Error
}

func (no *notificationRepository) SaveTrx(payload *model.Notification, tx *gorm.DB) error {
	return tx.Create(payload).Error
}

func (no *notificationRepository) Update(payload *model.Notification) error {
	updateFields := make(map[string]interface{})

	// Add fields to be updated based on the payload
	if payload.Body != "" {
		updateFields["body"] = payload.Body
	}

	if payload.Type != "" {
		updateFields["type"] = payload.Type
	}

	if payload.TypeID != "" {
		updateFields["type_id"] = payload.TypeID
	}

	return no.db.Model(&model.Notification{}).Where("id = ?", payload.ID).Updates(updateFields).Error
}

func (no *notificationRepository) UpdateTrx(payload *model.Notification, tx *gorm.DB) error {
	updateFields := make(map[string]interface{})

	// Add fields to be updated based on the payload
	if payload.Body != "" {
		updateFields["body"] = payload.Body
	}

	if payload.Type != "" {
		updateFields["type"] = payload.Type
	}

	if payload.TypeID != "" {
		updateFields["type_id"] = payload.TypeID
	}

	return tx.Model(&model.Notification{}).Where("id = ?", payload.ID).Updates(updateFields).Error
}

func (no *notificationRepository) Search(by map[string]interface{}) ([]model.Notification, error) {
	var notification []model.Notification
	query := no.db
	for key, value := range by {
		// Perform case-insensitive search using ilike
		query = query.Where(fmt.Sprintf("%s ilike ?", key), fmt.Sprintf("%%%v%%", value))
	}
	result := query.Find(&notification).Error
	if result != nil {
		return nil, result
	}
	return notification, nil
}

func (no *notificationRepository) CountData(fieldname string, id string) error {
	var count int64
	var result *gorm.DB

	if id != "" {
		result = no.db.Model(&model.Notification{}).Where("name ilike ? AND id <> ?", fieldname, id).Count(&count)
	} else {
		result = no.db.Model(&model.Notification{}).Where("name ilike ?", fieldname).Count(&count)
	}
	if result.Error != nil {
		return result.Error
	}
	if count > 0 {
		return fmt.Errorf("name %s already exist", fieldname)
	}
	return nil
}

func (no *notificationRepository) Paging(requestQueryParam dto.RequestQueryParams) ([]model.Notification, dto.Paging, error) {
	paginationQuery, orderQuery := pagingValidate(requestQueryParam)

	var notification []model.Notification
	query := no.db
	for key, value := range requestQueryParam.Filter {
		// Perform case-insensitive search using ilike
		query = query.Where(fmt.Sprintf("%s ilike ?", key), fmt.Sprintf("%%%v%%", value))
	}
	query = query.Or("user_id = ?", "0")

	// Add condition for user_id = 0

	var totalRows int64
	err := query.Model(model.Notification{}).Count(&totalRows).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}
	err = query.Order(orderQuery).Limit(paginationQuery.Take).Offset(paginationQuery.Skip).Find(&notification).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}
	return notification, common.Paginate(paginationQuery.Page, paginationQuery.Take, int(totalRows)), nil
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}
