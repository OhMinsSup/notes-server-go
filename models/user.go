package models

import (
	"encoding/json"
	"io"
)

type User struct {
	BaseModelLogicalDel
	Email               string `json:"email"`
	Username            string `json:"username"`
	Password            string `json:"password"`
	Salt                string `json:"salt"`
}

func (User) TableName() string {
	return "users"
}

func UserFormJSON(data io.Reader) *User {
	var item *User
	_ = json.NewDecoder(data).Decode(&item)
	return item
}
