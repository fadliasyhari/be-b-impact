package usecase

import (
	"fmt"
	"mime/multipart"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"
	"be-b-impact.com/csr/repository"
	"gorm.io/gorm"
)

type FileUseCase interface {
	BaseUseCase[model.File]
	BaseUseCasePaging[model.File]
	FirebaseUpload(file multipart.File) (string, error)
	SaveTrx(payload *model.File, tx *gorm.DB) error
	DeleteDataTrx(id string, tx *gorm.DB) error
}

type fileUseCase struct {
	repo repository.FileRepository
}

func (fi *fileUseCase) DeleteData(id string) error {
	file, err := fi.FindById(id)
	if err != nil {
		return fmt.Errorf("file with ID %s not found", id)
	}
	return fi.repo.Delete(file.ID)
}

func (fi *fileUseCase) DeleteDataTrx(id string, tx *gorm.DB) error {
	file, err := fi.FindById(id)
	if err != nil {
		return fmt.Errorf("file with ID %s not found", id)
	}
	return fi.repo.DeleteTrx(file.ID, tx)
}

func (fi *fileUseCase) FindAll() ([]model.File, error) {
	return fi.repo.List()
}

func (fi *fileUseCase) FindById(id string) (*model.File, error) {
	file, err := fi.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("file with ID %s not found", id)
	}
	return file, nil
}

func (fi *fileUseCase) FirebaseUpload(file multipart.File) (string, error) {
	return fi.repo.FirebaseSave(file)
}

func (fi *fileUseCase) SaveData(payload *model.File) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }
	return fi.repo.Save(payload)
}

func (fi *fileUseCase) SaveTrx(payload *model.File, tx *gorm.DB) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }
	return fi.repo.SaveTrx(payload, tx)
}

func (fi *fileUseCase) UpdateData(payload *model.File) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }
	// cek jika data sudah ada -> count > 0

	if payload.ID != "" {
		_, err := fi.FindById(payload.ID)
		if err != nil {
			return fmt.Errorf("file with ID %s not found", payload.ID)
		}
	}
	return fi.repo.Update(payload)
}

func (fi *fileUseCase) SearchBy(by map[string]interface{}) ([]model.File, error) {
	files, err := fi.repo.Search(by)
	if err != nil {
		return nil, fmt.Errorf("request invalid")
	}
	return files, nil
}

func (fi *fileUseCase) Pagination(requestQueryParams dto.RequestQueryParams) ([]model.File, dto.Paging, error) {
	if !requestQueryParams.QueryParams.IsSortValid() {
		return nil, dto.Paging{}, fmt.Errorf("invalid sort by: %s", requestQueryParams.QueryParams.Sort)
	}
	return fi.repo.Paging(requestQueryParams)
}

func NewFileUseCase(repo repository.FileRepository) FileUseCase {
	return &fileUseCase{repo: repo}
}
