package repository

import (
	"fmt"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"

	"be-b-impact.com/csr/utils/common"
	"gorm.io/gorm"
)

type TagRepository interface {
	BaseRepository[model.Tag]
	BaseRepositoryCount[model.Tag]
	BaseRepositoryPaging[model.Tag]
}
type tagRepository struct {
	db *gorm.DB
}

func (ta *tagRepository) Delete(id string) error {
	return ta.db.Delete(&model.Tag{}, "id=?", id).Error
}

func (ta *tagRepository) Get(id string) (*model.Tag, error) {
	var tag model.Tag
	result := ta.db.First(&tag, "id=?", id).Error
	if result != nil {
		return nil, result
	}
	return &tag, nil
}

func (ta *tagRepository) List() ([]model.Tag, error) {
	var tag []model.Tag
	result := ta.db.Find(&tag).Error
	if result != nil {
		return nil, result
	}
	return tag, nil
}

func (ta *tagRepository) Save(payload *model.Tag) error {
	return ta.db.Save(payload).Error
}

func (ta *tagRepository) Update(payload *model.Tag) error {
	updateFields := make(map[string]interface{})

	// Add fields to be updated based on the payload
	if payload.Name != "" {
		updateFields["name"] = payload.Name
	}

	return ta.db.Model(&model.Tag{}).Where("id = ?", payload.ID).Updates(updateFields).Error
}

func (ta *tagRepository) Search(by map[string]interface{}) ([]model.Tag, error) {
	var tag []model.Tag
	query := ta.db
	for key, value := range by {
		// Perform case-insensitive search using ilike
		query = query.Where(fmt.Sprintf("%s ilike ?", key), fmt.Sprintf("%%%v%%", value))
	}
	result := query.Find(&tag).Error
	if result != nil {
		return nil, result
	}
	return tag, nil
}

func (ta *tagRepository) CountData(fieldname string, id string) error {
	var count int64
	var result *gorm.DB

	if id != "" {
		result = ta.db.Model(&model.Tag{}).Where("name ilike ? AND id <> ?", fieldname, id).Count(&count)
	} else {
		result = ta.db.Model(&model.Tag{}).Where("name ilike ?", fieldname).Count(&count)
	}
	if result.Error != nil {
		return result.Error
	}
	if count > 0 {
		return fmt.Errorf("name %s already exist", fieldname)
	}
	return nil
}

func (ta *tagRepository) Paging(requestQueryParam dto.RequestQueryParams) ([]model.Tag, dto.Paging, error) {
	paginationQuery, orderQuery := pagingValidate(requestQueryParam)

	var tag []model.Tag
	query := ta.db
	for key, value := range requestQueryParam.Filter {
		// Perform case-insensitive search using ilike
		query = query.Where(fmt.Sprintf("%s ilike ?", key), fmt.Sprintf("%%%v%%", value))
	}
	err := query.Order(orderQuery).Limit(paginationQuery.Take).Offset(paginationQuery.Skip).Find(&tag).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}
	var totalRows int64
	err = ta.db.Model(model.Tag{}).Count(&totalRows).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}
	return tag, common.Paginate(paginationQuery.Page, paginationQuery.Take, int(totalRows)), nil
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}
