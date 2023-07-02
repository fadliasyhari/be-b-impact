package authenticator

import (
	"context"
	"errors"
	"fmt"
	"time"

	"be-b-impact.com/csr/config"
	"be-b-impact.com/csr/model"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type AccessToken interface {
	CreateAccessToken(cred *model.User) (TokenDetail, error)
	VerifyAccessToken(tokenString string) (AccessDetail, error)
	StoreAccessToken(username string, tokenDetail TokenDetail) error
	FetchAccessToken(accessDetail AccessDetail) error
	DeleteAccessToken(accessUUID string) error
}

type accessToken struct {
	Config config.TokenConfig
	client *redis.Client
}

// DeleteAccessToken implements AccessToken
func (t *accessToken) DeleteAccessToken(accessUUID string) error {
	rowAffected, err := t.client.Del(context.Background(), accessUUID).Result()
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		return errors.New("failed removing token")
	}
	return nil
}

// untuk menyimpan data ke dalam redis dan set time expired
// StoreAccessToken implements AccessToken
func (t *accessToken) StoreAccessToken(email string, tokenDetail TokenDetail) error {
	at := time.Unix(tokenDetail.AtExpired, 0)
	// membuat context operasi redis yang baru
	err := t.client.Set(
		context.Background(),
		tokenDetail.AccessUUID,
		email,
		time.Until(at),
	).Err()
	if err != nil {
		return err
	}
	return nil
}

// cek uuid dalam redis dan get data
// FetchAccessToken implements AccessToken
func (t *accessToken) FetchAccessToken(accessDetail AccessDetail) error {
	email, err := t.client.Get(
		context.Background(),
		accessDetail.AccessUUID,
	).Result()
	if err != nil {
		return err
	}

	if email == "" {
		return errors.New("invalid token")
	}

	return nil
}

func (t *accessToken) CreateAccessToken(cred *model.User) (TokenDetail, error) {

	now := time.Now().UTC()
	end := now.Add(t.Config.AccessTokenlifeTime)

	tokenDetail := TokenDetail{
		AccessUUID: uuid.New().String(),
		AtExpired:  end.Unix(),
	}
	claims := MyClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer: t.Config.ApplicationName,
		},
		UserId:     cred.ID,
		Username:   cred.Username,
		Name:       cred.Name,
		Email:      cred.Email,
		Role:       cred.Role,
		Status:     cred.Status,
		AccessUUID: tokenDetail.AccessUUID,
	}

	claims.IssuedAt = now.Unix()
	claims.ExpiresAt = end.Unix()

	// membuat accessToken dengan signin method method JwtSigninMethod (HS256)
	token := jwt.NewWithClaims(
		t.Config.JwtSigningMethod,
		claims,
	)

	newToken, err := token.SignedString([]byte(t.Config.JwtSigantureKey))
	if err != nil {
		return TokenDetail{}, err
	}
	tokenDetail.AccessToken = newToken

	return tokenDetail, nil
}

func (t *accessToken) VerifyAccessToken(tokenString string) (AccessDetail, error) {
	accessToken, err := jwt.Parse(tokenString, func(accessToken *jwt.Token) (interface{}, error) {
		if method, ok := accessToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method invalid")
		} else if method != t.Config.JwtSigningMethod {
			return nil, fmt.Errorf("signing method invalid")
		}
		return []byte(t.Config.JwtSigantureKey), nil
	})

	claims, ok := accessToken.Claims.(jwt.MapClaims)
	if !ok || !accessToken.Valid || claims["iss"] != t.Config.ApplicationName {
		return AccessDetail{}, err
	}
	username := claims["Username"].(string)
	name := claims["Name"].(string)
	email := claims["Email"].(string)
	uuid := claims["AccessUUID"].(string)
	role := claims["Role"].(string)
	status := claims["Status"].(string)
	userId := claims["UserId"].(string)
	return AccessDetail{
		AccessUUID: uuid,
		UserId:     userId,
		Name:       name,
		Email:      email,
		Username:   username,
		Role:       role,
		Status:     status,
	}, nil
}

func NewTokenService(config config.TokenConfig, client *redis.Client) AccessToken {
	return &accessToken{
		Config: config,
		client: client,
	}
}
