package stores

import (
	"log"

	sqlstore "github.com/OhMinsSup/notes-server-go"
)

func (s *Store) GetUserByEmailOrUsername(email string, username string) (*sqlstore.UserModel, error) {
	prisma := s.db

	user, err := prisma.User.FindFirst(
		sqlstore.User.Or(
			sqlstore.User.Email.Equals(email),
			sqlstore.User.Username.Equals(username),
		),
	).Exec(s.ctx)

	log.Println("err: ", err)

	if err != nil {
		return nil, err
	}

	return user, nil
}
