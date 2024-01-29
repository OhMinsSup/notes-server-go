package stores

import (
	"errors"

	sqlstore "github.com/OhMinsSup/notes-server-go"
	api_errors "github.com/OhMinsSup/notes-server-go/tools/errors"
)

// UserAccountStatus 사용자 계정 상태
type UserAccountStatus int

const (
	// UserAccountStatusDeactivated 사용자 계정 상태: 동결
	UserAccountStatusDeactivated UserAccountStatus = 0
	// UserAccountStatusActive 사용자 계정 상태: 사용
	UserAccountStatusActive UserAccountStatus = 1
	// UserAccountStatusSuspended 사용자 계정 상태: 일시중지
	UserAccountStatusSuspended UserAccountStatus = 2
)

var userAccountStatuses = map[UserAccountStatus]bool{
	UserAccountStatusDeactivated: true,
	UserAccountStatusActive:      true,
	UserAccountStatusSuspended:   true,
}

func (v UserAccountStatus) Valid() bool {
	return userAccountStatuses[v]
}

func (v UserAccountStatus) Int() int {
	return int(v)
}

// UserAccountType 사용자 계정 타입
type UserType int

const (
	UserTypeHuman UserType = iota
	UserTypeBot
	UserTypeWebhook
)

type CreateUserParams struct {
	Name         string `json:"name"`
	DisplayName  string `json:"displayName"`
	Role         string `json:"role"`
	Icon         string `json:"icon"`
	PasswordHash string `json:"passwordHash"`
	Salt         string `json:"salt"`
}

func (s *Store) CreateUser(data *CreateUserParams) (*sqlstore.UserModel, error) {
	client := s.db

	exist, err := client.User.FindUnique(
		sqlstore.User.Name.Equals(data.Name),
	).Exec(s.ctx)

	if err != nil && !errors.Is(err, sqlstore.ErrNotFound) {
		return nil, err
	}

	if exist != nil {
		return nil, api_errors.ErrorAlreadyExists
	}

	txUser := client.User.CreateOne(
		sqlstore.User.Name.Set(data.Name),
		sqlstore.User.DisplayName.Set(data.DisplayName),
		sqlstore.User.Password.Set(data.PasswordHash),
		sqlstore.User.Salt.Set(data.Salt),
		sqlstore.User.Icon.Set(data.Icon),
	).Tx()

	if err := client.Prisma.Transaction(txUser).Exec(s.ctx); err != nil {
		return nil, err
	}

	result := txUser.Result()

	txProfile := client.UserProfile.CreateOne(
		sqlstore.UserProfile.FkUserID.Set(result.ID),
	).Tx()

	if err := client.Prisma.Transaction(txProfile).Exec(s.ctx); err != nil {
		return nil, err
	}

	return result, nil

	// user, err := client.User.CreateOne(
	// 	sqlstore.User.Name.Set(data.Name),
	// 	sqlstore.User.DisplayName.Set(data.DisplayName),
	// 	sqlstore.User.Password.Set(data.PasswordHash),
	// 	sqlstore.User.Salt.Set(data.Salt),
	// 	sqlstore.User.Icon.Set(data.Icon),
	// ).Exec(s.ctx)

	// if err != nil {
	// 	return nil, err
	// }

	// if _, err := client.UserProfile.CreateOne(
	// 	sqlstore.UserProfile.FkUserID.Set(user.ID),
	// ).Exec(s.ctx); err != nil {
	// 	return nil, err
	// }

	// log.Println("profile", profile.Result())

	// if err := prisma.Prisma.Transaction(user, profile).Exec(s.ctx); err != nil {
	// 	log.Println(err)
	// 	return nil, err
	// }

	// result := user.Result()

	// return result, nil
	// return user, nil
}
