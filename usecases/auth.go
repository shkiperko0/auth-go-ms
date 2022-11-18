package usecases

import (
	"github.com/shkiperko0/auth-go-ms/iteractor"
)

type IAuthUseCase interface {
}

type AuthUseCase struct {
	JwtIteractor  iteractor.JwtIteractor
	UserIteractor iteractor.UserIteractor
}

func newAuthUseCase(JwtIteractor iteractor.JwtIteractor, UserIteractor iteractor.UserIteractor) IAuthUseCase {
	return &AuthUseCase{
		JwtIteractor:  JwtIteractor,
		UserIteractor: UserIteractor,
	}
}
