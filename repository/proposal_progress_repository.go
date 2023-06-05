package repository

import (
	"fmt"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"

	"be-b-impact.com/csr/utils/common"
	"gorm.io/gorm"
)

type ProposalProgressRepository interface {
	BaseRepository[model.ProposalProgress]
	BaseRepositoryCount[model.ProposalProgress]
	BaseRepositoryPaging[model.ProposalProgress]
}
type proposalProgressRepository struct {
	db *gorm.DB
}

func (pp *proposalProgressRepository) Delete(id string) error {
	return pp.db.Delete(&model.ProposalProgress{}, "id=?", id).Error
}

func (pp *proposalProgressRepository) Get(id string) (*model.ProposalProgress, error) {
	var proposalProgress model.ProposalProgress
	result := pp.db.First(&proposalProgress, "id=?", id).Error
	if result != nil {
		return nil, result
	}
	return &proposalProgress, nil
}

func (pp *proposalProgressRepository) List() ([]model.ProposalProgress, error) {
	var proposalProgress []model.ProposalProgress
	result := pp.db.Find(&proposalProgress).Error
	if result != nil {
		return nil, result
	}
	return proposalProgress, nil
}

func (pp *proposalProgressRepository) Save(payload *model.ProposalProgress) error {
	return pp.db.Save(payload).Error
}

func (pp *proposalProgressRepository) Update(payload *model.ProposalProgress) error {
	updateFields := make(map[string]interface{})

	return pp.db.Model(&model.ProposalProgress{}).Where("id = ?", payload.ID).Updates(updateFields).Error
}

func (pp *proposalProgressRepository) Search(by map[string]interface{}) ([]model.ProposalProgress, error) {
	var proposalProgress []model.ProposalProgress
	query := pp.db
	for key, value := range by {
		// Perform case-insensitive search using ilike
		query = query.Where(fmt.Sprintf("%s ilike ?", key), fmt.Sprintf("%%%v%%", value))
	}
	result := query.Find(&proposalProgress).Error
	if result != nil {
		return nil, result
	}
	return proposalProgress, nil
}

func (pp *proposalProgressRepository) CountData(fieldname string, id string) error {
	var count int64
	var result *gorm.DB

	if id != "" {
		result = pp.db.Model(&model.ProposalProgress{}).Where("name ilike ? AND id <> ?", fieldname, id).Count(&count)
	} else {
		result = pp.db.Model(&model.ProposalProgress{}).Where("name ilike ?", fieldname).Count(&count)
	}
	if result.Error != nil {
		return result.Error
	}
	if count > 0 {
		return fmt.Errorf("name %s already exist", fieldname)
	}
	return nil
}

func (pp *proposalProgressRepository) Paging(requestQueryParam dto.RequestQueryParams) ([]model.ProposalProgress, dto.Paging, error) {
	paginationQuery, orderQuery := pagingValidate(requestQueryParam)

	var proposalProgress []model.ProposalProgress
	query := pp.db
	for key, value := range requestQueryParam.Filter {
		// Perform case-insensitive search using ilike
		query = query.Where(fmt.Sprintf("%s ilike ?", key), fmt.Sprintf("%%%v%%", value))
	}
	err := query.Order(orderQuery).Limit(paginationQuery.Take).Offset(paginationQuery.Skip).Find(&proposalProgress).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}
	var totalRows int64
	err = pp.db.Model(model.ProposalProgress{}).Count(&totalRows).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}
	return proposalProgress, common.Paginate(paginationQuery.Page, paginationQuery.Take, int(totalRows)), nil
}

func NewProposalProgressRepository(db *gorm.DB) ProposalProgressRepository {
	return &proposalProgressRepository{db: db}
}
