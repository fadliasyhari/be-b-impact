package usecase

import (
	"fmt"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"
	"be-b-impact.com/csr/repository"
)

type ProgressUseCase interface {
	BaseUseCase[model.Progress]
	BaseUseCasePaging[model.Progress]
}

type progressUseCase struct {
	repo repository.ProgressRepository
}

func (pg *progressUseCase) DeleteData(id string) error {
	progress, err := pg.FindById(id)
	if err != nil {
		return fmt.Errorf("progress with ID %s not found", id)
	}
	return pg.repo.Delete(progress.ID)
}

func (pg *progressUseCase) FindAll() ([]model.Progress, error) {
	return pg.repo.List()
}

func (pg *progressUseCase) FindById(id string) (*model.Progress, error) {
	progress, err := pg.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("progress with ID %s not found", id)
	}
	return progress, nil
}

func (pg *progressUseCase) SaveData(payload *model.Progress) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }
	// cek jika data sudah ada -> count > 0
	err := pg.repo.CountData(payload.Name, payload.ID)
	if err != nil {
		return err
	}

	if payload.ID != "" {
		_, err := pg.FindById(payload.ID)
		if err != nil {
			return fmt.Errorf("progress with ID %s not found", payload.ID)
		}
	}
	return pg.repo.Save(payload)
}

func (pg *progressUseCase) UpdateData(payload *model.Progress) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }
	// cek jika data sudah ada -> count > 0

	if payload.Name != "" {
		err := pg.repo.CountData(payload.Name, payload.ID)
		if err != nil {
			return err
		}
	}

	if payload.ID != "" {
		_, err := pg.FindById(payload.ID)
		if err != nil {
			return fmt.Errorf("progress with ID %s not found", payload.ID)
		}
	}
	return pg.repo.Update(payload)
}

func (pg *progressUseCase) SearchBy(by map[string]interface{}) ([]model.Progress, error) {
	progresss, err := pg.repo.Search(by)
	if err != nil {
		return nil, fmt.Errorf("request invalid")
	}
	return progresss, nil
}

func (pg *progressUseCase) Pagination(requestQueryParams dto.RequestQueryParams) ([]model.Progress, dto.Paging, error) {
	if !requestQueryParams.QueryParams.IsSortValid() {
		return nil, dto.Paging{}, fmt.Errorf("invalid sort by: %s", requestQueryParams.QueryParams.Sort)
	}
	return pg.repo.Paging(requestQueryParams)
}

func NewProgressUseCase(repo repository.ProgressRepository) ProgressUseCase {
	return &progressUseCase{repo: repo}
}
