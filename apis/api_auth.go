package apis

import (
	"net/http"
	"time"

	"github.com/OhMinsSup/notes-server-go/stores"
	"github.com/OhMinsSup/notes-server-go/tools/security"
	"github.com/OhMinsSup/notes-server-go/tools/tokens"
	"github.com/labstack/echo/v5"
)

func bindAuthApi(app App, rg *echo.Group) {
	api := authApi{app: app}

	subGroup := rg.Group("/auth")
	subGroup.POST("/signup", api.signup)
	subGroup.POST("/signin", api.signin)
}

type authApi struct {
	app App
}

type signinResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		UserId int    `json:"userId"`
		Token  string `json:"token"`
	} `json:"data"`
}

type signinForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (api *authApi) signin(c echo.Context) error {
	form := new(signinForm)
	if err := c.Bind(form); err != nil {
		return NewBadRequestError("Invalid form data.", err)
	}

	store := api.app.Store()
	user, err := store.GetUserByEmail(form.Email)
	if err != nil {
		return NewNotFoundError("User not found.", err)
	}

	if !security.ComparePassword(user.PasswordHash, form.Password) {
		type Data struct {
			Key string `json:"key"`
		}
		var data Data
		data.Key = "password"
		return NewBadRequestError("Invalid password.", data)
	}

	token, tokenErr := tokens.NewRecordAuthToken(user.ID, api.app.Settings().RecordAuthToken.Secret, api.app.Settings().RecordAuthToken.Duration)
	if tokenErr != nil {
		return NewBadRequestError("Failed to auth token.", tokenErr)
	}

	resp := new(signinResponse)
	resp.Code = http.StatusOK
	resp.Message = "Success"
	resp.Data.UserId = user.ID
	resp.Data.Token = token

	cookie := new(http.Cookie)
	cookie.Name = api.app.Config().TokenName
	cookie.Value = token
	cookie.Path = "/"
	cookie.HttpOnly = true
	// 14 days
	cookie.Expires = time.Now().Add(14 * 24 * time.Hour)
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, resp)
}

type signupResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		UserId int    `json:"userId"`
		Token  string `json:"token"`
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
	user, _ := store.GetUserByEmailOrUsername(form.Email, form.Username)
	if user != nil {
		type Data struct {
			Exists bool   `json:"exists"`
			Key    string `json:"key"`
		}

		if user.Email == form.Email {
			var data Data
			data.Exists = true
			data.Key = "email"
			return NewBadRequestError(
				"Email already exists.",
				data,
			)
		}
		if user.Username == form.Username {
			var data Data
			data.Exists = true
			data.Key = "username"
			return NewBadRequestError(
				"Username already exists.",
				data,
			)
		}
		var data Data
		data.Exists = true
		data.Key = "user"
		return NewBadRequestError(
			"User already exists.",
			data,
		)
	}

	hash, err := security.HashPassword(form.Password)
	if err != nil {
		return NewInternalServerError("Failed to hash password.", err)
	}

	data := stores.CreateUserParams{
		Email:        form.Email,
		Username:     form.Username,
		PasswordHash: hash,
	}

	user, err = store.CreateUser(&data)
	if err != nil {
		return NewInternalServerError("Failed to create user.", err)
	}

	token, tokenErr := tokens.NewRecordAuthToken(user.ID, api.app.Settings().RecordAuthToken.Secret, api.app.Settings().RecordAuthToken.Duration)
	if tokenErr != nil {
		return NewBadRequestError("Failed to create auth token.", tokenErr)
	}

	resp := new(signupResponse)
	resp.Code = http.StatusOK
	resp.Message = "Success"
	resp.Data.UserId = user.ID
	resp.Data.Token = token

	cookie := new(http.Cookie)
	cookie.Name = api.app.Config().TokenName
	cookie.Value = token
	cookie.Path = "/"
	cookie.HttpOnly = true
	// 14 days
	cookie.Expires = time.Now().Add(14 * 24 * time.Hour)
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, resp)
}
