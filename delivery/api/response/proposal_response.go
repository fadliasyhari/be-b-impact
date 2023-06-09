package response

import (
	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"
)

func MapProposalToResponse(proposal *model.Proposal) dto.ProposalDTO {
	res := dto.ProposalDTO{
		ID:      proposal.ID,
		OrgName: proposal.OrgName,
		OrganizatonType: dto.CategoryDTO{
			ID:   proposal.OrganizatonType.ID,
			Name: proposal.OrganizatonType.Name,
		},
		Email:            proposal.Email,
		Phone:            proposal.Phone,
		PICName:          proposal.PICName,
		City:             proposal.City,
		PostalCode:       proposal.PostalCode,
		Address:          proposal.Address,
		Description:      proposal.Description,
		Status:           proposal.Status,
		CurrentProgress:  proposal.CurrentProgress,
		ProposalDetailID: proposal.ProposalDetail.ID,
		ProjectName:      proposal.ProposalDetail.ProjectName,
		PartnershipType: dto.CategoryDTO{
			ID:   proposal.ProposalDetail.PartnershipType.ID,
			Name: proposal.ProposalDetail.PartnershipType.Name,
		},
		StartDate:  proposal.ProposalDetail.StartDate,
		EndDate:    proposal.ProposalDetail.EndDate,
		Objective:  proposal.ProposalDetail.Objective,
		Alignment:  proposal.ProposalDetail.Alignment,
		CreatedBy:  proposal.CreatedBy,
		Reviewer:   proposal.ReviewerID,
		Files:      []dto.FileDTO{},
		Progresses: []dto.ProgressDTO{},
		CreatedAt:  proposal.CreatedAt,
		UpdatedAt:  proposal.UpdatedAt,
	}
	for _, v := range proposal.File {
		res.Files = append(res.Files, dto.FileDTO{
			ID:        v.ID,
			Label:     v.Label,
			FileURL:   v.FileURL,
			CreatedAt: v.CreatedAt,
		})
	}

	for _, v := range proposal.ProposalProgress {
		if v.Progress.Label == "approved" {
			res.Progresses = append(res.Progresses, dto.ProgressDTO{
				ID:             v.ID,
				Name:           v.Progress.Name,
				Label:          v.Progress.Label,
				Status:         v.Status,
				Note:           v.Note,
				ReviewLocation: v.ReviewLocation,
				ReviewDate:     v.ReviewDate,
				ReviewCP:       v.ReviewCP,
				CreatedAt:      v.CreatedAt,
				UpdatedAt:      v.UpdatedAt,
			})
		} else if v.Progress.Label == "rejected" {
			res.Progresses = append(res.Progresses, dto.ProgressDTO{
				ID:             v.ID,
				Name:           v.Progress.Name,
				Label:          v.Progress.Label,
				Status:         v.Status,
				Note:           v.Note,
				ReviewFeedback: v.ReviewFeedback,
				CreatedAt:      v.CreatedAt,
				UpdatedAt:      v.UpdatedAt,
			})
		} else {
			res.Progresses = append(res.Progresses, dto.ProgressDTO{
				ID:        v.ID,
				Name:      v.Progress.Name,
				Label:     v.Progress.Label,
				Status:    v.Status,
				Note:      v.Note,
				CreatedAt: v.CreatedAt,
				UpdatedAt: v.UpdatedAt,
			})
		}

	}
	return res
}
