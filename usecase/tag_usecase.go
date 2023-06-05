package usecase

import (
	"fmt"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"
	"be-b-impact.com/csr/repository"
)

type TagUseCase interface {
	BaseUseCase[model.Tag]
	BaseUseCasePaging[model.Tag]
}

type tagUseCase struct {
	repo repository.TagRepository
}

func (ta *tagUseCase) DeleteData(id string) error {
	tag, err := ta.FindById(id)
	if err != nil {
		return fmt.Errorf("tag with ID %s not found", id)
	}
	return ta.repo.Delete(tag.ID)
}

func (ta *tagUseCase) FindAll() ([]model.Tag, error) {
	return ta.repo.List()
}

func (ta *tagUseCase) FindById(id string) (*model.Tag, error) {
	tag, err := ta.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("tag with ID %s not found", id)
	}
	return tag, nil
}

func (ta *tagUseCase) SaveData(payload *model.Tag) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }
	// cek jika data sudah ada -> count > 0
	err := ta.repo.CountData(payload.Name, payload.ID)
	if err != nil {
		return err
	}

	if payload.ID != "" {
		_, err := ta.FindById(payload.ID)
		if err != nil {
			return fmt.Errorf("tag with ID %s not found", payload.ID)
		}
	}
	return ta.repo.Save(payload)
}

func (ta *tagUseCase) UpdateData(payload *model.Tag) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }
	// cek jika data sudah ada -> count > 0

	if payload.Name != "" {
		err := ta.repo.CountData(payload.Name, payload.ID)
		if err != nil {
			return err
		}
	}

	if payload.ID != "" {
		_, err := ta.FindById(payload.ID)
		if err != nil {
			return fmt.Errorf("tag with ID %s not found", payload.ID)
		}
	}
	return ta.repo.Update(payload)
}

func (ta *tagUseCase) SearchBy(by map[string]interface{}) ([]model.Tag, error) {
	tags, err := ta.repo.Search(by)
	if err != nil {
		return nil, fmt.Errorf("request invalid")
	}
	return tags, nil
}

func (ta *tagUseCase) Pagination(requestQueryParams dto.RequestQueryParams) ([]model.Tag, dto.Paging, error) {
	if !requestQueryParams.QueryParams.IsSortValid() {
		return nil, dto.Paging{}, fmt.Errorf("invalid sort by: %s", requestQueryParams.QueryParams.Sort)
	}
	return ta.repo.Paging(requestQueryParams)
}

func NewTagUseCase(repo repository.TagRepository) TagUseCase {
	return &tagUseCase{repo: repo}
}
