package apis

import (
	"github.com/OhMinsSup/notes-server-go/core"
	"github.com/labstack/echo/v4"
)

func InitApi(app core.App) (*echo.Echo, error) {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	return e, nil
}