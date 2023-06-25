package usecase

import (
	"fmt"
	"mime/multipart"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"
	"be-b-impact.com/csr/repository"
	"gorm.io/gorm"
)

type ProposalUseCase interface {
	BaseUseCase[model.Proposal]
	BaseUseCasePaging[model.Proposal]
	FindPropById(id string) (*dto.Proposal, error)
	SavePropo(payload *model.Proposal, payloadDetail *model.ProposalDetail, org_profile multipart.File, propo_doc multipart.File) error
	UpdateMinor(payload *model.Proposal, tx *gorm.DB) error
	UpdatePropo(payload *model.Proposal, payloadDetail *model.ProposalDetail, org_profile multipart.File, propo_doc multipart.File) error
	FindByIdTx(id string, tx *gorm.DB) (*model.Proposal, error)
}

type proposalUseCase struct {
	repo            repository.ProposalRepository
	propoDetailUC   ProposalDetailUseCase
	fileUC          FileUseCase
	progressUC      ProgressUseCase
	propoProgressUC ProposalProgressUseCase
	userUC          UsersUseCase
}

// SaveData implements ProposalUseCase.
func (*proposalUseCase) SaveData(payload *model.Proposal) error {
	panic("unimplemented")
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

func (pr *proposalUseCase) FindByIdTx(id string, tx *gorm.DB) (*model.Proposal, error) {
	proposal, err := pr.repo.GetTx(id, tx)
	if err != nil {
		return nil, fmt.Errorf("proposal with ID %s not found", id)
	}
	return proposal, nil
}

func (pr *proposalUseCase) SavePropo(payload *model.Proposal, payloadDetail *model.ProposalDetail, org_profile multipart.File, propo_doc multipart.File) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }

	if payload.Status == "1" && (payload.OrgName == "" || payload.OrganizationTypeID == nil || payload.Email == "" || payload.Phone == "" || payload.PICName == "" || payload.City == "" || payload.PostalCode == "" || payload.Address == "" || payload.Description == "" || payloadDetail.ProjectName == "" || payloadDetail.PartnershipTypeID == nil || payloadDetail.StartDate.IsZero() || payloadDetail.EndDate.IsZero() || payloadDetail.Objective == "" || payloadDetail.Alignment == "") {
		return fmt.Errorf("form is not completed")
	}

	tx := pr.repo.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := pr.repo.Save(payload); err != nil {
		tx.Rollback()
		return err
	}

	payloadDetail.ProposalID = payload.ID
	if err := pr.propoDetailUC.SavePropoDetail(payloadDetail, tx); err != nil {
		tx.Rollback()
		return err
	}

	if org_profile != nil {
		fileUrl, err := pr.fileUC.FirebaseUpload(org_profile)
		if err != nil {
			tx.Rollback()
			return err
		}
		// Create the file payload
		filePayload := model.File{
			BaseModel:  model.BaseModel{},
			Label:      "organization profile",
			FileURL:    fileUrl,
			ProposalID: payload.ID,
		}
		if err := pr.fileUC.SaveTrx(&filePayload, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if propo_doc != nil {
		fileUrl, err := pr.fileUC.FirebaseUpload(propo_doc)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Create the file payload
		filePayload := model.File{
			BaseModel:  model.BaseModel{},
			Label:      "proposal document",
			FileURL:    fileUrl,
			ProposalID: payload.ID,
		}

		if err := pr.fileUC.SaveTrx(&filePayload, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if payload.Status == "1" {

		progressFilter := make(map[string]interface{})
		progressFilter["label"] = "received"

		progress, err := pr.progressUC.SearchBy(progressFilter)
		if err != nil {
			tx.Rollback()
			return err
		}

		propoProgressPayload := model.ProposalProgress{
			ProposalID: payload.ID,
			ProgressID: progress[0].ID,
			Note:       progress[0].Name,
			Status:     "1",
		}

		if err := pr.propoProgressUC.SaveTrx(&propoProgressPayload, tx); err != nil {
			tx.Rollback()
			return err
		}

		updateProposal := model.Proposal{
			BaseModel:       model.BaseModel{ID: payload.ID},
			CurrentProgress: propoProgressPayload.Note,
		}

		if err := pr.UpdateMinor(&updateProposal, tx); err != nil {
			tx.Rollback()
			return err
		}

		progressFilter = make(map[string]interface{})
		progressFilter["label"] = "review"

		progress, err = pr.progressUC.SearchBy(progressFilter)
		if err != nil {
			tx.Rollback()
			return err
		}

		propoProgressPayload = model.ProposalProgress{
			ProposalID: payload.ID,
			ProgressID: progress[0].ID,
			Note:       progress[0].Name,
			Status:     "0",
		}

		if err := pr.propoProgressUC.SaveTrx(&propoProgressPayload, tx); err != nil {
			tx.Rollback()
			return err
		}

		// assign least assigned admin as reviewer
		filter := make(map[string]interface{})
		filter["role"] = "admin"

		admins, err := pr.userUC.SearchBy(filter)
		if err != nil {
			tx.Rollback()
			return err
		}

		if len(admins) > 0 {
			proposalPayload := model.Proposal{
				BaseModel: model.BaseModel{
					ID: payload.ID,
				},
				ReviewerID: admins[0].ID,
			}
			if err := pr.UpdateMinor(&proposalPayload, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
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
	return pr.repo.UpdateBasic(payload)
}

func (pr *proposalUseCase) UpdatePropo(payload *model.Proposal, payloadDetail *model.ProposalDetail, org_profile multipart.File, propo_doc multipart.File) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }

	if payload.ID != "" {
		_, err := pr.FindById(payload.ID)
		if err != nil {
			return fmt.Errorf("proposal with ID %s not found", payload.ID)
		}
	}

	tx := pr.repo.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := pr.repo.Update(payload); err != nil {
		tx.Rollback()
		return err
	}

	if err := pr.propoDetailUC.UpdatePropoDetail(payloadDetail, tx); err != nil {
		tx.Rollback()
		return err
	}

	if org_profile != nil {
		if len(payload.File) > 0 {
			for _, v := range payload.File {
				if v.Label == "organization profile" {
					if err := pr.fileUC.DeleteDataTrx(v.ID, tx); err != nil {
						tx.Rollback()
						return err
					}
				}
			}
		}
		fileUrl, err := pr.fileUC.FirebaseUpload(org_profile)
		if err != nil {
			tx.Rollback()
			return err
		}
		// Create the file payload
		filePayload := model.File{
			BaseModel:  model.BaseModel{},
			Label:      "organization profile",
			FileURL:    fileUrl,
			ProposalID: payload.ID,
		}
		if err := pr.fileUC.SaveTrx(&filePayload, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if propo_doc != nil {
		if len(payload.File) > 0 {
			for _, v := range payload.File {
				if v.Label == "proposal document" {
					if err := pr.fileUC.DeleteDataTrx(v.ID, tx); err != nil {
						tx.Rollback()
						return err
					}
				}
			}
		}
		fileUrl, err := pr.fileUC.FirebaseUpload(propo_doc)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Create the file payload
		filePayload := model.File{
			BaseModel:  model.BaseModel{},
			Label:      "proposal document",
			FileURL:    fileUrl,
			ProposalID: payload.ID,
		}
		if err := pr.fileUC.SaveTrx(&filePayload, tx); err != nil {
			tx.Rollback()
			return err
		}

	}

	if payload.Status == "1" {

		progressFilter := make(map[string]interface{})
		progressFilter["label"] = "received"

		progress, err := pr.progressUC.SearchBy(progressFilter)
		if err != nil {
			tx.Rollback()
			return err
		}

		propoProgressPayload := model.ProposalProgress{
			ProposalID: payload.ID,
			ProgressID: progress[0].ID,
			Note:       progress[0].Name,
			Status:     "1",
		}

		if err := pr.propoProgressUC.SaveTrx(&propoProgressPayload, tx); err != nil {
			tx.Rollback()
			return err
		}

		updateProposal := model.Proposal{
			BaseModel:       model.BaseModel{ID: payload.ID},
			CurrentProgress: propoProgressPayload.Note,
		}

		if err := pr.UpdateMinor(&updateProposal, tx); err != nil {
			tx.Rollback()
			return err
		}

		progressFilter = make(map[string]interface{})
		progressFilter["label"] = "review"

		progress, err = pr.progressUC.SearchBy(progressFilter)
		if err != nil {
			tx.Rollback()
			return err
		}

		propoProgressPayload = model.ProposalProgress{
			ProposalID: payload.ID,
			ProgressID: progress[0].ID,
			Note:       progress[0].Name,
			Status:     "0",
		}

		if err := pr.propoProgressUC.SaveTrx(&propoProgressPayload, tx); err != nil {
			tx.Rollback()
			return err
		}

		// assign least assigned admin as reviewer
		filter := make(map[string]interface{})
		filter["role"] = "admin"

		admins, err := pr.userUC.SearchBy(filter)
		if err != nil {
			tx.Rollback()
			return err
		}

		if len(admins) > 0 {
			payload.ReviewerID = admins[0].ID
			if err := pr.UpdateMinor(payload, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	_, err := pr.FindByIdTx(payload.ID, tx)
	if err != nil {
		return fmt.Errorf("proposal with ID %s not found", payload.ID)
	}

	if payload.Status == "1" && (payload.OrgName == "" || payload.OrganizationTypeID == nil || payload.Email == "" || payload.Phone == "" || payload.PICName == "" || payload.City == "" || payload.PostalCode == "" || payload.Address == "" || payload.Description == "" || payloadDetail.ProjectName == "" || payloadDetail.PartnershipTypeID == nil || payloadDetail.StartDate.IsZero() || payloadDetail.EndDate.IsZero() || payloadDetail.Objective == "" || payloadDetail.Alignment == "") {
		return fmt.Errorf("form is not completed")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (pr *proposalUseCase) UpdateMinor(payload *model.Proposal, tx *gorm.DB) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }
	// cek jika data sudah ada -> count > 0

	return pr.repo.UpdateMinor(payload, tx)
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

func NewProposalUseCase(repo repository.ProposalRepository, propoDetailUC ProposalDetailUseCase, fileUC FileUseCase, progressUC ProgressUseCase, propoProgressUC ProposalProgressUseCase, userUC UsersUseCase) ProposalUseCase {
	return &proposalUseCase{
		repo:            repo,
		propoDetailUC:   propoDetailUC,
		fileUC:          fileUC,
		progressUC:      progressUC,
		propoProgressUC: propoProgressUC,
		userUC:          userUC,
	}
}
