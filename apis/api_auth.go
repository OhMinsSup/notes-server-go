package apis

import (
	"log"
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

type signupForm struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (api *authApi) signup(c echo.Context) error {
	form := new(signupForm)
	if err := c.Bind(form); err != nil {
		return NewBadRequestError("Invalid form data.", err)
	}

	store := api.app.Store()
	user, exists := store.GetUserByEmail(form.Email)
	log.Println("exists", exists, user)

	resp := new(signinResponse)
	resp.Code = http.StatusOK
	resp.Message = "API is healthy."

	return c.JSON(http.StatusOK, resp)
}
