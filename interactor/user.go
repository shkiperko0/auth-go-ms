package interactor

import (
	"github.com/shkiperko0/auth-go-ms/models"
	"github.com/shkiperko0/auth-go-ms/repositories"
)

type IUserInteractor interface {
	GetById(id uint) (*models.User, error)
}

type UserInteractor struct {
	UserRepo repositories.IUserRepository
}

func (a *UserInteractor) GetById(id uint) (*models.User, error) {
	return a.UserRepo.GetById(id)
}

func NewUserInteractor(UserRepo repositories.IUserRepository) *UserInteractor {
	return &UserInteractor{
		UserRepo: UserRepo,
	}
}
