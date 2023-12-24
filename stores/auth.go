package stores

import (
	sqlstore "github.com/OhMinsSup/notes-server-go"
	"github.com/OhMinsSup/notes-server-go/tools/security"
	"github.com/spf13/cast"
)

// Returns an error if the JWT token is invalid, expired or not associated to an auth collection record.
func (s *Store) FindAuthRecordByToken(token string, baseTokenKey string) (*sqlstore.UserModel, error) {
	unverifiedClaims, err := security.ParseUnverifiedJWT(token)
	if err != nil {
		return nil, err
	}

	// check required claims
	var id = unverifiedClaims["id"]

	// idx := cast.ToInt(id)
	idx := cast.ToString(id)

	prisma := s.db

	// find auth record
	record, err := prisma.User.FindUnique(
		sqlstore.User.ID.Equals(idx),
	).Exec(s.ctx)

	if err != nil {
		return nil, err
	}

	// verify token signature
	if _, err := security.ParseJWT(token, baseTokenKey); err != nil {
		return nil, err
	}

	return record, nil
}
