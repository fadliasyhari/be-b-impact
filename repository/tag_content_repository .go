package repository

import (
	"fmt"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"

	"be-b-impact.com/csr/utils/common"
	"gorm.io/gorm"
)

type TagsContentRepository interface {
	BaseRepository[model.TagsContent]
	BaseRepositoryCount[model.TagsContent]
	BaseRepositoryPaging[model.TagsContent]
}
type tagsContentRepository struct {
	db *gorm.DB
}

func (tc *tagsContentRepository) Delete(id string) error {
	return tc.db.Delete(&model.TagsContent{}, "id=?", id).Error
}

func (tc *tagsContentRepository) Get(id string) (*model.TagsContent, error) {
	var tagsContent model.TagsContent
	result := tc.db.First(&tagsContent, "id=?", id).Error
	if result != nil {
		return nil, result
	}
	return &tagsContent, nil
}

func (tc *tagsContentRepository) List() ([]model.TagsContent, error) {
	var tagsContent []model.TagsContent
	result := tc.db.Find(&tagsContent).Error
	if result != nil {
		return nil, result
	}
	return tagsContent, nil
}

func (tc *tagsContentRepository) Save(payload *model.TagsContent) error {
	return tc.db.Save(payload).Error
}

func (tc *tagsContentRepository) Update(payload *model.TagsContent) error {
	return tc.db.Model(&model.TagsContent{}).Where("id = ?", payload.ContentID).Updates(payload).Error
}

func (tc *tagsContentRepository) Search(by map[string]interface{}) ([]model.TagsContent, error) {
	var tagsContent []model.TagsContent
	result := tc.db.Where(by).Find(&tagsContent).Error
	if result != nil {
		return nil, result
	}
	return tagsContent, nil
}

func (tc *tagsContentRepository) CountData(fieldname string, id string) error {
	var count int64
	var result *gorm.DB

	if id != "" {
		result = tc.db.Model(&model.TagsContent{}).Where("name ilike ? AND id <> ?", fieldname, id).Count(&count)
	} else {
		result = tc.db.Model(&model.TagsContent{}).Where("name ilike ?", fieldname).Count(&count)
	}
	if result.Error != nil {
		return result.Error
	}
	if count > 0 {
		return fmt.Errorf("name %s already exist", fieldname)
	}
	return nil
}

func (tc *tagsContentRepository) Paging(requestQueryParam dto.RequestQueryParams) ([]model.TagsContent, dto.Paging, error) {
	paginationQuery, orderQuery := pagingValidate(requestQueryParam)

	var tagsContent []model.TagsContent
	err := tc.db.Order(orderQuery).Limit(paginationQuery.Take).Offset(paginationQuery.Skip).Find(&tagsContent).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}
	var totalRows int64
	err = tc.db.Model(model.TagsContent{}).Count(&totalRows).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}
	return tagsContent, common.Paginate(paginationQuery.Page, paginationQuery.Take, int(totalRows)), nil
}

func NewTagsContentRepository(db *gorm.DB) TagsContentRepository {
	return &tagsContentRepository{db: db}
}
