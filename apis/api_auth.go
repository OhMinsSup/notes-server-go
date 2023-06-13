package apis

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func bindAuthApi(app App, rg *echo.Group) {
	api := authApi{app: app}

	subGroup := rg.Group("/auth")
	subGroup.POST("/signup", api.signup)
}

type authApi struct {
	app App
}

type signinResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
	} `json:"data"`
}

func (api *authApi) signup(c echo.Context) error {
	resp := new(signinResponse)
	resp.Code = http.StatusOK
	resp.Message = "API is healthy."

	return c.JSON(http.StatusOK, resp)
}
