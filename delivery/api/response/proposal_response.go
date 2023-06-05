package response

import (
	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"
)

func MapProposalToResponse(proposal *model.Proposal, role string) dto.ProposalDTO {
	res := dto.ProposalDTO{
		ID:              proposal.ID,
		OrgName:         proposal.OrgName,
		OrganizatonType: proposal.OrganizatonType.Name,
		Email:           proposal.Email,
		Phone:           proposal.Phone,
		PICName:         proposal.PICName,
		City:            proposal.City,
		PostalCode:      proposal.PostalCode,
		Address:         proposal.Address,
		Description:     proposal.Description,
		Status:          proposal.Status,
		ProjectName:     proposal.ProposalDetail.ProjectName,
		PartnershipType: proposal.ProposalDetail.PartnershipType.Name,
		StartDate:       proposal.ProposalDetail.StartDate,
		EndDate:         proposal.ProposalDetail.EndDate,
		Objective:       proposal.ProposalDetail.Objective,
		Alignment:       proposal.ProposalDetail.Alignment,
		CreatedBy:       proposal.CreatedBy,
		Reviewer:        proposal.ReviewerID,
		CreatedAt:       proposal.CreatedAt,
		UpdatedAt:       proposal.UpdatedAt,
	}
	for _, v := range proposal.File {
		res.Files = append(res.Files, dto.FileDTO{
			Label:     v.Label,
			FileURL:   v.FileURL,
			CreatedAt: v.CreatedAt,
		})
	}

	for _, v := range proposal.ProposalProgress {
		if role == "admin" {
			res.Progresses = append(res.Progresses, dto.ProgressDTO{
				Name:      v.Progress.Name,
				Label:     v.Progress.Label,
				Status:    v.Status,
				Note:      v.Note,
				CreatedAt: v.CreatedAt,
				UpdatedAt: v.UpdatedAt,
			})
		} else {
			if v.Status == "1" {
				res.Progresses = append(res.Progresses, dto.ProgressDTO{
					Name:      v.Progress.Name,
					Label:     v.Progress.Label,
					Status:    v.Status,
					Note:      v.Note,
					CreatedAt: v.CreatedAt,
					UpdatedAt: v.UpdatedAt,
				})
			}
		}
	}
	lastArr := len(res.Progresses) - 1
	if role == "admin" && res.Progresses[lastArr].Label != "completed" {
		lastArr--
		res.CurrentProgress = dto.ProgressDTO{
			Name:      res.Progresses[lastArr].Name,
			Label:     res.Progresses[lastArr].Label,
			Status:    res.Progresses[lastArr].Status,
			Note:      res.Progresses[lastArr].Note,
			CreatedAt: res.Progresses[lastArr].CreatedAt,
			UpdatedAt: res.Progresses[lastArr].UpdatedAt,
		}
	} else {
		res.CurrentProgress = dto.ProgressDTO{
			Name:      res.Progresses[lastArr].Name,
			Label:     res.Progresses[lastArr].Label,
			Status:    res.Progresses[lastArr].Status,
			Note:      res.Progresses[lastArr].Note,
			CreatedAt: res.Progresses[lastArr].CreatedAt,
			UpdatedAt: res.Progresses[lastArr].UpdatedAt,
		}
	}

	return res
}
