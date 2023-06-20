package tokens

import (
	"github.com/OhMinsSup/notes-server-go/tools/security"
	"github.com/golang-jwt/jwt/v4"
)

// NewRecordAuthToken generates and returns a new auth record authentication token.
func NewRecordAuthToken(userId int, secret string, duration int64) (string, error) {
	return security.NewToken(
		jwt.MapClaims{
			"id":   userId,
			"type": TypeAuthRecord,
		},
		secret,
		duration,
	)
}
