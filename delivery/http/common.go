package http

import (
	"net/http"

	"github.com/shkiperko0/auth-go-ms/common"

	"github.com/labstack/echo/v4"
)

type CommonHTTPHandler struct{}

func NewCommonHTTPHandler(e *echo.Echo) {
	CommonHTTPHandler := CommonHTTPHandler{}

	e.GET(common.API_VER_NO+"/health", CommonHTTPHandler.Health)
	e.GET("/", CommonHTTPHandler.Health)
}

func CORS_Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		res := c.Response()
		//req := c.Request()
		//fmt.Println(req.Method, " ", req.URL.Path, " ", req.Host, " ", c.Request().Header.Get("Origin"))

		if c.Request().Method == "OPTIONS" {
			res.Header().Set("Access-Control-Allow-Origin", "*")
			res.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,PUT,PATCH,POST,DELETE")
			res.Header().Set("Access-Control-Allow-Headers", "*")
			res.WriteHeader(http.StatusNoContent)
			return nil
		}

		return next(c)
	}
}

func (*CommonHTTPHandler) Health(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
