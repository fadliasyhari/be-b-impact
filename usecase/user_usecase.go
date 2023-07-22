package usecase

import (
	"fmt"
	"mime/multipart"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"
	"be-b-impact.com/csr/repository"
	"be-b-impact.com/csr/utils"
)

type UsersUseCase interface {
	BaseUseCase[model.User]
	BaseUseCasePaging[model.User]
	UpdateUser(payload *model.User, image multipart.File) error
}

type usersUseCase struct {
	repo         repository.UsersRepository
	userDetailUC UserDetailUseCase
}

// UpdateData implements UsersUseCase.
func (*usersUseCase) UpdateData(payload *model.User) error {
	panic("unimplemented")
}

func (us *usersUseCase) DeleteData(id string) error {
	users, err := us.FindById(id)
	if err != nil {
		return fmt.Errorf("users with ID %s not found", id)
	}
	return us.repo.Delete(users.ID)
}

func (us *usersUseCase) FindAll() ([]model.User, error) {
	return us.repo.List()
}

func (us *usersUseCase) FindById(id string) (*model.User, error) {
	users, err := us.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("users with ID %s not found", id)
	}
	return users, nil
}

func (us *usersUseCase) SaveData(payload *model.User) error {
	err := payload.Validate()
	if err != nil {
		return err
	}
	// cek jika data sudah ada -> count > 0
	err = us.repo.CountData(payload.Email, payload.ID)
	if err != nil {
		return err
	}

	if payload.Password != "" {
		password, err := utils.HashPassword(payload.Password)
		if err != nil {
			return err
		}
		payload.Password = password
	}

	if payload.ID != "" {
		_, err := us.FindById(payload.ID)
		if err != nil {
			return fmt.Errorf("users with ID %s not found", payload.ID)
		}
	}
	return us.repo.Save(payload)
}

func (us *usersUseCase) UpdateUser(payload *model.User, image multipart.File) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }
	// cek jika data sudah ada -> count > 0

	if payload.Email != "" {
		err := us.repo.CountData(payload.Email, payload.ID)
		if err != nil {
			return err
		}
	}

	if payload.Password != "" {
		password, err := utils.HashPassword(payload.Password)
		if err != nil {
			return err
		}
		payload.Password = password
	}

	if payload.ID != "" {
		_, err := us.FindById(payload.ID)
		if err != nil {
			return fmt.Errorf("users with ID %s not found", payload.ID)
		}
	}

	tx := us.repo.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := us.repo.Update(payload); err != nil {
		tx.Rollback()
		return err
	}
	if image != nil {
		imageURL, err := us.userDetailUC.FirebaseUpload(image)
		if err != nil {
			tx.Rollback()
			return err
		}

		imagePayload := model.UserDetail{
			UserID:   payload.ID,
			ImageURL: imageURL,
		}

		if err := us.userDetailUC.SaveUserDetail(&imagePayload, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (us *usersUseCase) SearchBy(by map[string]interface{}) ([]model.User, error) {
	users, err := us.repo.Search(by)
	if err != nil {
		return nil, fmt.Errorf("request invalid")
	}
	return users, nil
}

func (us *usersUseCase) Pagination(requestQueryParams dto.RequestQueryParams) ([]model.User, dto.Paging, error) {
	if !requestQueryParams.QueryParams.IsSortValid() {
		return nil, dto.Paging{}, fmt.Errorf("invalid sort by: %s", requestQueryParams.QueryParams.Sort)
	}
	return us.repo.Paging(requestQueryParams)
}

func NewUsersUseCase(repo repository.UsersRepository, userDetailUC UserDetailUseCase) UsersUseCase {
	return &usersUseCase{
		repo:         repo,
		userDetailUC: userDetailUC,
	}
}
