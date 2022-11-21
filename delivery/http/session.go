package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shkiperko0/auth-go-ms/common"
	"github.com/shkiperko0/auth-go-ms/models/dto"
	"github.com/shkiperko0/auth-go-ms/usecases"
)

type SessionHttpHandler struct {
	SessionUC usecases.ISessionUseCase
}

func (h *SessionHttpHandler) Check(c echo.Context) error {
	var model dto.SessionCheck
	err := c.Bind(&model)
	if err != nil {
		return err
	}

	token := c.Request().Header.Get("authorization")
	app := c.Request().Header.Get("app")

	res, checkErr := h.SessionUC.Check(&token, app, model.Url)
	if checkErr != nil {
		return checkErr.Err
	}

	return c.JSON(http.StatusOK, dto.Check_isPublic{
		IsPublic: res.IsPublic,
	})
}

func (h *SessionHttpHandler) Clear(c echo.Context) error {
	var model dto.StringRequest
	err := c.Bind(&model)
	if err != nil {
		return err
	}

	h.SessionUC.Clear(model.Value)
	return c.NoContent(http.StatusOK)
}

func NewSessionHttpHandler(e *echo.Echo, SessionUC usecases.ISessionUseCase) {
	handler := SessionHttpHandler{
		SessionUC: SessionUC,
	}

	e.POST(common.API_VER_1+"/session/check", handler.Check)
	e.POST(common.API_VER_1+"/session/clear", handler.Clear)
}
