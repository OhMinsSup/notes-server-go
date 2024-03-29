package apis

import (
	"errors"
	"net/http"
	"time"

	"github.com/OhMinsSup/notes-server-go/stores"
	api_errors "github.com/OhMinsSup/notes-server-go/tools/errors"
	"github.com/OhMinsSup/notes-server-go/tools/security"
	"github.com/OhMinsSup/notes-server-go/tools/tokens"
	"github.com/labstack/echo/v5"
	"github.com/spf13/cast"
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
	Code       int    `json:"code"`
	ResultCode int    `json:"resultCode"`
	Message    string `json:"message"`
	Result     struct {
		UserId    string `json:"userId"`
		AuthToken string `json:"authToken"`
	} `json:"result"`
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

	// store := api.app.Store()
	// user, err := store.GetUserByEmail(form.Email)
	// if err != nil {
	// 	return NewNotFoundError("User not found.", err)
	// }

	// userPassword, _ := user.UserPassword()
	// if userPassword == nil {
	// 	return NewInternalServerError("Failed to get user password.", err)
	// }

	// if !security.ComparePassword(userPassword.PasswordHash, form.Password) {
	// 	type Data struct {
	// 		Key string `json:"key"`
	// 	}
	// 	var data Data
	// 	data.Key = "password"
	// 	return NewBadRequestError("Invalid password.", data)
	// }

	// token, tokenErr := tokens.NewRecordAuthToken(user.ID, api.app.Settings().RecordAuthToken.Secret, api.app.Settings().RecordAuthToken.Duration)
	// if tokenErr != nil {
	// 	return NewBadRequestError("Failed to auth token.", tokenErr)
	// }

	resp := new(signinResponse)
	// resp.Code = http.StatusOK
	// resp.ResultCode = http.StatusOK
	// resp.Message = "Success"
	// resp.Result.UserId = user.ID
	// resp.Result.AuthToken = token

	// cookie := new(http.Cookie)
	// cookie.Name = api.app.Config().TokenName
	// cookie.Value = token
	// cookie.Path = "/"
	// cookie.HttpOnly = true
	// // 14 days
	// cookie.Expires = time.Now().Add(14 * 24 * time.Hour)
	// c.SetCookie(cookie)

	return c.JSON(http.StatusOK, resp)
}

type signupResponse struct {
	Code       int    `json:"code"`
	ResultCode int    `json:"resultCode"`
	Message    string `json:"message"`
	Result     struct {
		UserId    string `json:"userId"`
		AuthToken string `json:"authToken"`
	} `json:"result"`
}

type signupForm struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Role        string `json:"role"`
	Icon        string `json:"icon"`
	Password    string `json:"password"`
}

func (api *authApi) signup(c echo.Context) error {
	form := new(signupForm)
	if err := c.Bind(form); err != nil {
		return NewBadRequestError("Invalid form data.", err)
	}

	store := api.app.Store()

	hash, err := security.HashPassword(form.Password)
	if err != nil {
		return NewBadRequestError("Failed to hash password.", err)
	}

	data := stores.CreateUserParams{
		Name:         form.Name,
		DisplayName:  form.DisplayName,
		Role:         form.Role,
		Icon:         form.Icon,
		PasswordHash: hash,
		Salt:         cast.ToString(security.PasswordHashStrength),
	}

	result, err := store.CreateUser(&data)
	if err != nil {
		if errors.Is(err, api_errors.ErrorAlreadyExists) {
			return NewBadRequestError("User already exists.", err)
		}
		return NewInternalServerError("Failed to create user.", err)
	}

	token, tokenErr := tokens.NewRecordAuthToken(result.ID, api.app.Settings().RecordAuthToken.Secret, api.app.Settings().RecordAuthToken.Duration)
	if tokenErr != nil {
		return NewBadRequestError("Failed to create auth token.", tokenErr)
	}

	resp := new(signupResponse)
	resp.Code = http.StatusOK
	resp.ResultCode = http.StatusOK
	resp.Message = "Success"
	resp.Result.UserId = result.ID
	resp.Result.AuthToken = token

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
