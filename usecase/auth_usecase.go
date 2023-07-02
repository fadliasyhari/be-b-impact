package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/repository"
	"be-b-impact.com/csr/utils/authenticator"
	"github.com/go-redis/redis/v8"
)

type AuthUseCase interface {
	Login(email string, password string) (currentUser *model.User, token string, err error)
	Logout(token string) error
	TokenRegister(user *model.User) (token string, err error)
	StoreOTP(ctx context.Context, email, otp string) error
	GetOTP(ctx context.Context, email string) (string, error)
	DeleteOTP(ctx context.Context, email string) error
}

type authUseCase struct {
	tokenService authenticator.AccessToken
	repo         repository.UsersRepository
	redisClient  *redis.Client
}

func (a *authUseCase) DeleteOTP(ctx context.Context, email string) error {
	err := a.redisClient.Del(ctx, email).Err()
	if err != nil {
		return err
	}
	return nil
}

func (a *authUseCase) StoreOTP(ctx context.Context, email, otp string) error {
	err := a.redisClient.Set(ctx, email, otp, time.Minute*5).Err()
	if err != nil {
		return err
	}
	return nil
}

func (a *authUseCase) GetOTP(ctx context.Context, email string) (string, error) {
	otp, err := a.redisClient.Get(ctx, email).Result()
	if err != nil {
		if err == redis.Nil {
			return "", errors.New("OTP not found")
		}
		return "", err
	}
	return otp, nil
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

func (a *authUseCase) Login(email string, password string) (currentUser *model.User, token string, err error) {
	user, err := a.repo.GetByEmailPassword(email, password)
	if err != nil {
		return nil, "", fmt.Errorf("user with email %s not found", email)
	}

	tokenDetail, _ := a.tokenService.CreateAccessToken(user)
	err = a.tokenService.StoreAccessToken(user.Email, tokenDetail)
	if err != nil {
		return nil, "", err
	}
	return user, tokenDetail.AccessToken, nil
}

func (a *authUseCase) TokenRegister(user *model.User) (token string, err error) {
	tokenDetail, _ := a.tokenService.CreateAccessToken(user)
	err = a.tokenService.StoreAccessToken(user.Email, tokenDetail)
	if err != nil {
		return "", err
	}
	return tokenDetail.AccessToken, nil
}

func NewAuthUseCase(service authenticator.AccessToken, repo repository.UsersRepository, redisClient *redis.Client) AuthUseCase {
	return &authUseCase{
		tokenService: service,
		repo:         repo,
		redisClient:  redisClient,
	}
}
