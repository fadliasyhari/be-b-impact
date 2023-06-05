package usecase

import (
	"fmt"
	"mime/multipart"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"
	"be-b-impact.com/csr/repository"
)

type ImageUseCase interface {
	BaseUseCase[model.Image]
	BaseUseCasePaging[model.Image]
	FirebaseUpload(file multipart.File) (string, error)
}

type imageUseCase struct {
	repo repository.ImageRepository
}

func (im *imageUseCase) DeleteData(id string) error {
	image, err := im.FindById(id)
	if err != nil {
		return fmt.Errorf("image with ID %s not found", id)
	}
	return im.repo.Delete(image.ID)
}

func (im *imageUseCase) FindAll() ([]model.Image, error) {
	return im.repo.List()
}

func (im *imageUseCase) FindById(id string) (*model.Image, error) {
	image, err := im.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("image with ID %s not found", id)
	}
	return image, nil
}

func (im *imageUseCase) FirebaseUpload(file multipart.File) (string, error) {
	return im.repo.FirebaseSave(file)
}

func (im *imageUseCase) SaveData(payload *model.Image) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }
	return im.repo.Save(payload)
}

func (im *imageUseCase) UpdateData(payload *model.Image) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }
	// cek jika data sudah ada -> count > 0

	if payload.ID != "" {
		_, err := im.FindById(payload.ID)
		if err != nil {
			return fmt.Errorf("image with ID %s not found", payload.ID)
		}
	}
	return im.repo.Update(payload)
}

func (im *imageUseCase) SearchBy(by map[string]interface{}) ([]model.Image, error) {
	images, err := im.repo.Search(by)
	if err != nil {
		return nil, fmt.Errorf("request invalid")
	}
	return images, nil
}

func (im *imageUseCase) Pagination(requestQueryParams dto.RequestQueryParams) ([]model.Image, dto.Paging, error) {
	if !requestQueryParams.QueryParams.IsSortValid() {
		return nil, dto.Paging{}, fmt.Errorf("invalid sort by: %s", requestQueryParams.QueryParams.Sort)
	}
	return im.repo.Paging(requestQueryParams)
}

func NewImageUseCase(repo repository.ImageRepository) ImageUseCase {
	return &imageUseCase{repo: repo}
}
