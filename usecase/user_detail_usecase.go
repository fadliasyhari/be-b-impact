package usecase

import (
	"fmt"
	"mime/multipart"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/repository"
	"gorm.io/gorm"
)

type UserDetailUseCase interface {
	FirebaseUpload(file multipart.File) (string, error)
	SaveUserDetail(payload *model.UserDetail, tx *gorm.DB) error
	DeleteDataTrx(id string, tx *gorm.DB) error
}

type userDetailUseCase struct {
	repo repository.UserDetailRepository
}

func (ud *userDetailUseCase) DeleteDataTrx(id string, tx *gorm.DB) error {
	userDetail, err := ud.FindById(id)
	if err != nil {
		return fmt.Errorf("userDetail with ID %s not found", id)
	}
	return ud.repo.DeleteTrx(userDetail.ID, tx)
}

func (ud *userDetailUseCase) FindById(id string) (*model.UserDetail, error) {
	userDetail, err := ud.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("userDetail with ID %s not found", id)
	}
	return userDetail, nil
}

func (ud *userDetailUseCase) FirebaseUpload(file multipart.File) (string, error) {
	return ud.repo.FirebaseSave(file)
}

func (ud *userDetailUseCase) SaveUserDetail(payload *model.UserDetail, tx *gorm.DB) error {
	// Save the UserDetail using the provided transaction
	if err := ud.repo.SaveTrx(payload, tx); err != nil {
		return err
	}
	return nil
}

func NewUserDetailUseCase(repo repository.UserDetailRepository) UserDetailUseCase {
	return &userDetailUseCase{repo: repo}
}
