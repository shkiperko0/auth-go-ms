package interactor

import (
	"time"

	"github.com/shkiperko0/auth-go-ms/models"
	"github.com/shkiperko0/auth-go-ms/repositories"
)

type ITokenInteractor interface {
	Create(token *models.Token, duration time.Duration) error
	Update(token *models.Token, duration time.Duration) error
	UpdateAccessToken(userId uint, accessToken string, duration time.Duration) error
	GetById(id string, userId uint) (*models.Token, error)
	Remove(token *models.Token) error
	RemoveAllOfUser(token *models.Token) error
}

type TokenInteractor struct {
	TokenRepo repositories.ITokenRepository
	UserRepo  repositories.IUserRepository
}

func NewTokenInteractor(TokenRepo repositories.ITokenRepository, UserRepo repositories.IUserRepository) ITokenInteractor {
	return &TokenInteractor{
		TokenRepo: TokenRepo,
		UserRepo:  UserRepo,
	}
}

func (a *TokenInteractor) Remove(token *models.Token) error {
	if token == nil {
		return nil
	}

	return a.TokenRepo.Remove(token)
}

func (a *TokenInteractor) RemoveAllOfUser(token *models.Token) error {
	return a.TokenRepo.RemoveAllOfUser(token)
}

func (a *TokenInteractor) Create(token *models.Token, duration time.Duration) error {
	return a.TokenRepo.Create(token, duration)
}

func (a *TokenInteractor) Update(item *models.Token, duration time.Duration) error {
	return a.TokenRepo.Update(item, duration)
}

func (a *TokenInteractor) UpdateAccessToken(userId uint, accessToken string, duration time.Duration) error {
	return a.TokenRepo.UpdateAccessToken(userId, accessToken, duration)
}

func (a *TokenInteractor) GetById(id string, userId uint) (*models.Token, error) {
	return a.TokenRepo.Get(id, userId)
}
