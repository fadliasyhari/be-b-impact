package repository

import (
	"errors"
	"fmt"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"
	"be-b-impact.com/csr/utils"

	"be-b-impact.com/csr/utils/common"
	"gorm.io/gorm"
)

type UsersRepository interface {
	BaseRepository[model.User]
	BaseRepositoryCount[model.User]
	BaseRepositoryPaging[model.User]
	GetByUsernamePassword(username string, password string) (*model.User, error)
}
type usersRepository struct {
	db *gorm.DB
}

func (us *usersRepository) Delete(id string) error {
	return us.db.Delete(&model.User{}, "id=?", id).Error
}

func (us *usersRepository) Get(id string) (*model.User, error) {
	var users model.User
	result := us.db.Select("id, username, email, role, status, created_at, updated_at").First(&users, "id=?", id).Error
	if result != nil {
		return nil, result
	}
	return &users, nil
}

func (us *usersRepository) List() ([]model.User, error) {
	var users []model.User
	result := us.db.Find(&users).Error
	if result != nil {
		return nil, result
	}
	return users, nil
}

func (us *usersRepository) Save(payload *model.User) error {
	return us.db.Save(payload).Error
}

func (us *usersRepository) Update(payload *model.User) error {
	updateFields := make(map[string]interface{})

	// Add fields to be updated based on the payload
	if payload.Email != "" {
		updateFields["email"] = payload.Email
	}

	if payload.Username != "" {
		updateFields["username"] = payload.Username
	}

	if payload.Status != "" {
		updateFields["status"] = payload.Status
	}

	if payload.Role != "" {
		updateFields["role"] = payload.Role
	}

	if payload.Password != "" {
		updateFields["password"] = payload.Password
	}

	return us.db.Model(&model.User{}).Where("id = ?", payload.ID).Updates(updateFields).Error
}

func (us *usersRepository) Search(by map[string]interface{}) ([]model.User, error) {
	var users []model.User
	query := us.db
	for key, value := range by {
		// Perform case-insensitive search using ilike
		if key == "role" && value == "admin" {
			query = query.Where(fmt.Sprintf("%s ilike ?", key), fmt.Sprintf("%%%v%%", value))
			query = query.Joins("LEFT JOIN proposals ON users.id = proposals.reviewer_id").
				Group("users.id").
				Order("COUNT(proposals.id) ASC")

		} else {
			query = query.Where(fmt.Sprintf("%s ilike ?", key), fmt.Sprintf("%%%v%%", value))
		}
	}
	result := query.Find(&users).Error
	if result != nil {
		return nil, result
	}
	return users, nil
}

func (us *usersRepository) CountData(fieldname string, id string) error {
	var count int64
	var result *gorm.DB

	if id != "" {
		result = us.db.Model(&model.User{}).Where("username ilike ? AND id <> ?", fieldname, id).Count(&count)
	} else {
		result = us.db.Model(&model.User{}).Where("username ilike ?", fieldname).Count(&count)
	}
	if result.Error != nil {
		return result.Error
	}
	if count > 0 {
		return fmt.Errorf("username %s already exist", fieldname)
	}
	return nil
}

func (us *usersRepository) Paging(requestQueryParam dto.RequestQueryParams) ([]model.User, dto.Paging, error) {
	paginationQuery, orderQuery := pagingValidate(requestQueryParam)

	var users []model.User
	query := us.db
	for key, value := range requestQueryParam.Filter {
		// Perform case-insensitive search using ilike
		query = query.Where(fmt.Sprintf("%s ilike ?", key), fmt.Sprintf("%%%v%%", value))
	}
	
	err := query.Select("id, username, email, role, status, created_at, updated_at").Order(orderQuery).Limit(paginationQuery.Take).Offset(paginationQuery.Skip).Find(&users).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}
	var totalRows int64
	err = us.db.Model(model.User{}).Count(&totalRows).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}
	return users, common.Paginate(paginationQuery.Page, paginationQuery.Take, int(totalRows)), nil
}

func (us *usersRepository) GetByUsername(username string) (*model.User, error) {
	var userCredential model.User
	result := us.db.Where("username = ?", username).First(&userCredential)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with username '%s' not found", username)
		}
		return nil, fmt.Errorf("failed to get user with username '%s': %v", username, err)
	}
	return &userCredential, nil
}

func (us *usersRepository) GetByUsernamePassword(username string, password string) (*model.User, error) {
	user, err := us.GetByUsername(username)
	if err != nil {
		return &model.User{}, err
	}

	pwdCheck := utils.CheckPasswordHash(password, user.Password)
	if !pwdCheck {
		return &model.User{}, fmt.Errorf("password don't match")
	}

	user.Password = ""
	return user, nil
}

func NewUsersRepository(db *gorm.DB) UsersRepository {
	return &usersRepository{db: db}
}
