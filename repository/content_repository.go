package repository

import (
	"fmt"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"

	"be-b-impact.com/csr/utils/common"
	"gorm.io/gorm"
)

type ContentRepository interface {
	BaseRepository[model.Content]
	BaseRepositoryCount[model.Content]
	BaseRepositoryPaging[model.Content]
	BeginTransaction() *gorm.DB
}
type contentRepository struct {
	db *gorm.DB
}

func (co *contentRepository) Delete(id string) error {
	return co.db.Delete(&model.Content{}, "id=?", id).Error
}

func (co *contentRepository) Get(id string) (*model.Content, error) {
	var content model.Content
	result := co.db.Preload("Image").Preload("Category").Preload("TagsContent").Preload("TagsContent.Tag").First(&content, "id=?", id).Error
	if result != nil {
		return nil, result
	}
	return &content, nil
}

func (co *contentRepository) List() ([]model.Content, error) {
	var content []model.Content
	result := co.db.Find(&content).Error
	if result != nil {
		return nil, result
	}
	return content, nil
}

func (co *contentRepository) BeginTransaction() *gorm.DB {
	return co.db.Begin()
}

func (co *contentRepository) Save(payload *model.Content) error {
	return co.db.Save(payload).Error
}

func (co *contentRepository) Update(payload *model.Content) error {
	updateFields := make(map[string]interface{})

	// Add fields to be updated based on the payload
	if payload.Title != "" {
		updateFields["title"] = payload.Title
	}

	if payload.Body != "" {
		updateFields["body"] = payload.Body
	}

	if payload.Status != "" {
		updateFields["status"] = payload.Status
	}

	if payload.CategoryID != "" {
		updateFields["category_id"] = payload.CategoryID
	}

	if payload.DeletedBy != "" {
		updateFields["deleted_by"] = payload.DeletedBy
	}

	if payload.DeletedAt.Time.String() != "" {
		updateFields["deleted_at"] = payload.DeletedAt
	}

	if payload.ProposalID != "" {
		updateFields["proposal_id"] = payload.ProposalID
	}

	if payload.Author != "" {
		updateFields["author"] = payload.Author
	}

	if payload.Excerpt != "" {
		updateFields["excerpt"] = payload.Excerpt
	}

	return co.db.Model(&model.Content{}).Where("id = ?", payload.ID).Updates(updateFields).Error
}

func (co *contentRepository) Search(by map[string]interface{}) ([]model.Content, error) {
	var content []model.Content
	query := co.db
	for key, value := range by {
		// Perform case-insensitive search using ilike
		query = query.Where(fmt.Sprintf("%s ilike ?", key), fmt.Sprintf("%%%v%%", value))
	}
	result := query.Find(&content).Error
	if result != nil {
		return nil, result
	}
	return content, nil
}

func (co *contentRepository) CountData(fieldname string, id string) error {
	var count int64
	var result *gorm.DB

	if id != "" {
		result = co.db.Model(&model.Content{}).Where("name ilike ? AND id <> ?", fieldname, id).Count(&count)
	} else {
		result = co.db.Model(&model.Content{}).Where("name ilike ?", fieldname).Count(&count)
	}
	if result.Error != nil {
		return result.Error
	}
	if count > 0 {
		return fmt.Errorf("name %s already exist", fieldname)
	}
	return nil
}

func (co *contentRepository) Paging(requestQueryParam dto.RequestQueryParams) ([]model.Content, dto.Paging, error) {
	paginationQuery, orderQuery := pagingValidate(requestQueryParam)

	var content []model.Content

	query := co.db.Preload("Image").Preload("Tag").Preload("Category")
	for key, value := range requestQueryParam.Filter {
		// Perform case-insensitive search using ilike
		if key == "category" {
			query = query.Joins("JOIN categories ON categories.id = contents.category_id").
				Where("categories.name = ?", value)
		} else {
			query = query.Where(fmt.Sprintf("%s ilike ?", key), fmt.Sprintf("%%%v%%", value))
		}
	}
	err := query.Order(orderQuery).Limit(paginationQuery.Take).Offset(paginationQuery.Skip).Find(&content).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}
	var totalRows int64
	err = co.db.Model(model.Content{}).Count(&totalRows).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}
	return content, common.Paginate(paginationQuery.Page, paginationQuery.Take, int(totalRows)), nil
}

func NewContentRepository(db *gorm.DB) ContentRepository {
	return &contentRepository{db: db}
}
