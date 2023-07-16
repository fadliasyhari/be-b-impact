package usecase

import (
	"fmt"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"
	"be-b-impact.com/csr/repository"
	"gorm.io/gorm"
)

type ProposalProgressUseCase interface {
	BaseUseCase[model.ProposalProgress]
	BaseUseCasePaging[model.ProposalProgress]
	SaveTrx(payload *model.ProposalProgress, tx *gorm.DB) error
}

type proposalProgressUseCase struct {
	repo           repository.ProposalProgressRepository
	notificationUC NotificationUseCase
}

func (pp *proposalProgressUseCase) DeleteData(id string) error {
	proposalProgress, err := pp.FindById(id)
	if err != nil {
		return fmt.Errorf("proposalProgress with ID %s not found", id)
	}
	return pp.repo.Delete(proposalProgress.ID)
}

func (pp *proposalProgressUseCase) FindAll() ([]model.ProposalProgress, error) {
	return pp.repo.List()
}

func (pp *proposalProgressUseCase) FindById(id string) (*model.ProposalProgress, error) {
	proposalProgress, err := pp.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("proposalProgress with ID %s not found", id)
	}
	return proposalProgress, nil
}

func (pp *proposalProgressUseCase) SaveData(payload *model.ProposalProgress) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }

	tx := pp.repo.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := pp.repo.SaveTrx(payload, tx); err != nil {
		tx.Rollback()
		return err
	}

	progress, err := pp.repo.GetProposalById(payload.ProposalID)
	if err != nil {
		return fmt.Errorf("proposalProgress with ID %s not found", payload.ID)
	}

	notifPayload := model.Notification{
		Body:   fmt.Sprintf("There's new update on your proposal <b>%s</b>", progress.ProposalDetail.ProjectName),
		UserID: progress.CreatedBy,
		Type:   "proposal",
		TypeID: progress.ID,
	}

	err = pp.notificationUC.SaveNotifDetail(&notifPayload, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil

}

func (pp *proposalProgressUseCase) SaveTrx(payload *model.ProposalProgress, tx *gorm.DB) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }
	return pp.repo.SaveTrx(payload, tx)
}

func (pp *proposalProgressUseCase) UpdateData(payload *model.ProposalProgress) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }
	// cek jika data sudah ada -> count > 0

	if payload.ID != "" {
		_, err := pp.FindById(payload.ID)
		if err != nil {
			return fmt.Errorf("proposalProgress with ID %s not found", payload.ID)
		}
	}
	return pp.repo.Update(payload)
}

func (pp *proposalProgressUseCase) SearchBy(by map[string]interface{}) ([]model.ProposalProgress, error) {
	proposalProgresss, err := pp.repo.Search(by)
	if err != nil {
		return nil, fmt.Errorf("request invalid")
	}
	return proposalProgresss, nil
}

func (pp *proposalProgressUseCase) Pagination(requestQueryParams dto.RequestQueryParams) ([]model.ProposalProgress, dto.Paging, error) {
	if !requestQueryParams.QueryParams.IsSortValid() {
		return nil, dto.Paging{}, fmt.Errorf("invalid sort by: %s", requestQueryParams.QueryParams.Sort)
	}
	return pp.repo.Paging(requestQueryParams)
}

func NewProposalProgressUseCase(repo repository.ProposalProgressRepository, notificationUC NotificationUseCase) ProposalProgressUseCase {
	return &proposalProgressUseCase{
		repo:           repo,
		notificationUC: notificationUC,
	}
}
