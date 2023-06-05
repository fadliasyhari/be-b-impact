package usecase

import (
	"fmt"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"
	"be-b-impact.com/csr/repository"
)

type CategoryUseCase interface {
	BaseUseCase[model.Category]
	BaseUseCasePaging[model.Category]
}

type categoryUseCase struct {
	repo repository.CategoryRepository
}

func (ca *categoryUseCase) DeleteData(id string) error {
	category, err := ca.FindById(id)
	if err != nil {
		return fmt.Errorf("category with ID %s not found", id)
	}
	return ca.repo.Delete(category.ID)
}

func (ca *categoryUseCase) FindAll() ([]model.Category, error) {
	return ca.repo.List()
}

func (ca *categoryUseCase) FindById(id string) (*model.Category, error) {
	category, err := ca.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("category with ID %s not found", id)
	}
	return category, nil
}

func (ca *categoryUseCase) SaveData(payload *model.Category) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }
	// cek jika data sudah ada -> count > 0
	err := ca.repo.CountData(payload.Name, payload.ID)
	if err != nil {
		return err
	}

	if payload.ID != "" {
		_, err := ca.FindById(payload.ID)
		if err != nil {
			return fmt.Errorf("category with ID %s not found", payload.ID)
		}
	}
	return ca.repo.Save(payload)
}

func (ca *categoryUseCase) UpdateData(payload *model.Category) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }
	// cek jika data sudah ada -> count > 0

	if payload.Name != "" {
		err := ca.repo.CountData(payload.Name, payload.ID)
		if err != nil {
			return err
		}
	}

	if payload.ID != "" {
		_, err := ca.FindById(payload.ID)
		if err != nil {
			return fmt.Errorf("category with ID %s not found", payload.ID)
		}
	}
	return ca.repo.Update(payload)
}

func (ca *categoryUseCase) SearchBy(by map[string]interface{}) ([]model.Category, error) {
	categorys, err := ca.repo.Search(by)
	if err != nil {
		return nil, fmt.Errorf("request invalid")
	}
	return categorys, nil
}

func (ca *categoryUseCase) Pagination(requestQueryParams dto.RequestQueryParams) ([]model.Category, dto.Paging, error) {
	if !requestQueryParams.QueryParams.IsSortValid() {
		return nil, dto.Paging{}, fmt.Errorf("invalid sort by: %s", requestQueryParams.QueryParams.Sort)
	}
	return ca.repo.Paging(requestQueryParams)
}

func NewCategoryUseCase(repo repository.CategoryRepository) CategoryUseCase {
	return &categoryUseCase{repo: repo}
}
