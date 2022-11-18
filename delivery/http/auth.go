package http

import (
	"net/http"

	"github.com/shkiperko0/auth-go-ms/common"
	"github.com/shkiperko0/auth-go-ms/usecases"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHTTPHandler struct {
	AuthUC usecases.AuthUseCase
	UserUC usecases.UserUseCase
}

type UserRegisterModel struct {
	NickName  string
	Email     string
	Password  string
	PromoCode *string
	ReffID    *string
}

type UserLoginModel struct {
	NickName *string
	Email    *string
	Password string
}

func (handler *AuthHTTPHandler) Login(c echo.Context) error {
	var model UserLoginModel
	err := c.Bind(&model)
	if err != nil {
		return err
	}

	data, err := handler.AuthUC.Login()
	if err != nil {
		return status.Errorf(codes.NotFound, err.Error())
	}

	return c.JSON(http.StatusOK, data)
}

// func (handler *AuthHTTPHandler) Register(c echo.Context) error {
// 	var model UserRegisterModel
// 	err := c.Bind(&model)
// 	if err != nil {
// 		return err
// 	}

// 	user, err := h.userUC.Register(model.Type, &model)
// 	if err != nil {
// 		return err
// 	}

// 	res := &pb.RegisterResp{}

// 	if user.Extra != nil {
// 		extraString, err := user.Extra.Value()

// 		if err != nil {
// 			return err
// 		}

// 		err = json.Unmarshal([]byte((extraString).(string)), &res)

// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return c.JSON(http.StatusOK, res)
// }

func newAuthHTTPHandler(e *echo.Echo, AuthUC usecases.AuthUseCase, UserUC usecases.UserUseCase) {
	handler := AuthHTTPHandler{
		AuthUC: AuthUC,
		UserUC: UserUC,
	}

	e.POST(common.API_VER_1+"/auth/register", handler.Register)
	e.POST(common.API_VER_1+"/auth/login", handler.Login)
}
