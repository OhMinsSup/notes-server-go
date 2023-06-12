package apis

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

// bindHealthApi registers the health api endpoint.
func bindHealthApi(app App, rg *echo.Group) {
	api := healthApi{app: app}

	subGroup := rg.Group("/health")
	subGroup.GET("", api.healthCheck)
}

type healthApi struct {
	app App
}

type healthCheckResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		CanBackup bool `json:"canBackup"`
	} `json:"data"`
}

// healthCheck returns a 200 OK response if the server is healthy.
func (api *healthApi) healthCheck(c echo.Context) error {
	resp := new(healthCheckResponse)
	resp.Code = http.StatusOK
	resp.Message = "API is healthy."

	return c.JSON(http.StatusOK, resp)
}
