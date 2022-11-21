package usecases

import (
	"errors"
	"strings"

	"github.com/shkiperko0/auth-go-ms/interactor"
	"github.com/shkiperko0/auth-go-ms/models/dto"
)

type ISessionUseCase interface {
	Clear(token string) error
	Check(refreshToken *string, app string, url string) (*dto.Check_isPublic, *dto.RequestError)
	CheckPublic(app string, url string) (*dto.Check_isPublic, *dto.RequestError)
}

type SessionUseCase struct {
	TokenInteractor interactor.ITokenInteractor
	JwtInteractor   interactor.IJwtInteractor
	UserInteractor  interactor.IUserInteractor
}

func (s *SessionUseCase) Clear(token string) error {
	jwtModel, err := s.JwtInteractor.VerifyRefresh(token)
	if err != nil {
		return err
	}

	sessionToken, _ := s.TokenInteractor.GetById(jwtModel.SessionID, jwtModel.ID)
	if err != nil {
		return nil
	}

	return s.TokenInteractor.Remove(sessionToken)
}

func (s *SessionUseCase) CheckPublic(app string, url string) (*dto.Check_isPublic, *dto.RequestError) {

	return &dto.Check_isPublic{
		IsPublic: true,
	}, nil
}

func (s *SessionUseCase) checkToken(app string, refreshToken string, url string) (*dto.Check_isPublic, *dto.RequestError) {
	jwtModel, err := s.JwtInteractor.VerifyRefresh(refreshToken)
	if err != nil {
		return nil, &dto.RequestError{
			StatusCode: 405,
			Err:        errors.New(err.Error()),
		}
	}

	tokenObj, err := s.TokenInteractor.GetById(jwtModel.SessionID, jwtModel.ID)
	if err != nil || tokenObj == nil {
		return nil, &dto.RequestError{
			StatusCode: 405,
			Err:        errors.New("auth-ms.session.not-found"),
		}
	}

	accessModel, err := s.JwtInteractor.VerifyAccess(tokenObj.AccessToken)
	if err != nil {
		return nil, &dto.RequestError{
			StatusCode: 405,
			Err:        errors.New("auth-ms.session.not-found"),
		}
	}

	User, err := s.UserInteractor.GetById(accessModel.ID)
	if err != nil {
		return nil, &dto.RequestError{
			StatusCode: 405,
			Err:        errors.New("auth-ms.user.not-found"),
		}
	}

	return &dto.Check_isPublic{
		IsPublic: false,
		User:     User,
	}, nil
}

const SEPARATOR = " "

func (s *SessionUseCase) Check(refreshToken *string, app string, url string) (*dto.Check_isPublic, *dto.RequestError) {
	if refreshToken != nil && len(*refreshToken) > 0 {
		if strings.ContainsAny(*refreshToken, SEPARATOR) {
			vals := strings.Split(*refreshToken, SEPARATOR)
			if len(vals) > 1 {
				if vals[0] == "Bearer" {
					return s.checkToken(app, vals[1], url)
				}
			}

			return nil, &dto.RequestError{
				StatusCode: 405,
				Err:        errors.New("auth-ms.session.not-found"),
			}
		} else if *refreshToken == "Bearer" {
			return s.CheckPublic(app, url)
		}

		return s.checkToken(app, *refreshToken, url)
	}

	return s.CheckPublic(app, url)
}

func NewSessionUseCase(TokenInteractor interactor.ITokenInteractor, JwtInteractor interactor.IJwtInteractor, UserInteractor interactor.IUserInteractor) ISessionUseCase {
	return &SessionUseCase{
		JwtInteractor:   JwtInteractor,
		TokenInteractor: TokenInteractor,
		UserInteractor:  UserInteractor,
	}
}
