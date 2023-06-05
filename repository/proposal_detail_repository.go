package repository

import (
	"fmt"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"

	"be-b-impact.com/csr/utils/common"
	"gorm.io/gorm"
)

type ProposalDetailRepository interface {
	BaseRepository[model.ProposalDetail]
	BaseRepositoryCount[model.ProposalDetail]
	BaseRepositoryPaging[model.ProposalDetail]
}
type proposalDetailRepository struct {
	db *gorm.DB
}

func (pd *proposalDetailRepository) Delete(id string) error {
	return pd.db.Delete(&model.ProposalDetail{}, "id=?", id).Error
}

func (pd *proposalDetailRepository) Get(id string) (*model.ProposalDetail, error) {
	var proposalDetail model.ProposalDetail
	result := pd.db.First(&proposalDetail, "id=?", id).Error
	if result != nil {
		return nil, result
	}
	return &proposalDetail, nil
}

func (pd *proposalDetailRepository) List() ([]model.ProposalDetail, error) {
	var proposalDetail []model.ProposalDetail
	result := pd.db.Find(&proposalDetail).Error
	if result != nil {
		return nil, result
	}
	return proposalDetail, nil
}

func (pd *proposalDetailRepository) Save(payload *model.ProposalDetail) error {
	return pd.db.Save(payload).Error
}

func (pd *proposalDetailRepository) Update(payload *model.ProposalDetail) error {
	updateFields := make(map[string]interface{})

	// Add fields to be updated based on the payload
	if payload.ProjectName != "" {
		updateFields["project_name"] = payload.ProjectName
	}

	if payload.PartnershipTypeID != "" {
		updateFields["partnership_type_id"] = payload.PartnershipTypeID
	}

	if !payload.StartDate.IsZero() {
		updateFields["start_date"] = payload.StartDate
	}

	if !payload.EndDate.IsZero() {
		updateFields["end_date"] = payload.EndDate
	}

	if payload.Objective != "" {
		updateFields["objective"] = payload.Objective
	}

	if payload.Alignment != "" {
		updateFields["alignment"] = payload.Alignment
	}

	return pd.db.Model(&model.ProposalDetail{}).Where("id = ?", payload.ID).Updates(updateFields).Error
}

func (pd *proposalDetailRepository) Search(by map[string]interface{}) ([]model.ProposalDetail, error) {
	var proposalDetail []model.ProposalDetail
	query := pd.db
	for key, value := range by {
		// Perform case-insensitive search using ilike
		query = query.Where(fmt.Sprintf("%s ilike ?", key), fmt.Sprintf("%%%v%%", value))
	}
	result := query.Find(&proposalDetail).Error
	if result != nil {
		return nil, result
	}
	return proposalDetail, nil
}

func (pd *proposalDetailRepository) CountData(fieldname string, id string) error {
	var count int64
	var result *gorm.DB

	if id != "" {
		result = pd.db.Model(&model.ProposalDetail{}).Where("name ilike ? AND id <> ?", fieldname, id).Count(&count)
	} else {
		result = pd.db.Model(&model.ProposalDetail{}).Where("name ilike ?", fieldname).Count(&count)
	}
	if result.Error != nil {
		return result.Error
	}
	if count > 0 {
		return fmt.Errorf("name %s already exist", fieldname)
	}
	return nil
}

func (pd *proposalDetailRepository) Paging(requestQueryParam dto.RequestQueryParams) ([]model.ProposalDetail, dto.Paging, error) {
	paginationQuery, orderQuery := pagingValidate(requestQueryParam)

	var proposalDetail []model.ProposalDetail
	query := pd.db
	for key, value := range requestQueryParam.Filter {
		// Perform case-insensitive search using ilike
		query = query.Where(fmt.Sprintf("%s ilike ?", key), fmt.Sprintf("%%%v%%", value))
	}
	err := query.Order(orderQuery).Limit(paginationQuery.Take).Offset(paginationQuery.Skip).Find(&proposalDetail).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}
	var totalRows int64
	err = pd.db.Model(model.ProposalDetail{}).Count(&totalRows).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}
	return proposalDetail, common.Paginate(paginationQuery.Page, paginationQuery.Take, int(totalRows)), nil
}

func NewProposalDetailRepository(db *gorm.DB) ProposalDetailRepository {
	return &proposalDetailRepository{db: db}
}
