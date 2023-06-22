package repository

import (
	"fmt"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"

	"be-b-impact.com/csr/utils/common"
	"gorm.io/gorm"
)

type ProposalRepository interface {
	BaseRepository[model.Proposal]
	BaseRepositoryCount[model.Proposal]
	BaseRepositoryPaging[model.Proposal]
	GetByID(id string) (*dto.Proposal, error)
	BeginTransaction() *gorm.DB
	UpdateMinor(payload *model.Proposal, tx *gorm.DB) error
	GetTx(id string, tx *gorm.DB) (*model.Proposal, error)
	UpdateBasic(payload *model.Proposal) error
}
type proposalRepository struct {
	tx *gorm.DB
	db *gorm.DB
}

func (pr *proposalRepository) Delete(id string) error {
	return pr.db.Delete(&model.Proposal{}, "id=?", id).Error
}

func (pr *proposalRepository) GetByID(id string) (*dto.Proposal, error) {
	var proposal dto.Proposal
	result := pr.db.Preload("ProposalDetail").Preload("ProposalDetail.PartnershipType").Preload("File").Preload("OrganizatonType").Preload("ProposalProgress").Preload("ProposalProgress.Progress").First(&proposal, "id=?", id).Error
	if result != nil {
		return nil, result
	}
	return &proposal, nil
}

func (pr *proposalRepository) Get(id string) (*model.Proposal, error) {
	var proposal model.Proposal
	result := pr.db.Preload("ProposalDetail").Preload("ProposalDetail.PartnershipType").Preload("File").Preload("OrganizatonType").Preload("ProposalProgress").Preload("ProposalProgress.Progress").First(&proposal, "id=?", id).Error
	if result != nil {
		return nil, result
	}
	return &proposal, nil
}

func (pr *proposalRepository) GetTx(id string, tx *gorm.DB) (*model.Proposal, error) {
	var proposal model.Proposal
	result := tx.Preload("ProposalDetail").Preload("ProposalDetail.PartnershipType").Preload("File").Preload("OrganizatonType").Preload("ProposalProgress").Preload("ProposalProgress.Progress").First(&proposal, "id=?", id).Error
	if result != nil {
		return nil, result
	}
	return &proposal, nil
}

func (pr *proposalRepository) List() ([]model.Proposal, error) {
	var proposal []model.Proposal
	result := pr.db.Find(&proposal).Error
	if result != nil {
		return nil, result
	}
	return proposal, nil
}

func (pr *proposalRepository) BeginTransaction() *gorm.DB {
	pr.tx = pr.db.Begin()
	return pr.tx
}

func (pr *proposalRepository) Save(payload *model.Proposal) error {
	return pr.tx.Create(payload).Error
}

func (pr *proposalRepository) Update(payload *model.Proposal) error {
	updateFields := make(map[string]interface{})

	// Add fields to be updated based on the payload
	if payload.OrgName != "" {
		updateFields["org_name"] = payload.OrgName
	}

	if payload.OrganizationTypeID != nil {
		updateFields["organization_type_id"] = payload.OrganizationTypeID
	}

	if payload.Email != "" {
		updateFields["email"] = payload.Email
	}

	if payload.Phone != "" {
		updateFields["phone"] = payload.Phone
	}

	if payload.PICName != "" {
		updateFields["pic_name"] = payload.PICName
	}

	if payload.City != "" {
		updateFields["city"] = payload.City
	}

	if payload.PostalCode != "" {
		updateFields["postal_code"] = payload.PostalCode
	}

	if payload.Address != "" {
		updateFields["address"] = payload.Address
	}

	if payload.Description != "" {
		updateFields["description"] = payload.Description
	}

	if payload.Status != "" {
		updateFields["status"] = payload.Status
	}

	if payload.DeletedBy != "" {
		updateFields["deleted_by"] = payload.DeletedBy
	}

	if payload.ReviewerID != "" {
		updateFields["reviewer_id"] = payload.ReviewerID
	}

	if payload.CurrentProgress != "" {
		updateFields["current_progress"] = payload.CurrentProgress
	}

	return pr.tx.Model(&model.Proposal{}).Where("id = ?", payload.ID).Updates(updateFields).Error
}

func (pr *proposalRepository) UpdateBasic(payload *model.Proposal) error {
	updateFields := make(map[string]interface{})

	// Add fields to be updated based on the payload
	if payload.OrgName != "" {
		updateFields["org_name"] = payload.OrgName
	}

	if payload.OrganizationTypeID != nil {
		updateFields["organization_type_id"] = payload.OrganizationTypeID
	}

	if payload.Email != "" {
		updateFields["email"] = payload.Email
	}

	if payload.Phone != "" {
		updateFields["phone"] = payload.Phone
	}

	if payload.PICName != "" {
		updateFields["pic_name"] = payload.PICName
	}

	if payload.City != "" {
		updateFields["city"] = payload.City
	}

	if payload.PostalCode != "" {
		updateFields["postal_code"] = payload.PostalCode
	}

	if payload.Address != "" {
		updateFields["address"] = payload.Address
	}

	if payload.Description != "" {
		updateFields["description"] = payload.Description
	}

	if payload.Status != "" {
		updateFields["status"] = payload.Status
	}

	if payload.DeletedBy != "" {
		updateFields["deleted_by"] = payload.DeletedBy
	}

	if payload.ReviewerID != "" {
		updateFields["reviewer_id"] = payload.ReviewerID
	}

	if payload.CurrentProgress != "" {
		updateFields["current_progress"] = payload.CurrentProgress
	}

	return pr.db.Model(&model.Proposal{}).Where("id = ?", payload.ID).Updates(updateFields).Error
}

func (pr *proposalRepository) UpdateMinor(payload *model.Proposal, tx *gorm.DB) error {
	updateFields := make(map[string]interface{})

	// Add fields to be updated based on the payload
	if payload.OrgName != "" {
		updateFields["org_name"] = payload.OrgName
	}

	if payload.OrganizationTypeID != nil {
		updateFields["organization_type_id"] = payload.OrganizationTypeID
	}

	if payload.Email != "" {
		updateFields["email"] = payload.Email
	}

	if payload.Phone != "" {
		updateFields["phone"] = payload.Phone
	}

	if payload.PICName != "" {
		updateFields["pic_name"] = payload.PICName
	}

	if payload.City != "" {
		updateFields["city"] = payload.City
	}

	if payload.PostalCode != "" {
		updateFields["postal_code"] = payload.PostalCode
	}

	if payload.Address != "" {
		updateFields["address"] = payload.Address
	}

	if payload.Description != "" {
		updateFields["description"] = payload.Description
	}

	if payload.Status != "" {
		updateFields["status"] = payload.Status
	}

	if payload.DeletedBy != "" {
		updateFields["deleted_by"] = payload.DeletedBy
	}

	if payload.ReviewerID != "" {
		updateFields["reviewer_id"] = payload.ReviewerID
	}

	if payload.CurrentProgress != "" {
		updateFields["current_progress"] = payload.CurrentProgress
	}

	return tx.Model(&model.Proposal{}).Where("id = ?", payload.ID).Updates(updateFields).Error
}

func (pr *proposalRepository) Search(by map[string]interface{}) ([]model.Proposal, error) {
	var proposal []model.Proposal
	query := pr.db
	for key, value := range by {
		// Perform case-insensitive search using ilike
		query = query.Where(fmt.Sprintf("%s ilike ?", key), fmt.Sprintf("%%%v%%", value))
	}
	result := query.Find(&proposal).Error
	if result != nil {
		return nil, result
	}
	return proposal, nil
}

func (pr *proposalRepository) CountData(fieldname string, id string) error {
	var count int64
	var result *gorm.DB

	if id != "" {
		result = pr.db.Model(&model.Proposal{}).Where("name ilike ? AND id <> ?", fieldname, id).Count(&count)
	} else {
		result = pr.db.Model(&model.Proposal{}).Where("name ilike ?", fieldname).Count(&count)
	}
	if result.Error != nil {
		return result.Error
	}
	if count > 0 {
		return fmt.Errorf("name %s already exist", fieldname)
	}
	return nil
}

func (pr *proposalRepository) Paging(requestQueryParam dto.RequestQueryParams) ([]model.Proposal, dto.Paging, error) {
	paginationQuery, orderQuery := pagingValidate(requestQueryParam)

	var proposal []model.Proposal
	query := pr.db.Preload("ProposalDetail").Preload("ProposalProgress").Preload("ProposalProgress.Progress").Preload("File")

	for key, value := range requestQueryParam.Filter {
		// Perform case-insensitive search using ilike
		query = query.Where(fmt.Sprintf("%s ilike ?", key), fmt.Sprintf("%%%v%%", value))
	}

	err := query.Order(orderQuery).Limit(paginationQuery.Take).Offset(paginationQuery.Skip).Find(&proposal).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}
	var totalRows int64
	err = query.Model(model.Proposal{}).Count(&totalRows).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}
	return proposal, common.Paginate(paginationQuery.Page, paginationQuery.Take, int(totalRows)), nil
}

func NewProposalRepository(db *gorm.DB) ProposalRepository {
	return &proposalRepository{db: db}
}
