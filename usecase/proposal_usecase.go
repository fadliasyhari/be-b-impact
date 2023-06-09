package usecase

import (
	"fmt"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"
	"be-b-impact.com/csr/repository"
)

type ProposalUseCase interface {
	BaseUseCase[model.Proposal]
	BaseUseCasePaging[model.Proposal]
	FindPropById(id string) (*dto.Proposal, error)
}

type proposalUseCase struct {
	repo repository.ProposalRepository
}

func (pr *proposalUseCase) DeleteData(id string) error {
	proposal, err := pr.FindById(id)
	if err != nil {
		return fmt.Errorf("proposal with ID %s not found", id)
	}
	return pr.repo.Delete(proposal.ID)
}

func (pr *proposalUseCase) FindAll() ([]model.Proposal, error) {
	return pr.repo.List()
}

func (pr *proposalUseCase) FindPropById(id string) (*dto.Proposal, error) {
	proposal, err := pr.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("proposal with ID %s not found", id)
	}
	return proposal, nil
}

func (pr *proposalUseCase) FindById(id string) (*model.Proposal, error) {
	proposal, err := pr.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("proposal with ID %s not found", id)
	}
	return proposal, nil
}

func (pr *proposalUseCase) SaveData(payload *model.Proposal) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }
	// cek jika data sudah ada -> count > 0

	if payload.ID != "" {
		_, err := pr.FindById(payload.ID)
		if err != nil {
			return fmt.Errorf("proposal with ID %s not found", payload.ID)
		}
	}
	return pr.repo.Save(payload)
}

func (pr *proposalUseCase) UpdateData(payload *model.Proposal) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }
	// cek jika data sudah ada -> count > 0

	if payload.ID != "" {
		_, err := pr.FindById(payload.ID)
		if err != nil {
			return fmt.Errorf("proposal with ID %s not found", payload.ID)
		}
	}
	return pr.repo.Update(payload)
}

func (pr *proposalUseCase) SearchBy(by map[string]interface{}) ([]model.Proposal, error) {
	proposals, err := pr.repo.Search(by)
	if err != nil {
		return nil, fmt.Errorf("request invalid")
	}
	return proposals, nil
}

func (pr *proposalUseCase) Pagination(requestQueryParams dto.RequestQueryParams) ([]model.Proposal, dto.Paging, error) {
	if !requestQueryParams.QueryParams.IsSortValid() {
		return nil, dto.Paging{}, fmt.Errorf("invalid sort by: %s", requestQueryParams.QueryParams.Sort)
	}
	return pr.repo.Paging(requestQueryParams)
}

func NewProposalUseCase(repo repository.ProposalRepository) ProposalUseCase {
	return &proposalUseCase{repo: repo}
}
