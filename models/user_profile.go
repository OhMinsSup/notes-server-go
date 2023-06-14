package models

import (
	"database/sql"
	"encoding/json"
	"io"
)

type UserProfile struct {
	BaseModelLogicalDel
	UserId   int64          `json:"user_id"`
	Bio      sql.NullString `json:"bio"`
	Location sql.NullString `json:"location"`
	Birthday sql.NullTime   `json:"birthday"`
	Country  sql.NullString `json:"country"`
	Gender   sql.NullString `json:"gender"`
}

func (UserProfile) TableName() string {
	return "user_profiles"
}

func UserProfileFormJSON(data io.Reader) *UserProfile {
	var item *UserProfile
	_ = json.NewDecoder(data).Decode(&item)
	return item
}
