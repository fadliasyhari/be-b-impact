package repository

import (
	"fmt"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"

	"be-b-impact.com/csr/utils/common"
	"gorm.io/gorm"
)

type ProgressRepository interface {
	BaseRepository[model.Progress]
	BaseRepositoryCount[model.Progress]
	BaseRepositoryPaging[model.Progress]
}
type progressRepository struct {
	db *gorm.DB
}

func (pg *progressRepository) Delete(id string) error {
	return pg.db.Delete(&model.Progress{}, "id=?", id).Error
}

func (pg *progressRepository) Get(id string) (*model.Progress, error) {
	var progress model.Progress
	result := pg.db.First(&progress, "id=?", id).Error
	if result != nil {
		return nil, result
	}
	return &progress, nil
}

func (pg *progressRepository) List() ([]model.Progress, error) {
	var progress []model.Progress
	result := pg.db.Find(&progress).Error
	if result != nil {
		return nil, result
	}
	return progress, nil
}

func (pg *progressRepository) Save(payload *model.Progress) error {
	return pg.db.Save(payload).Error
}

func (pg *progressRepository) Update(payload *model.Progress) error {
	updateFields := make(map[string]interface{})

	// Add fields to be updated based on the payload
	if payload.Name != "" {
		updateFields["name"] = payload.Name
	}

	return pg.db.Model(&model.Progress{}).Where("id = ?", payload.ID).Updates(updateFields).Error
}

func (pg *progressRepository) Search(by map[string]interface{}) ([]model.Progress, error) {
	var progress []model.Progress
	query := pg.db
	for key, value := range by {
		// Perform case-insensitive search using ilike
		query = query.Where(fmt.Sprintf("%s ilike ?", key), fmt.Sprintf("%%%v%%", value))
	}
	result := query.Find(&progress).Error
	if result != nil {
		return nil, result
	}
	return progress, nil
}

func (pg *progressRepository) CountData(fieldname string, id string) error {
	var count int64
	var result *gorm.DB

	if id != "" {
		result = pg.db.Model(&model.Progress{}).Where("name ilike ? AND id <> ?", fieldname, id).Count(&count)
	} else {
		result = pg.db.Model(&model.Progress{}).Where("name ilike ?", fieldname).Count(&count)
	}
	if result.Error != nil {
		return result.Error
	}
	if count > 0 {
		return fmt.Errorf("name %s already exist", fieldname)
	}
	return nil
}

func (pg *progressRepository) Paging(requestQueryParam dto.RequestQueryParams) ([]model.Progress, dto.Paging, error) {
	paginationQuery, orderQuery := pagingValidate(requestQueryParam)

	var progress []model.Progress
	query := pg.db
	for key, value := range requestQueryParam.Filter {
		// Perform case-insensitive search using ilike
		query = query.Where(fmt.Sprintf("%s ilike ?", key), fmt.Sprintf("%%%v%%", value))
	}
	err := query.Order(orderQuery).Limit(paginationQuery.Take).Offset(paginationQuery.Skip).Find(&progress).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}
	var totalRows int64
	err = pg.db.Model(model.Progress{}).Count(&totalRows).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}
	return progress, common.Paginate(paginationQuery.Page, paginationQuery.Take, int(totalRows)), nil
}

func NewProgressRepository(db *gorm.DB) ProgressRepository {
	return &progressRepository{db: db}
}
