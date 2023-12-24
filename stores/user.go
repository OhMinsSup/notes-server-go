package stores

import (
	"time"

	sqlstore "github.com/OhMinsSup/notes-server-go"
)

func (s *Store) GetUserByUsername(username string) (*sqlstore.UserProfileModel, error) {
	prisma := s.db

	profile, err := prisma.UserProfile.FindUnique(
		sqlstore.UserProfile.Username.Equals(username),
	).Exec(s.ctx)

	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *Store) GetUserByEmail(email string) (*sqlstore.UserModel, error) {
	prisma := s.db

	user, err := prisma.User.FindUnique(
		sqlstore.User.Email.Equals(email),
	).With(
		sqlstore.User.UserPassword.Fetch(),
	).Exec(s.ctx)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) GetUserById(userId string) (*sqlstore.UserModel, error) {
	prisma := s.db

	userPassword, err := prisma.User.FindUnique(
		sqlstore.User.ID.Equals(userId),
	).With(
		sqlstore.User.UserPassword.Fetch(),
	).Exec(s.ctx)

	if err != nil {
		return nil, err
	}

	return userPassword, nil
}

type CreateUserParams struct {
	Email        string
	Username     string
	Nickname     string
	PasswordHash string
	Salt         string
}

func (s *Store) CreateUser(data *CreateUserParams) (*sqlstore.UserModel, error) {
	prisma := s.db

	user, err := prisma.User.CreateOne(
		sqlstore.User.Email.Set(data.Email),
	).Exec(s.ctx)

	if err != nil {
		return nil, err
	}

	userId := user.ID

	if _, err := prisma.UserProfile.CreateOne(
		sqlstore.UserProfile.Username.Set(data.Username),
		sqlstore.UserProfile.Nickname.Set(data.Nickname),
		sqlstore.UserProfile.User.Link(
			sqlstore.User.ID.Equals(userId),
		),
	).Exec(s.ctx); err != nil {
		return nil, err
	}

	if _, err := prisma.UserSocial.CreateOne(
		sqlstore.UserSocial.User.Link(
			sqlstore.User.ID.Equals(userId),
		),
	).Exec(s.ctx); err != nil {
		return nil, err
	}

	if _, err := prisma.UserPassword.CreateOne(
		sqlstore.UserPassword.PasswordHash.Set(data.PasswordHash),
		sqlstore.UserPassword.Salt.Set(data.Salt),
		sqlstore.UserPassword.User.Link(
			sqlstore.User.ID.Equals(userId),
		),
	).Exec(s.ctx); err != nil {
		return nil, err
	}

	if _, err := prisma.Notification.CreateOne(
		sqlstore.Notification.Type.Set("WELCOME"),
		sqlstore.Notification.Message.Set(data.Username+"님, 환영합니다!"),
		sqlstore.Notification.User.Link(
			sqlstore.User.ID.Equals(userId),
		),
	).Exec(s.ctx); err != nil {
		return nil, err
	}

	now := time.Now()
	if _, err := prisma.History.CreateOne(
		sqlstore.History.Text.Set("Joined Hashnode"),
		sqlstore.History.ItemType.Set("JOIN_HASHNODE"),
		sqlstore.History.DateAddedAt.Set(now),
		sqlstore.History.User.Link(
			sqlstore.User.ID.Equals(userId),
		),
	).Exec(s.ctx); err != nil {
		return nil, err
	}

	return user, nil
}
