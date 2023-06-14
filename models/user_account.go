package models

import (
	"encoding/json"
	"io"
)

type UserAccount struct {
	BaseModelLogicalDel
	UserId int64 `json:"user_id"`
}

func (UserAccount) TableName() string {
	return "user_accounts"
}

func UserAccountFormJSON(data io.Reader) *UserAccount {
	var item *UserAccount
	_ = json.NewDecoder(data).Decode(&item)
	return item
}
