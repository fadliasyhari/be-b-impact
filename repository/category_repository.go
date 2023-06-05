package repository

import (
	"fmt"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"

	"be-b-impact.com/csr/utils/common"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	BaseRepository[model.Category]
	BaseRepositoryCount[model.Category]
	BaseRepositoryPaging[model.Category]
}
type categoryRepository struct {
	db *gorm.DB
}

func (ca *categoryRepository) Delete(id string) error {
	return ca.db.Delete(&model.Category{}, "id=?", id).Error
}

func (ca *categoryRepository) Get(id string) (*model.Category, error) {
	var category model.Category
	result := ca.db.First(&category, "id=?", id).Error
	if result != nil {
		return nil, result
	}
	return &category, nil
}

func (ca *categoryRepository) List() ([]model.Category, error) {
	var category []model.Category
	result := ca.db.Find(&category).Error
	if result != nil {
		return nil, result
	}
	return category, nil
}

func (ca *categoryRepository) Save(payload *model.Category) error {
	return ca.db.Save(payload).Error
}

func (ca *categoryRepository) Update(payload *model.Category) error {
	updateFields := make(map[string]interface{})

	// Add fields to be updated based on the payload
	if payload.Parent != "" {
		updateFields["parent"] = payload.Parent
	}

	if payload.Name != "" {
		updateFields["name"] = payload.Name
	}

	if payload.Status != "" {
		updateFields["status"] = payload.Status
	}

	return ca.db.Model(&model.Category{}).Where("id = ?", payload.ID).Updates(updateFields).Error
}

func (ca *categoryRepository) Search(by map[string]interface{}) ([]model.Category, error) {
	var category []model.Category
	query := ca.db
	for key, value := range by {
		// Perform case-insensitive search using ilike
		query = query.Where(fmt.Sprintf("%s ilike ?", key), fmt.Sprintf("%%%v%%", value))
	}
	result := query.Find(&category).Error
	if result != nil {
		return nil, result
	}
	return category, nil
}

func (ca *categoryRepository) CountData(fieldname string, id string) error {
	var count int64
	var result *gorm.DB

	if id != "" {
		result = ca.db.Model(&model.Category{}).Where("name ilike ? AND id <> ?", fieldname, id).Count(&count)
	} else {
		result = ca.db.Model(&model.Category{}).Where("name ilike ?", fieldname).Count(&count)
	}
	if result.Error != nil {
		return result.Error
	}
	if count > 0 {
		return fmt.Errorf("name %s already exist", fieldname)
	}
	return nil
}

func (ca *categoryRepository) Paging(requestQueryParam dto.RequestQueryParams) ([]model.Category, dto.Paging, error) {
	paginationQuery, orderQuery := pagingValidate(requestQueryParam)

	var category []model.Category
	query := ca.db
	for key, value := range requestQueryParam.Filter {
		// Perform case-insensitive search using ilike
		query = query.Where(fmt.Sprintf("%s ilike ?", key), fmt.Sprintf("%%%v%%", value))
	}
	err := query.Order(orderQuery).Limit(paginationQuery.Take).Offset(paginationQuery.Skip).Find(&category).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}
	var totalRows int64
	err = ca.db.Model(model.Category{}).Count(&totalRows).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}
	return category, common.Paginate(paginationQuery.Page, paginationQuery.Take, int(totalRows)), nil
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}
