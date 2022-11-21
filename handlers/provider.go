package handlers

import (
	"errors"

	"github.com/shkiperko0/auth-go-ms/handlers/idp"
	"github.com/shkiperko0/auth-go-ms/interactor"
	"github.com/shkiperko0/auth-go-ms/models"
)

type IAuthProvider interface {
	Register(data interface{}) (*models.User, error)
	Login(data interface{}) (*models.User, error)
}

type AuthProvider struct {
	UserInteractor     interactor.IUserInteractor
	PasswordInteractor interactor.IPasswordInteractor
	ProviderRegistry   idp.HandlerProviderRegistry
}

func (e AuthProvider) Register(data interface{}) (*models.User, error) {

	user := &models.User{
		//FirstName:      &extra.UserInfo.Username,
		Role: "user",
		//Email:          extra.UserInfo.Email,
		//Password:       e.PasswordInteractor.Hash(uuid.NewV4().String(), int(extra.Organization.PasswordSaltSize)),
		Verified: true,
	}

	return user, nil
}

type SLoginData struct {
	UserID uint
}

func (e AuthProvider) Login(idata interface{}) (*models.User, error) {
	data, ok := idata.(SLoginData)
	if !ok {
		return nil, errors.New("SLoginData?")
	}

	user, err := e.UserInteractor.GetById(data.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func NewAuthProvider(UserInteractor interactor.IUserInteractor, PasswordInteractor interactor.IPasswordInteractor) IAuthProvider {
	return &AuthProvider{
		UserInteractor:     UserInteractor,
		PasswordInteractor: PasswordInteractor,
	}
}
