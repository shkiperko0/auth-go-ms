package http

import (
	"net/http"

	"github.com/shkiperko0/auth-go-ms/common"
	"github.com/shkiperko0/auth-go-ms/models/dto"
	"github.com/shkiperko0/auth-go-ms/usecases"

	"github.com/labstack/echo/v4"
)

type AuthHTTPHandler struct {
	AuthUC usecases.IAuthUseCase
	UserUC usecases.IUserUseCase
}

func (handler *AuthHTTPHandler) Login(c echo.Context) error {
	var model dto.UserLoginModel
	err := c.Bind(&model)
	if err != nil {
		return err
	}

	user, err := handler.AuthUC.Login(&model)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func (h *AuthHTTPHandler) Register(c echo.Context) error {
	var model dto.UserRegisterModel
	err := c.Bind(&model)
	if err != nil {
		return err
	}

	user, err := h.AuthUC.Register(&model)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func NewAuthHTTPHandler(e *echo.Echo, AuthUC usecases.IAuthUseCase, UserUC usecases.IUserUseCase) {
	handler := AuthHTTPHandler{
		AuthUC: AuthUC,
		UserUC: UserUC,
	}

	e.POST(common.API_VER_1+"/auth/register", handler.Register)
	e.POST(common.API_VER_1+"/auth/login", handler.Login)
}
