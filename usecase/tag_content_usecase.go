package usecase

import (
	"fmt"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"
	"be-b-impact.com/csr/repository"
	"gorm.io/gorm"
)

type TagsContentUseCase interface {
	BaseUseCase[model.TagsContent]
	BaseUseCasePaging[model.TagsContent]
	SaveTagsContent(payload *model.TagsContent, tx *gorm.DB) error
	DeleteDataTrx(id string, tx *gorm.DB) error
}

type tagsContentUseCase struct {
	repo repository.TagsContentRepository
}

func (tc *tagsContentUseCase) DeleteData(id string) error {
	tagsContent, err := tc.FindById(id)
	if err != nil {
		return fmt.Errorf("tagsContent with ID %s not found", id)
	}
	return tc.repo.Delete(tagsContent.ID)
}

func (tc *tagsContentUseCase) DeleteDataTrx(id string, tx *gorm.DB) error {
	tagsContent, err := tc.FindById(id)
	if err != nil {
		return fmt.Errorf("tagsContent with ID %s not found", id)
	}
	return tc.repo.DeleteTrx(tagsContent.ID, tx)
}

func (tc *tagsContentUseCase) FindAll() ([]model.TagsContent, error) {
	return tc.repo.List()
}

func (tc *tagsContentUseCase) FindById(id string) (*model.TagsContent, error) {
	tagsContent, err := tc.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("tagsContent with ID %s not found", id)
	}
	return tagsContent, nil
}

func (tc *tagsContentUseCase) SaveData(payload *model.TagsContent) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }
	// cek jika data sudah ada -> count > 0

	if payload.ID != "" {
		_, err := tc.FindById(payload.ID)
		if err != nil {
			return fmt.Errorf("tagsContent with ID %s not found", payload.ID)
		}
	}
	return tc.repo.Save(payload)
}

func (tc *tagsContentUseCase) SaveTagsContent(payload *model.TagsContent, tx *gorm.DB) error {
	// Save the TagsContent using the provided transaction
	if err := tc.repo.SaveTrx(payload, tx); err != nil {
		return err
	}
	return nil
}

func (tc *tagsContentUseCase) UpdateData(payload *model.TagsContent) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }
	// cek jika data sudah ada -> count > 0

	if payload.ID != "" {
		_, err := tc.FindById(payload.ID)
		if err != nil {
			return fmt.Errorf("tagsContent with ID %s not found", payload.ID)
		}
	}
	return tc.repo.Update(payload)
}

func (tc *tagsContentUseCase) SearchBy(by map[string]interface{}) ([]model.TagsContent, error) {
	tagsContents, err := tc.repo.Search(by)
	if err != nil {
		return nil, fmt.Errorf("request invalid")
	}
	return tagsContents, nil
}

func (tc *tagsContentUseCase) Pagination(requestQueryParams dto.RequestQueryParams) ([]model.TagsContent, dto.Paging, error) {
	if !requestQueryParams.QueryParams.IsSortValid() {
		return nil, dto.Paging{}, fmt.Errorf("invalid sort by: %s", requestQueryParams.QueryParams.Sort)
	}
	return tc.repo.Paging(requestQueryParams)
}

func NewTagsContentUseCase(repo repository.TagsContentRepository) TagsContentUseCase {
	return &tagsContentUseCase{repo: repo}
}
