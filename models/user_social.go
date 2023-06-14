package models

import (
	"database/sql"
	"encoding/json"
	"io"
)

type UserSocial struct {
	BaseModelLogicalDel `xorm:"extends"`
	UserId              int64          `json:"user_id"`
	WebSite             sql.NullString `json:"website"`
}

func (UserSocial) TableName() string {
	return "user_socials"
}

func UserSocialFormJSON(data io.Reader) *UserSocial {
	var item *UserSocial
	_ = json.NewDecoder(data).Decode(&item)
	return item
}
