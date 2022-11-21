package usecases

import (
	"github.com/shkiperko0/auth-go-ms/interactor"
	"github.com/shkiperko0/auth-go-ms/models"
)

type UserUseCase struct {
	JwtInteractor  interactor.IJwtInteractor
	UserInteractor interactor.IUserInteractor
}

type IUserUseCase interface {
	Get(id uint) (*models.User, error)
}

func (u *UserUseCase) Get(id uint) (*models.User, error) {
	return u.UserInteractor.GetById(id)
}

func NewUserUseCase(JwtInteractor interactor.IJwtInteractor, UserInteractor interactor.IUserInteractor) IUserUseCase {
	return &UserUseCase{
		JwtInteractor:  JwtInteractor,
		UserInteractor: UserInteractor,
	}
}
