package apis

import (
	"github.com/OhMinsSup/notes-server-go/tools/security"
	"github.com/OhMinsSup/notes-server-go/tools/tokens"
	"github.com/labstack/echo/v5"
	"github.com/spf13/cast"
)

const (
	ContextAdminKey      string = "admin"
	ContextAuthRecordKey string = "authRecord"
	ContextCollectionKey string = "collection"
)

func LoadAuthContext(app App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(app.Config().TokenName)
			if err != nil {
				return next(c)
			}

			token := cookie.Value

			claims, _ := security.ParseUnverifiedJWT(token)
			tokenType := cast.ToString(claims["type"])

			switch tokenType {
			case tokens.TypeAdmin:
			case tokens.TypeAuthRecord:
				user, err := app.Store().FindAuthRecordByToken(token, app.Settings().RecordAuthToken.Secret)
				if err == nil && user != nil {
					c.Set(ContextAuthRecordKey, user)
				}
			}
			return next(c)
		}
	}
}
