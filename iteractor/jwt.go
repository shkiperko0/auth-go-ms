package iteractor

import (
	//"encoding/json"
	"eam-auth-go-ms/models"
	"eam-auth-go-ms/repositories"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"
)

const TK_Ver1 = "v1"
const TK_ASecret_Only_v1 = "keklol no secrets got iter guy (access man)"
const TK_RSecret_Only_v1 = "keklol no secrets got iter guy (refresh man)"

type JwtIteractor struct {
	UserRepository repositories.UserRepository
}

/*
	ToDo
	Every organization have secrets and expire time
*/

const (
	FL_ID        = "user_id"
	FL_SessionID = "session_id"
	FL_Role      = "role"
	FL_OrgSlug   = "org_slug"
	FL_Version   = "version"
)

type JwtErrors interface{}

const (
	ERR_VerifyFailed = "jwt.verify.failed"
)

func (iter *JwtIteractor) GetTokens(user_id uint) (*models.JwtTokens, error) {

	user, err := iter.UserRepository.GetById(user_id)
	if err != nil {
		return nil, err
	}

	refreshData := models.RefreshJwtToken{
		Role:      user.Role,
		ID:        user.ID,
		Version:   TK_Ver1,
		OrgSlug:   "",
		SessionID: "",
	}

	refreshToken, err := iter.SignRefresh(refreshData)
	if err != nil {
		return nil, err
	}

	accessData := models.AccessJwtToken{
		Role:      user.Role,
		ID:        user.ID,
		Version:   TK_Ver1,
		OrgSlug:   "",
		SessionID: "",
	}

	accessToken, err := iter.SignAccess(accessData)
	if err != nil {
		return nil, err
	}

	data := &models.JwtTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return data, err
}

func (iter *JwtIteractor) SignRefresh(data models.RefreshJwtToken) (string, error) {

	payload := jwt.MapClaims{
		FL_ID:        data.ID,
		FL_SessionID: data.SessionID,
		FL_OrgSlug:   data.OrgSlug,
		FL_Role:      data.Role,
		FL_Version:   TK_Ver1,
	}

	//payload.VerifyExpiresAt(org.JwtExpiresIn, true)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(TK_RSecret_Only_v1))
}

func (iter *JwtIteractor) SignAccess(data models.AccessJwtToken) (string, error) {
	payload := jwt.MapClaims{
		FL_ID:        data.ID,
		FL_SessionID: data.SessionID,
		FL_OrgSlug:   data.OrgSlug,
		FL_Role:      data.Role,
		FL_Version:   TK_Ver1,
	}

	//payload.VerifyExpiresAt(org.JwtExpiresIn, true)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(TK_ASecret_Only_v1))
}

func (iter *JwtIteractor) VerifyRefresh(tokenString string) (*models.RefreshJwtToken, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		claims := token.Claims.(jwt.MapClaims)
		version, ok := claims[FL_Version].(string)
		if ok {
			if version == TK_Ver1 {
				return []byte(TK_RSecret_Only_v1), nil
			}
		}

		return nil, errors.New(ERR_VerifyFailed)
	})

	if err != nil {
		return nil, err
	}

	if token == nil {
		return nil, errors.New(ERR_VerifyFailed)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId, ok := claims[FL_ID].(float64)
		if !ok {
			return nil, errors.New(ERR_VerifyFailed)
		}

		sessionId, ok := claims[FL_SessionID].(string)
		if !ok {
			return nil, errors.New(ERR_VerifyFailed)
		}

		orgSlug, ok := claims[FL_OrgSlug].(string)
		if !ok {
			return nil, errors.New(ERR_VerifyFailed)
		}

		role, ok := claims[FL_Role].(string)
		if !ok {
			return nil, errors.New(ERR_VerifyFailed)
		}

		data := models.RefreshJwtToken{
			ID:        uint(userId),
			SessionID: sessionId,
			OrgSlug:   orgSlug,
			Role:      role,
		}

		return &data, err
	}

	return nil, err
}

func (iter *JwtIteractor) VerifyAccess(tokenString string) (*models.AccessJwtToken, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		claims := token.Claims.(jwt.MapClaims)
		version, ok := claims[FL_Version].(string)
		if ok {
			if version == TK_Ver1 {
				return []byte(TK_RSecret_Only_v1), nil
			}
		}

		return nil, errors.New(ERR_VerifyFailed)
	})

	if err != nil {
		return nil, err
	}

	if token == nil {
		return nil, errors.New(ERR_VerifyFailed)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId, ok := claims[FL_ID].(float64)
		if !ok {
			return nil, errors.New(ERR_VerifyFailed)
		}

		sessionId, ok := claims[FL_SessionID].(string)
		if !ok {
			return nil, errors.New(ERR_VerifyFailed)
		}

		orgSlug, ok := claims[FL_OrgSlug].(string)
		if !ok {
			return nil, errors.New(ERR_VerifyFailed)
		}

		role, ok := claims[FL_Role].(string)
		if !ok {
			return nil, errors.New(ERR_VerifyFailed)
		}

		data := models.AccessJwtToken{
			ID:        uint(userId),
			SessionID: sessionId,
			OrgSlug:   orgSlug,
			Role:      role,
		}

		return &data, err
	}

	return nil, err
}
