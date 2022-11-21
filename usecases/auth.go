package usecases

import (
	"github.com/shkiperko0/auth-go-ms/handlers"
	"github.com/shkiperko0/auth-go-ms/interactor"
	"github.com/shkiperko0/auth-go-ms/models"
	"github.com/shkiperko0/auth-go-ms/models/dto"
)

type IAuthUseCase interface {
	Register(data *dto.UserRegisterModel) (*models.User, error)
	Login(data *dto.UserLoginModel) (*models.User, error)
}

type AuthUseCase struct {
	JwtInteractor  interactor.IJwtInteractor
	UserInteractor interactor.IUserInteractor
	Provider       handlers.IAuthProvider
}

func NewAuthUseCase(JwtInteractor interactor.IJwtInteractor, UserInteractor interactor.IUserInteractor, Provider handlers.IAuthProvider) IAuthUseCase {
	return &AuthUseCase{
		JwtInteractor:  JwtInteractor,
		UserInteractor: UserInteractor,
		Provider:       Provider,
	}
}

func (u *AuthUseCase) Register(data *dto.UserRegisterModel) (*models.User, error) {

	user, err := u.Provider.Register(data)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *AuthUseCase) Login(data *dto.UserLoginModel) (*models.User, error) {

	user, err := u.Provider.Login(data)
	if err != nil {
		return nil, err
	}

	return user, nil
}
