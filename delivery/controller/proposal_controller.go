package controller

import (
	"net/http"
	"strconv"
	"time"

	"be-b-impact.com/csr/delivery/api"
	"be-b-impact.com/csr/delivery/api/middleware"
	"be-b-impact.com/csr/delivery/api/response"
	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"
	"be-b-impact.com/csr/usecase"
	"be-b-impact.com/csr/utils"
	"github.com/gin-gonic/gin"
)

type ProposalController struct {
	router          *gin.Engine
	useCase         usecase.ProposalUseCase
	propoDetailUC   usecase.ProposalDetailUseCase
	fileUC          usecase.FileUseCase
	progressUC      usecase.ProgressUseCase
	propoProgressUC usecase.ProposalProgressUseCase
	userUC          usecase.UsersUseCase
	api.BaseApi
}

func (pr *ProposalController) createHandler(c *gin.Context) {
	userTyped := utils.AccessInsideToken(pr.BaseApi, c)

	// Parse the form data
	if err := c.Request.ParseMultipartForm(10 * 1024 * 1024); err != nil {
		pr.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get the form values
	orgName := c.Request.FormValue("org_name")
	orgTypeID := c.Request.FormValue("org_type_id")
	email := c.Request.FormValue("email")
	phone := c.Request.FormValue("phone")
	picName := c.Request.FormValue("pic_name")
	city := c.Request.FormValue("city")
	postalCode := c.Request.FormValue("postal_code")
	address := c.Request.FormValue("address")
	description := c.Request.FormValue("description")
	status := c.Request.FormValue("status")

	// Create the proposal payload
	proposalPayload := model.Proposal{
		OrgName:            orgName,
		OrganizationTypeID: orgTypeID,
		Email:              email,
		Phone:              phone,
		PICName:            picName,
		City:               city,
		PostalCode:         postalCode,
		Address:            address,
		Description:        description,
		Status:             status,
		CreatedBy:          userTyped.UserId,
	}

	if err := pr.useCase.SaveData(&proposalPayload); err != nil {
		pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	projectName := c.Request.FormValue("project_name")
	partTypeID := c.Request.FormValue("part_type_id")
	startDateStr := c.Request.FormValue("start_date")
	endDateStr := c.Request.FormValue("end_date")
	objective := c.Request.FormValue("objective")
	alignment := c.Request.FormValue("alignment")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		pr.NewFailedResponse(c, http.StatusBadRequest, "failed to parse start date")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		pr.NewFailedResponse(c, http.StatusBadRequest, "failed to parse end date")
		return
	}

	// Create the proposal detail payload
	proposalDetailPayload := model.ProposalDetail{
		BaseModel:         model.BaseModel{},
		ProposalID:        proposalPayload.ID,
		ProjectName:       projectName,
		PartnershipTypeID: partTypeID,
		StartDate:         startDate,
		EndDate:           endDate,
		Objective:         objective,
		Alignment:         alignment,
	}

	if err := pr.propoDetailUC.SaveData(&proposalDetailPayload); err != nil {
		pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	org_profile, _, err := c.Request.FormFile("org_profile")
	if err == nil {
		fileUrl, err := pr.fileUC.FirebaseUpload(org_profile)
		if err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		// Create the file payload
		filePayload := model.File{
			BaseModel:  model.BaseModel{},
			Label:      "organization profile",
			FileURL:    fileUrl,
			ProposalID: proposalPayload.ID,
		}
		if err := pr.fileUC.SaveData(&filePayload); err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

	}

	propo_doc, _, err := c.Request.FormFile("propo_doc")
	if err == nil {
		fileUrl, err := pr.fileUC.FirebaseUpload(propo_doc)
		if err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		// Create the file payload
		filePayload := model.File{
			BaseModel:  model.BaseModel{},
			Label:      "proposal document",
			FileURL:    fileUrl,
			ProposalID: proposalPayload.ID,
		}
		if err := pr.fileUC.SaveData(&filePayload); err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

	}

	if proposalPayload.Status == "1" {

		progressFilter := make(map[string]interface{})
		progressFilter["label"] = "received"

		progress, err := pr.progressUC.SearchBy(progressFilter)
		if err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		propoProgressPayload := model.ProposalProgress{
			ProposalID: proposalPayload.ID,
			ProgressID: progress[0].ID,
			Note:       progress[0].Name,
			Status:     "1",
		}

		if err := pr.propoProgressUC.SaveData(&propoProgressPayload); err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		updateProposal := model.Proposal{
			BaseModel:       model.BaseModel{ID: proposalPayload.ID},
			CurrentProgress: propoProgressPayload.Note,
		}

		if err := pr.useCase.UpdateData(&updateProposal); err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		progressFilter = make(map[string]interface{})
		progressFilter["label"] = "review"

		progress, err = pr.progressUC.SearchBy(progressFilter)
		if err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		propoProgressPayload = model.ProposalProgress{
			ProposalID: proposalPayload.ID,
			ProgressID: progress[0].ID,
			Note:       progress[0].Name,
			Status:     "0",
		}

		if err := pr.propoProgressUC.SaveData(&propoProgressPayload); err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		// assign least assigned admin as reviewer
		filter := make(map[string]interface{})
		filter["role"] = "admin"

		admins, err := pr.userUC.SearchBy(filter)
		if err != nil {
			pr.NewFailedResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		if len(admins) > 0 {
			proposalPayload.ReviewerID = admins[0].ID
			if err := pr.useCase.UpdateData(&proposalPayload); err != nil {
				pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
				return
			}
		}
	}

	pr.NewSuccessSingleResponse(c, "OK", proposalPayload)
}

func (pr *ProposalController) listHandler(c *gin.Context) {
	userTyped := utils.AccessInsideToken(pr.BaseApi, c)
	filter := make(map[string]interface{})

	// Iterate over the query parameters
	for key, values := range c.Request.URL.Query() {
		// Skip if the key is empty or has multiple values
		if key == "page" || key == "limit" || key == "order" || key == "sort" || key == "" || len(values) != 1 {
			continue
		}

		// Add key-value pair to the filter
		filter[key] = values[0]
	}
	if userTyped.Role == "member" {
		filter["created_by"] = userTyped.UserId
	}
	if userTyped.Role == "admin" {
		filter["reviewer_id"] = userTyped.UserId
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		pr.NewFailedResponse(c, http.StatusBadRequest, "invalid page number")
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		pr.NewFailedResponse(c, http.StatusBadRequest, "invalid limit number")
		return
	}
	order := c.DefaultQuery("order", "id")
	sort := c.DefaultQuery("sort", "ASC")
	requestQueryParams := dto.RequestQueryParams{
		QueryParams: dto.QueryParams{
			Sort:  sort,
			Order: order,
		},
		PaginationParam: dto.PaginationParam{
			Page:  page,
			Limit: limit,
		},
		Filter: filter,
	}
	proposal, paging, err := pr.useCase.Pagination(requestQueryParams)
	if err != nil {
		pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	var proposalInterface []interface{}
	for _, pr := range proposal {
		res := response.MapProposalToResponse(&pr)
		proposalInterface = append(proposalInterface, res)
	}
	pr.NewSuccessPagedResponse(c, "OK", proposalInterface, paging)
}

func (pr *ProposalController) getHandler(c *gin.Context) {
	userTyped := utils.AccessInsideToken(pr.BaseApi, c)

	id := c.Param("id")
	proposal, err := pr.useCase.FindById(id)
	if err != nil {
		pr.NewFailedResponse(c, http.StatusNotFound, err.Error())
		return
	}
	if userTyped.Role != "super" {
		if userTyped.UserId != proposal.CreatedBy && userTyped.UserId != proposal.ReviewerID {
			pr.NewFailedResponse(c, http.StatusForbidden, "access denied")
			return
		}
	}
	res := response.MapProposalToResponse(proposal)
	pr.NewSuccessSingleResponse(c, "OK", res)
}

func (pr *ProposalController) searchHandler(c *gin.Context) {
	filter := make(map[string]interface{})

	// Iterate over the query parameters
	for key, values := range c.Request.URL.Query() {
		// Skip if the key is empty or has multiple values
		if key == "" || len(values) != 1 {
			continue
		}

		// Add key-value pair to the filter
		filter[key] = values[0]
	}
	proposal, err := pr.useCase.SearchBy(filter)
	if err != nil {
		pr.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	pr.NewSuccessSingleResponse(c, "OK", proposal)
}

func (pr *ProposalController) updateHandler(c *gin.Context) {
	userTyped := utils.AccessInsideToken(pr.BaseApi, c)
	// Parse the form data
	if err := c.Request.ParseMultipartForm(10 * 1024 * 1024); err != nil {
		pr.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id := c.Request.FormValue("id")
	orgName := c.Request.FormValue("org_name")
	orgTypeID := c.Request.FormValue("org_type_id")
	email := c.Request.FormValue("email")
	phone := c.Request.FormValue("phone")
	picName := c.Request.FormValue("pic_name")
	city := c.Request.FormValue("city")
	postalCode := c.Request.FormValue("postal_code")
	address := c.Request.FormValue("address")
	description := c.Request.FormValue("description")
	status := c.Request.FormValue("status")

	existingProposal, err := pr.useCase.FindById(id)
	if err != nil {
		pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if userTyped.Role != "super" {
		if userTyped.UserId != existingProposal.CreatedBy && userTyped.UserId != existingProposal.ReviewerID {
			pr.NewFailedResponse(c, http.StatusForbidden, "access denied")
			return
		}
	}

	existingProposal.OrgName = orgName
	existingProposal.OrganizationTypeID = orgTypeID
	existingProposal.Email = email
	existingProposal.Phone = phone
	existingProposal.PICName = picName
	existingProposal.City = city
	existingProposal.PostalCode = postalCode
	existingProposal.Address = address
	existingProposal.Description = description
	existingProposal.Status = status

	if err := pr.useCase.UpdateData(existingProposal); err != nil {
		pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	projectName := c.Request.FormValue("project_name")
	partTypeID := c.Request.FormValue("part_type_id")
	startDateStr := c.Request.FormValue("start_date")
	endDateStr := c.Request.FormValue("end_date")
	objective := c.Request.FormValue("objective")
	alignment := c.Request.FormValue("alignment")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		pr.NewFailedResponse(c, http.StatusBadRequest, "failed to parse start date")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		pr.NewFailedResponse(c, http.StatusBadRequest, "failed to parse end date")
		return
	}

	filter := make(map[string]interface{})
	filter["proposal_id"] = existingProposal.ID

	proposalDetail, err := pr.propoDetailUC.SearchBy(filter)
	if err != nil {
		pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Create the proposal detail payload
	proposalDetailPayload := model.ProposalDetail{
		BaseModel: model.BaseModel{
			ID: proposalDetail[0].ID,
		},
		ProjectName:       projectName,
		PartnershipTypeID: partTypeID,
		StartDate:         startDate,
		EndDate:           endDate,
		Objective:         objective,
		Alignment:         alignment,
	}

	if err := pr.propoDetailUC.UpdateData(&proposalDetailPayload); err != nil {
		pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	org_profile, _, err := c.Request.FormFile("org_profile")
	if err == nil {
		if len(existingProposal.File) > 0 {
			for _, v := range existingProposal.File {
				if v.Label == "organization profile" {
					if err := pr.fileUC.DeleteData(v.ID); err != nil {
						pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
						return
					}
				}
			}
		}
		fileUrl, err := pr.fileUC.FirebaseUpload(org_profile)
		if err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		// Create the file payload
		filePayload := model.File{
			BaseModel:  model.BaseModel{},
			Label:      "organization profile",
			FileURL:    fileUrl,
			ProposalID: existingProposal.ID,
		}
		if err := pr.fileUC.SaveData(&filePayload); err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

	}

	propo_doc, _, err := c.Request.FormFile("propo_doc")
	if err == nil {
		if len(existingProposal.File) > 0 {
			for _, v := range existingProposal.File {
				if v.Label == "proposal document" {
					if err := pr.fileUC.DeleteData(v.ID); err != nil {
						pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
						return
					}
				}
			}
		}
		fileUrl, err := pr.fileUC.FirebaseUpload(propo_doc)
		if err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		// Create the file payload
		filePayload := model.File{
			BaseModel:  model.BaseModel{},
			Label:      "proposal document",
			FileURL:    fileUrl,
			ProposalID: existingProposal.ID,
		}
		if err := pr.fileUC.SaveData(&filePayload); err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

	}

	if existingProposal.Status == "1" {

		progressFilter := make(map[string]interface{})
		progressFilter["label"] = "received"

		progress, err := pr.progressUC.SearchBy(progressFilter)
		if err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		propoProgressPayload := model.ProposalProgress{
			ProposalID: existingProposal.ID,
			ProgressID: progress[0].ID,
			Note:       progress[0].Name,
			Status:     "1",
		}

		if err := pr.propoProgressUC.SaveData(&propoProgressPayload); err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		updateProposal := model.Proposal{
			BaseModel:       model.BaseModel{ID: existingProposal.ID},
			CurrentProgress: propoProgressPayload.Note,
		}

		if err := pr.useCase.UpdateData(&updateProposal); err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		progressFilter = make(map[string]interface{})
		progressFilter["label"] = "review"

		progress, err = pr.progressUC.SearchBy(progressFilter)
		if err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		propoProgressPayload = model.ProposalProgress{
			ProposalID: existingProposal.ID,
			ProgressID: progress[0].ID,
			Note:       progress[0].Name,
			Status:     "0",
		}

		if err := pr.propoProgressUC.SaveData(&propoProgressPayload); err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		// assign least assigned admin as reviewer
		filter := make(map[string]interface{})
		filter["role"] = "admin"

		admins, err := pr.userUC.SearchBy(filter)
		if err != nil {
			pr.NewFailedResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		if len(admins) > 0 {
			existingProposal.ReviewerID = admins[0].ID
			if err := pr.useCase.UpdateData(existingProposal); err != nil {
				pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
				return
			}
		}
	}

	pr.NewSuccessSingleResponse(c, "OK", existingProposal)
}

func (pr *ProposalController) deleteHandler(c *gin.Context) {
	userTyped := utils.AccessInsideToken(pr.BaseApi, c)
	id := c.Param("id")
	proposal, err := pr.useCase.FindById(id)
	if err != nil {
		pr.NewFailedResponse(c, http.StatusNotFound, err.Error())
		return
	}
	if userTyped.Role != "super" {
		if userTyped.UserId != proposal.CreatedBy && userTyped.UserId != proposal.ReviewerID {
			pr.NewFailedResponse(c, http.StatusForbidden, "access denied")
			return
		}
	}

	err = pr.useCase.DeleteData(id)
	if err != nil {
		pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "Success Delete")
}

func (pr *ProposalController) createProgressHandler(c *gin.Context) {
	userTyped := utils.AccessInsideToken(pr.BaseApi, c)
	if userTyped.Role != "admin" && userTyped.Role != "super" {
		pr.NewFailedResponse(c, http.StatusForbidden, "access denied")
		return
	}

	var payload model.ProposalProgress
	if err := pr.ParseRequestBody(c, &payload); err != nil {
		pr.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	progress, err := pr.progressUC.FindById(payload.ProgressID)
	if err != nil {
		pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	payload.Note = progress.Name

	if err := pr.propoProgressUC.SaveData(&payload); err != nil {
		pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	pr.NewSuccessSingleResponse(c, "OK", payload)
}

func (pr *ProposalController) updateProgressHandler(c *gin.Context) {
	userTyped := utils.AccessInsideToken(pr.BaseApi, c)
	if userTyped.Role != "admin" && userTyped.Role != "super" {
		pr.NewFailedResponse(c, http.StatusForbidden, "access denied")
		return
	}

	type ProgressReq struct {
		ProposalId         string    `json:"proposal_id"`
		ProposalProgressId string    `json:"proposal_progress_id"`
		Progress           string    `json:"progress"`
		ReviewLocation     string    `json:"review_location"`
		ReviewCP           string    `json:"review_cp"`
		ReviewDate         time.Time `json:"review_date"`
		ReviewFeedback     string    `json:"review_feedback"`
	}

	var payload ProgressReq
	if err := pr.ParseRequestBody(c, &payload); err != nil {
		pr.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	proposal, err := pr.useCase.FindById(payload.ProposalId)
	if err != nil {
		pr.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	// update to review
	if payload.Progress == "review" {
		updatePayload := model.ProposalProgress{
			BaseModel: model.BaseModel{
				ID: payload.ProposalProgressId,
			},
			Status: "1",
		}
		if err := pr.propoProgressUC.UpdateData(&updatePayload); err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		updateProposal := model.Proposal{
			BaseModel:       model.BaseModel{ID: payload.ProposalId},
			CurrentProgress: "Under Review",
		}

		if err := pr.useCase.UpdateData(&updateProposal); err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		progressFilter := make(map[string]interface{})
		progressFilter["label"] = "decision"

		progress, err := pr.progressUC.SearchBy(progressFilter)
		if err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		propoProgressPayload := model.ProposalProgress{
			ProposalID: proposal.ID,
			ProgressID: progress[0].ID,
			Note:       progress[0].Name,
			Status:     "0",
		}

		if err := pr.propoProgressUC.SaveData(&propoProgressPayload); err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}
	// update to approved
	if payload.Progress == "approved" {
		if err := pr.propoProgressUC.DeleteData(payload.ProposalProgressId); err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		// add approved progress
		progressFilter := make(map[string]interface{})
		progressFilter["label"] = payload.Progress

		progress, err := pr.progressUC.SearchBy(progressFilter)
		if err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		propoProgressPayload := model.ProposalProgress{
			ProposalID:     proposal.ID,
			ProgressID:     progress[0].ID,
			Note:           progress[0].Name,
			Status:         "1",
			ReviewLocation: payload.ReviewLocation,
			ReviewCP:       payload.ReviewCP,
			ReviewDate:     payload.ReviewDate,
		}

		if err := pr.propoProgressUC.SaveData(&propoProgressPayload); err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		updateProposal := model.Proposal{
			BaseModel:       model.BaseModel{ID: payload.ProposalId},
			CurrentProgress: propoProgressPayload.Note,
		}

		if err := pr.useCase.UpdateData(&updateProposal); err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		// add in progress
		progressFilter = make(map[string]interface{})
		progressFilter["label"] = "inProgress"

		progress, err = pr.progressUC.SearchBy(progressFilter)
		if err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		propoProgressPayload = model.ProposalProgress{
			ProposalID: proposal.ID,
			ProgressID: progress[0].ID,
			Note:       progress[0].Name,
			Status:     "0",
		}

		if err := pr.propoProgressUC.SaveData(&propoProgressPayload); err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}
	// update if rejected
	if payload.Progress == "rejected" {
		if err := pr.propoProgressUC.DeleteData(payload.ProposalProgressId); err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		// add approved progress
		progressFilter := make(map[string]interface{})
		progressFilter["label"] = payload.Progress

		progress, err := pr.progressUC.SearchBy(progressFilter)
		if err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		propoProgressPayload := model.ProposalProgress{
			ProposalID:     proposal.ID,
			ProgressID:     progress[0].ID,
			Note:           progress[0].Name,
			Status:         "1",
			ReviewFeedback: payload.ReviewFeedback,
		}

		if err := pr.propoProgressUC.SaveData(&propoProgressPayload); err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		updateProposal := model.Proposal{
			BaseModel:       model.BaseModel{ID: payload.ProposalId},
			CurrentProgress: propoProgressPayload.Note,
		}

		if err := pr.useCase.UpdateData(&updateProposal); err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		// updateProposal := model.Proposal{
		// 	BaseModel: model.BaseModel{
		// 		ID: proposal.ID,
		// 	},
		// 	Status: "0",
		// }
		// if err := pr.useCase.UpdateData(&updateProposal); err != nil {
		// 	pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
		// 	return
		// }
	}

	// update if inProgress
	if payload.Progress == "inProgress" {

		updatePayload := model.ProposalProgress{
			BaseModel: model.BaseModel{
				ID: payload.ProposalProgressId,
			},
			Status: "1",
		}

		if err := pr.propoProgressUC.UpdateData(&updatePayload); err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		updateProposal := model.Proposal{
			BaseModel:       model.BaseModel{ID: payload.ProposalId},
			CurrentProgress: "Project In Progress",
		}

		if err := pr.useCase.UpdateData(&updateProposal); err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		progressFilter := make(map[string]interface{})
		progressFilter["label"] = "completed"

		progress, err := pr.progressUC.SearchBy(progressFilter)
		if err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		propoProgressPayload := model.ProposalProgress{
			ProposalID: proposal.ID,
			ProgressID: progress[0].ID,
			Note:       progress[0].Name,
			Status:     "0",
		}

		if err := pr.propoProgressUC.SaveData(&propoProgressPayload); err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

	}

	// update to completed
	if payload.Progress == "completed" {
		updatePayload := model.ProposalProgress{
			BaseModel: model.BaseModel{
				ID: payload.ProposalProgressId,
			},
			Status: "1",
		}

		if err := pr.propoProgressUC.UpdateData(&updatePayload); err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		updateProposal := model.Proposal{
			BaseModel:       model.BaseModel{ID: payload.ProposalId},
			CurrentProgress: "Project Completed",
		}

		if err := pr.useCase.UpdateData(&updateProposal); err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

	}
	pr.NewSuccessSingleResponse(c, "OK", payload)
}

func (pr *ProposalController) uploadFileHandler(c *gin.Context) {
	userTyped := utils.AccessInsideToken(pr.BaseApi, c)
	if userTyped.Role != "admin" && userTyped.Role != "super" {
		pr.NewFailedResponse(c, http.StatusForbidden, "access denied")
		return
	}

	// Parse the form data
	if err := c.Request.ParseMultipartForm(10 * 1024 * 1024); err != nil {
		pr.NewFailedResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	file, _, err := c.Request.FormFile("file")
	proposalId := c.Request.FormValue("proposal_id")
	label := c.Request.FormValue("label")
	proposalProgressId := c.Request.FormValue("proposal_progress_id")
	if label == "" {
		label = "accountable_report"
	}

	if err == nil {
		fileUrl, err := pr.fileUC.FirebaseUpload(file)
		if err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		// Create the file payload
		filePayload := model.File{
			BaseModel:  model.BaseModel{},
			Label:      label,
			FileURL:    fileUrl,
			ProposalID: proposalId,
		}
		if err := pr.fileUC.SaveData(&filePayload); err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		// auto complete progress
		updatePayload := model.ProposalProgress{
			BaseModel: model.BaseModel{
				ID: proposalProgressId,
			},
			Status: "1",
		}

		if err := pr.propoProgressUC.UpdateData(&updatePayload); err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		updateProposal := model.Proposal{
			BaseModel:       model.BaseModel{ID: proposalId},
			CurrentProgress: "Project Completed",
		}

		if err := pr.useCase.UpdateData(&updateProposal); err != nil {
			pr.NewFailedResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		pr.NewSuccessSingleResponse(c, "OK", filePayload)
	}
}

func NewProposalController(r *gin.Engine, useCase usecase.ProposalUseCase, propoDetailUC usecase.ProposalDetailUseCase, fileUC usecase.FileUseCase, progressUC usecase.ProgressUseCase, propoProgressUC usecase.ProposalProgressUseCase, userUC usecase.UsersUseCase, tokenMdw middleware.AuthTokenMiddlerware) *ProposalController {
	controller := &ProposalController{
		router:          r,
		useCase:         useCase,
		propoDetailUC:   propoDetailUC,
		fileUC:          fileUC,
		progressUC:      progressUC,
		propoProgressUC: propoProgressUC,
		userUC:          userUC,
	}
	proposalGroup := r.Group("/proposal", tokenMdw.RequireToken())
	{
		proposalGroup.GET("", controller.listHandler)
		proposalGroup.GET("/:id", controller.getHandler)
		proposalGroup.GET("/search", controller.searchHandler)
		proposalGroup.POST("", controller.createHandler)
		proposalGroup.POST("/progress", controller.createProgressHandler)
		proposalGroup.PUT("/progress", controller.updateProgressHandler)
		proposalGroup.POST("/file", controller.uploadFileHandler)
		proposalGroup.PUT("", controller.updateHandler)
		proposalGroup.DELETE("/:id", controller.deleteHandler)
	}
	return controller
}
