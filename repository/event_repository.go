package repository

import (
	"fmt"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"

	"be-b-impact.com/csr/utils/common"
	"gorm.io/gorm"
)

type EventRepository interface {
	BaseRepository[model.Event]
	BaseRepositoryCount[model.Event]
	BaseRepositoryPaging[model.Event]
	PagingDto(requestQueryParam dto.RequestQueryParams, eventIdList []string) ([]dto.EventDTO, dto.Paging, error)
	GetDto(id string) (*dto.EventDTO, error)
	BeginTransaction() *gorm.DB
}
type eventRepository struct {
	db *gorm.DB
	tx *gorm.DB
}

// Get implements EventRepository.
func (ev *eventRepository) Get(id string) (*model.Event, error) {
	var event model.Event
	result := ev.db.Preload("EventImage").Preload("Category").First(&event, "id=?", id).Error
	if result != nil {
		return nil, result
	}
	return &event, nil
}

// Paging implements EventRepository.
func (ev *eventRepository) Paging(requestQueryParam dto.RequestQueryParams) ([]model.Event, dto.Paging, error) {
	paginationQuery, orderQuery := pagingValidate(requestQueryParam)

	var event []model.Event

	query := ev.db.Preload("EventImage").Preload("Category")

	for key, value := range requestQueryParam.Filter {
		// Perform case-insensitive search using ilike
		if key == "category" {
			query = query.Joins("JOIN categories ON categories.id = events.category_id").
				Where("categories.name = ?", value)
		} else {
			query = query.Where(fmt.Sprintf("%s ilike ?", key), fmt.Sprintf("%%%v%%", value))
		}
	}

	var totalRows int64
	err := query.Model(model.Event{}).Count(&totalRows).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}
	err = query.Order(orderQuery).Limit(paginationQuery.Take).Offset(paginationQuery.Skip).Find(&event).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}

	return event, common.Paginate(paginationQuery.Page, paginationQuery.Take, int(totalRows)), nil
}

func (ev *eventRepository) Delete(id string) error {
	return ev.db.Delete(&model.Event{}, "id=?", id).Error
}

func (ev *eventRepository) GetDto(id string) (*dto.EventDTO, error) {
	var event dto.EventDTO
	result := ev.db.Preload("EventImage").Preload("Category").First(&event, "id=?", id).Error
	if result != nil {
		return nil, result
	}
	return &event, nil
}

func (ev *eventRepository) List() ([]model.Event, error) {
	var event []model.Event
	result := ev.db.Find(&event).Error
	if result != nil {
		return nil, result
	}
	return event, nil
}

func (ev *eventRepository) BeginTransaction() *gorm.DB {
	ev.tx = ev.db.Begin()
	return ev.tx
}

func (ev *eventRepository) Save(payload *model.Event) error {
	return ev.tx.Save(payload).Error
}

func (ev *eventRepository) Update(payload *model.Event) error {
	updateFields := make(map[string]interface{})

	// Add fields to be updated based on the payload
	if payload.Title != "" {
		updateFields["title"] = payload.Title
	}

	if payload.Description != "" {
		updateFields["description"] = payload.Description
	}

	if payload.Status != "" {
		updateFields["status"] = payload.Status
	}

	if payload.CategoryID != nil {
		updateFields["category_id"] = payload.CategoryID
	}

	if payload.DeletedBy != "" {
		updateFields["deleted_by"] = payload.DeletedBy
	}

	if payload.DeletedAt.Time.String() != "" {
		updateFields["deleted_at"] = payload.DeletedAt
	}

	if payload.StartDate != "" {
		updateFields["start_date"] = payload.StartDate
	}

	if payload.EndDate != "" {
		updateFields["end_date"] = payload.EndDate
	}

	if payload.Location != "" {
		updateFields["location"] = payload.Location
	}

	return ev.tx.Model(&model.Event{}).Where("id = ?", payload.ID).Updates(updateFields).Error
}

func (ev *eventRepository) Search(by map[string]interface{}) ([]model.Event, error) {
	var event []model.Event
	query := ev.db
	for key, value := range by {
		// Perform case-insensitive search using ilike
		query = query.Where(fmt.Sprintf("%s ilike ?", key), fmt.Sprintf("%%%v%%", value))
	}
	result := query.Find(&event).Error
	if result != nil {
		return nil, result
	}
	return event, nil
}

func (ev *eventRepository) CountData(fieldname string, id string) error {
	var count int64
	var result *gorm.DB

	if id != "" {
		result = ev.db.Model(&model.Event{}).Where("name ilike ? AND id <> ?", fieldname, id).Count(&count)
	} else {
		result = ev.db.Model(&model.Event{}).Where("name ilike ?", fieldname).Count(&count)
	}
	if result.Error != nil {
		return result.Error
	}
	if count > 0 {
		return fmt.Errorf("name %s already exist", fieldname)
	}
	return nil
}

func (ev *eventRepository) PagingDto(requestQueryParam dto.RequestQueryParams, eventIdList []string) ([]dto.EventDTO, dto.Paging, error) {
	paginationQuery, orderQuery := pagingValidate(requestQueryParam)

	var event []dto.EventDTO

	query := ev.db.Preload("EventImage").Preload("Category")

	for key, value := range requestQueryParam.Filter {
		// Perform case-insensitive search using ilike
		if key == "category" {
			query = query.Joins("JOIN categories ON categories.id = events.category_id").
				Where("categories.name = ?, events.id ASC", value)
		} else if key == "user_id" {
			query = query.Where("events.id IN (?)", eventIdList)
		} else {
			query = query.Where(fmt.Sprintf("%s ilike ?", key), fmt.Sprintf("%%%v%%", value))
		}
	}

	subQuery := ev.db.Model(&model.EventParticipant{}).Select(`COUNT("event_participants"."event_id")`).Where(`"event_participants"."event_id" = "events"."id"`)

	var totalRows int64
	err := query.Model(model.Event{}).Count(&totalRows).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}
	if requestQueryParam.Order == "category" {
		if requestQueryParam.Sort == "DESC" {
			query = query.Joins("JOIN categories ON categories.id = events.category_id").
				Order("categories.name DESC, events.id ASC")
		} else {
			query = query.Joins("JOIN categories ON categories.id = events.category_id").
				Order("categories.name ASC, events.id ASC")
		}
		err = query.Select("*,events.id as id, (?) as total_participant", subQuery).Limit(paginationQuery.Take).Offset(paginationQuery.Skip).Find(&event).Error
		if err != nil {
			return nil, dto.Paging{}, err
		}
	} else {
		err = query.Select("*,events.id as id, (?) as total_participant", subQuery).Order(orderQuery).Limit(paginationQuery.Take).Offset(paginationQuery.Skip).Find(&event).Error
		if err != nil {
			return nil, dto.Paging{}, err
		}
	}

	return event, common.Paginate(paginationQuery.Page, paginationQuery.Take, int(totalRows)), nil
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}
