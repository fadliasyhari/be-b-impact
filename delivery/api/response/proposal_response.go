package response

import (
	"be-b-impact.com/csr/model/dto"
)

func MapProposalToResponse(proposal dto.Proposal, role string) dto.ProposalDTO {
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
		CurrentProgress: proposal.Current,
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

	for i, v := range proposal.ProposalProgress {
		if (role == "admin" || role == "super") && len(proposal.ProposalProgress)-1 == i && v.Progress.Label != "rejected" && v.Status == "0" {
			res.Progresses = append(res.Progresses, dto.ProgressDTO{
				ID:        v.ID,
				Name:      v.Progress.Name,
				Label:     v.Progress.Label,
				Status:    v.Status,
				Note:      v.Note,
				CreatedAt: v.CreatedAt,
				UpdatedAt: v.UpdatedAt,
			})
		} else if v.Status == "1" {
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
	}
	return res
}

func MapProposalToSingleResponse(proposal *dto.Proposal, role string) dto.ProposalDTO {
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
		CurrentProgress: proposal.Current,
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

	for i, v := range proposal.ProposalProgress {
		if (role == "admin" || role == "super") && len(proposal.ProposalProgress)-1 == i && v.Progress.Label != "rejected" && v.Status == "0" {
			res.Progresses = append(res.Progresses, dto.ProgressDTO{
				ID:        v.ID,
				Name:      v.Progress.Name,
				Label:     v.Progress.Label,
				Status:    v.Status,
				Note:      v.Note,
				CreatedAt: v.CreatedAt,
				UpdatedAt: v.UpdatedAt,
			})
		} else if v.Status == "1" {
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
	}

	return res
}
