package stores

import (
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

	if err != nil {
		return nil, err
	}

	return user, nil
}

type CreateUserParams struct {
	Email        string
	Username     string
	PasswordHash string
}

func (s *Store) CreateUser(data *CreateUserParams) (*sqlstore.UserModel, error) {
	prisma := s.db

	user, err := prisma.User.CreateOne(
		sqlstore.User.Email.Set(data.Email),
		sqlstore.User.Username.Set(data.Username),
		sqlstore.User.PasswordHash.Set(data.PasswordHash),
	).Exec(s.ctx)
	if err != nil {
		return nil, err
	}

	userId := user.ID

	profile := prisma.UserProfile.CreateOne(
		sqlstore.UserProfile.User.Link(
			sqlstore.User.ID.Equals(userId),
		),
	).Tx()

	setting := prisma.UserSetting.CreateOne(
		sqlstore.UserSetting.User.Link(
			sqlstore.User.ID.Equals(userId),
		),
	).Tx()

	account := prisma.UserAccount.CreateOne(
		sqlstore.UserAccount.User.Link(
			sqlstore.User.ID.Equals(userId),
		),
	).Tx()

	social := prisma.UserSocial.CreateOne(
		sqlstore.UserSocial.User.Link(
			sqlstore.User.ID.Equals(userId),
		),
	).Tx()

	if err := prisma.Prisma.Transaction(
		profile,
		setting,
		account,
		social,
	).Exec(s.ctx); err != nil {
		return nil, err
	}

	return user, nil
}
