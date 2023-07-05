package repository

import (
	"fmt"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"

	"be-b-impact.com/csr/utils/common"
	"gorm.io/gorm"
)

type EventParticipantRepository interface {
	BaseRepository[model.EventParticipant]
	BaseRepositoryCount[model.EventParticipant]
	BaseRepositoryPaging[model.EventParticipant]
	SaveTrx(payload *model.EventParticipant, tx *gorm.DB) error
	DeleteTrx(id string, tx *gorm.DB) error
	CountParticipant(id string) (int64, error)
}
type eventParticipantRepository struct {
	db *gorm.DB
}

// CountParticipant implements EventParticipantRepository.
func (ep *eventParticipantRepository) CountParticipant(id string) (int64, error) {
	var count int64
	ep.db.Model(&model.EventParticipant{}).Where("event_id = ?", id).Count(&count)
	return count, nil
}

func (ep *eventParticipantRepository) Delete(id string) error {
	return ep.db.Delete(&model.EventParticipant{}, "id=?", id).Error
}

func (ep *eventParticipantRepository) DeleteTrx(id string, tx *gorm.DB) error {
	return tx.Delete(&model.EventParticipant{}, "id=?", id).Error
}

func (ep *eventParticipantRepository) Get(id string) (*model.EventParticipant, error) {
	var eventParticipant model.EventParticipant
	result := ep.db.First(&eventParticipant, "id=?", id).Error
	if result != nil {
		return nil, result
	}
	return &eventParticipant, nil
}

func (ep *eventParticipantRepository) List() ([]model.EventParticipant, error) {
	var eventParticipant []model.EventParticipant
	result := ep.db.Find(&eventParticipant).Error
	if result != nil {
		return nil, result
	}
	return eventParticipant, nil
}

func (ep *eventParticipantRepository) Save(payload *model.EventParticipant) error {
	return ep.db.Save(payload).Error
}

func (r *eventParticipantRepository) SaveTrx(payload *model.EventParticipant, tx *gorm.DB) error {
	// If the provided transaction is not nil, use it for saving the EventParticipant
	if tx != nil {
		return tx.Create(payload).Error
	}

	// Otherwise, use the default DB connection for saving the EventParticipant
	return r.db.Create(payload).Error
}

func (ep *eventParticipantRepository) Update(payload *model.EventParticipant) error {
	return ep.db.Model(&model.EventParticipant{}).Where("id = ?", payload.EventID).Updates(payload).Error
}

func (ep *eventParticipantRepository) Search(by map[string]interface{}) ([]model.EventParticipant, error) {
	var eventParticipant []model.EventParticipant
	result := ep.db.Where(by).Find(&eventParticipant).Error
	if result != nil {
		return nil, result
	}
	return eventParticipant, nil
}

func (ep *eventParticipantRepository) CountData(fieldname string, id string) error {
	var count int64
	var result *gorm.DB

	if id != "" {
		result = ep.db.Model(&model.EventParticipant{}).Where("name ilike ? AND id <> ?", fieldname, id).Count(&count)
	} else {
		result = ep.db.Model(&model.EventParticipant{}).Where("name ilike ?", fieldname).Count(&count)
	}
	if result.Error != nil {
		return result.Error
	}
	if count > 0 {
		return fmt.Errorf("name %s already exist", fieldname)
	}
	return nil
}

func (ep *eventParticipantRepository) Paging(requestQueryParam dto.RequestQueryParams) ([]model.EventParticipant, dto.Paging, error) {
	paginationQuery, orderQuery := pagingValidate(requestQueryParam)

	var eventParticipant []model.EventParticipant
	query := ep.db
	for key, value := range requestQueryParam.Filter {
		// Perform case-insensitive search using ilike
		query = query.Where(fmt.Sprintf("%s ilike ?", key), fmt.Sprintf("%%%v%%", value))

	}
	var totalRows int64
	err := query.Model(model.EventParticipant{}).Count(&totalRows).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}
	err = query.Order(orderQuery).Limit(paginationQuery.Take).Offset(paginationQuery.Skip).Find(&eventParticipant).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}
	return eventParticipant, common.Paginate(paginationQuery.Page, paginationQuery.Take, int(totalRows)), nil
}

func NewEventParticipantRepository(db *gorm.DB) EventParticipantRepository {
	return &eventParticipantRepository{db: db}
}
