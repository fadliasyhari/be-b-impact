package usecase

import (
	"fmt"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/repository"
	"be-b-impact.com/csr/utils/authenticator"
)

type AuthUseCase interface {
	Login(username string, password string) (token string, err error)
	Logout(token string) error
	TokenRegister(user *model.User) (token string, err error)
}

type authUseCase struct {
	tokenService authenticator.AccessToken
	repo         repository.UsersRepository
}

// Logout implements AuthUseCase
func (a *authUseCase) Logout(token string) error {
	accountDetail, _ := a.tokenService.VerifyAccessToken(token)
	err := a.tokenService.DeleteAccessToken(accountDetail.AccessUUID)
	if err != nil {
		return err
	}
	return nil
}

func (a *authUseCase) Login(username string, password string) (token string, err error) {
	user, err := a.repo.GetByUsernamePassword(username, password)
	if err != nil {
		return "", fmt.Errorf("user with username %s not found", username)
	}

	tokenDetail, _ := a.tokenService.CreateAccessToken(user)
	err = a.tokenService.StoreAccessToken(user.Username, tokenDetail)
	if err != nil {
		return "", err
	}
	return tokenDetail.AccessToken, nil
}

func (a *authUseCase) TokenRegister(user *model.User) (token string, err error) {
	tokenDetail, _ := a.tokenService.CreateAccessToken(user)
	err = a.tokenService.StoreAccessToken(user.Username, tokenDetail)
	if err != nil {
		return "", err
	}
	return tokenDetail.AccessToken, nil
}

func NewAuthUseCase(service authenticator.AccessToken, repo repository.UsersRepository) AuthUseCase {
	return &authUseCase{
		tokenService: service,
		repo:         repo,
	}
}
