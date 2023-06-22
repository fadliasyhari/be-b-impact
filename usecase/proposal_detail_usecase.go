package usecase

import (
	"fmt"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"
	"be-b-impact.com/csr/repository"
	"gorm.io/gorm"
)

type ProposalDetailUseCase interface {
	BaseUseCase[model.ProposalDetail]
	BaseUseCasePaging[model.ProposalDetail]
	SavePropoDetail(payload *model.ProposalDetail, tx *gorm.DB) error
	UpdatePropoDetail(payload *model.ProposalDetail, tx *gorm.DB) error
}

type proposalDetailUseCase struct {
	repo repository.ProposalDetailRepository
}

func (pd *proposalDetailUseCase) DeleteData(id string) error {
	proposalDetail, err := pd.FindById(id)
	if err != nil {
		return fmt.Errorf("proposalDetail with ID %s not found", id)
	}
	return pd.repo.Delete(proposalDetail.ID)
}

func (pd *proposalDetailUseCase) FindAll() ([]model.ProposalDetail, error) {
	return pd.repo.List()
}

func (pd *proposalDetailUseCase) FindById(id string) (*model.ProposalDetail, error) {
	proposalDetail, err := pd.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("proposalDetail with ID %s not found", id)
	}
	return proposalDetail, nil
}

func (pd *proposalDetailUseCase) SaveData(payload *model.ProposalDetail) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }
	// cek jika data sudah ada -> count > 0

	if payload.ID != "" {
		_, err := pd.FindById(payload.ID)
		if err != nil {
			return fmt.Errorf("proposalDetail with ID %s not found", payload.ID)
		}
	}
	return pd.repo.Save(payload)
}

func (pd *proposalDetailUseCase) SavePropoDetail(payload *model.ProposalDetail, tx *gorm.DB) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }

	return pd.repo.SaveTrx(payload, tx)
}

func (pd *proposalDetailUseCase) UpdatePropoDetail(payload *model.ProposalDetail, tx *gorm.DB) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }

	return pd.repo.UpdateTrx(payload, tx)
}

func (pd *proposalDetailUseCase) UpdateData(payload *model.ProposalDetail) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }
	// cek jika data sudah ada -> count > 0

	if payload.ID != "" {
		_, err := pd.FindById(payload.ID)
		if err != nil {
			return fmt.Errorf("proposalDetail with ID %s not found", payload.ID)
		}
	}
	return pd.repo.Update(payload)
}

func (pd *proposalDetailUseCase) SearchBy(by map[string]interface{}) ([]model.ProposalDetail, error) {
	proposalDetails, err := pd.repo.Search(by)
	if err != nil {
		return nil, fmt.Errorf("request invalid")
	}
	return proposalDetails, nil
}

func (pd *proposalDetailUseCase) Pagination(requestQueryParams dto.RequestQueryParams) ([]model.ProposalDetail, dto.Paging, error) {
	if !requestQueryParams.QueryParams.IsSortValid() {
		return nil, dto.Paging{}, fmt.Errorf("invalid sort by: %s", requestQueryParams.QueryParams.Sort)
	}
	return pd.repo.Paging(requestQueryParams)
}

func NewProposalDetailUseCase(repo repository.ProposalDetailRepository) ProposalDetailUseCase {
	return &proposalDetailUseCase{repo: repo}
}
