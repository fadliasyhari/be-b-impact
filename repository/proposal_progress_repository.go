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
	SaveTrx(payload *model.ProposalProgress, tx *gorm.DB) error
	UpdateTrx(payload *model.ProposalProgress, tx *gorm.DB) error
	GetProposalById(id string) (*model.Proposal, error)
	BeginTransaction() *gorm.DB
}
type proposalProgressRepository struct {
	db *gorm.DB
	tx *gorm.DB
}

func (pp *proposalProgressRepository) Delete(id string) error {
	return pp.db.Delete(&model.ProposalProgress{}, "id=?", id).Error
}

func (pp *proposalProgressRepository) Get(id string) (*model.ProposalProgress, error) {
	var proposalProgress model.ProposalProgress
	result := pp.db.Preload("Proposal").Preload("Proposal.ProposalDetail").First(&proposalProgress, "id=?", id).Error
	if result != nil {
		return nil, result
	}
	return &proposalProgress, nil
}

func (pp *proposalProgressRepository) GetProposalById(id string) (*model.Proposal, error) {
	var proposal model.Proposal
	result := pp.db.Preload("ProposalDetail").Preload("ProposalDetail.PartnershipType").Preload("File").Preload("OrganizatonType").Preload("ProposalProgress").Preload("ProposalProgress.Progress").First(&proposal, "id=?", id).Error
	if result != nil {
		return nil, result
	}
	return &proposal, nil
}

func (pp *proposalProgressRepository) List() ([]model.ProposalProgress, error) {
	var proposalProgress []model.ProposalProgress
	result := pp.db.Find(&proposalProgress).Error
	if result != nil {
		return nil, result
	}
	return proposalProgress, nil
}

func (pp *proposalProgressRepository) BeginTransaction() *gorm.DB {
	pp.tx = pp.db.Begin()
	return pp.tx
}

func (pp *proposalProgressRepository) Save(payload *model.ProposalProgress) error {
	return pp.db.Save(payload).Error
}

func (pp *proposalProgressRepository) SaveTrx(payload *model.ProposalProgress, tx *gorm.DB) error {
	return tx.Create(payload).Error
}

func (pp *proposalProgressRepository) Update(payload *model.ProposalProgress) error {
	updateFields := make(map[string]interface{})

	if payload.Note != "" {
		updateFields["note"] = payload.Note
	}

	if payload.Status != "" {
		updateFields["status"] = payload.Status
	}

	if payload.ReviewLocation != "" {
		updateFields["review_location"] = payload.ReviewLocation
	}

	if !payload.ReviewDate.IsZero() {
		updateFields["review_date"] = payload.ReviewDate
	}

	if payload.ReviewCP != "" {
		updateFields["review_cp"] = payload.ReviewCP
	}

	if payload.ReviewFeedback != "" {
		updateFields["review_feedback"] = payload.ReviewFeedback
	}

	return pp.db.Model(&model.ProposalProgress{}).Where("id = ?", payload.ID).Updates(updateFields).Error
}

func (pp *proposalProgressRepository) UpdateTrx(payload *model.ProposalProgress, tx *gorm.DB) error {
	updateFields := make(map[string]interface{})

	if payload.Note != "" {
		updateFields["note"] = payload.Note
	}

	if payload.Status != "" {
		updateFields["status"] = payload.Status
	}

	if payload.ReviewLocation != "" {
		updateFields["review_location"] = payload.ReviewLocation
	}

	if !payload.ReviewDate.IsZero() {
		updateFields["review_date"] = payload.ReviewDate
	}

	if payload.ReviewCP != "" {
		updateFields["review_cp"] = payload.ReviewCP
	}

	if payload.ReviewFeedback != "" {
		updateFields["review_feedback"] = payload.ReviewFeedback
	}

	return tx.Model(&model.ProposalProgress{}).Where("id = ?", payload.ID).Updates(updateFields).Error
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
