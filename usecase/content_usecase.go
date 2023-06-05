package usecase

import (
	"fmt"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"
	"be-b-impact.com/csr/repository"
)

type ContentUseCase interface {
	BaseUseCase[model.Content]
	BaseUseCasePaging[model.Content]
}

type contentUseCase struct {
	repo repository.ContentRepository
}

func (co *contentUseCase) DeleteData(id string) error {
	content, err := co.FindById(id)
	if err != nil {
		return fmt.Errorf("content with ID %s not found", id)
	}
	return co.repo.Delete(content.ID)
}

func (co *contentUseCase) FindAll() ([]model.Content, error) {
	return co.repo.List()
}

func (co *contentUseCase) FindById(id string) (*model.Content, error) {
	content, err := co.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("content with ID %s not found", id)
	}
	return content, nil
}

func (co *contentUseCase) SaveData(payload *model.Content) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }
	// cek jika data sudah ada -> count > 0

	return co.repo.Save(payload)
}

func (co *contentUseCase) UpdateData(payload *model.Content) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }
	// cek jika data sudah ada -> count > 0

	if payload.ID != "" {
		_, err := co.FindById(payload.ID)
		if err != nil {
			return fmt.Errorf("content with ID %s not found", payload.ID)
		}
	}
	return co.repo.Update(payload)
}

func (co *contentUseCase) SearchBy(by map[string]interface{}) ([]model.Content, error) {
	contents, err := co.repo.Search(by)
	if err != nil {
		return nil, fmt.Errorf("request invalid")
	}
	return contents, nil
}

func (co *contentUseCase) Pagination(requestQueryParams dto.RequestQueryParams) ([]model.Content, dto.Paging, error) {
	if !requestQueryParams.QueryParams.IsSortValid() {
		return nil, dto.Paging{}, fmt.Errorf("invalid sort by: %s", requestQueryParams.QueryParams.Sort)
	}
	return co.repo.Paging(requestQueryParams)
}

func NewContentUseCase(repo repository.ContentRepository) ContentUseCase {
	return &contentUseCase{repo: repo}
}
