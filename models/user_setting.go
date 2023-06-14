package models

import (
	"encoding/json"
	"io"
)

type UserSetting struct {
	BaseModelLogicalDel
	UserId       int64 `json:"user_id"`
	PrivateTweet bool  `json:"private_tweet"`
}

func (UserSetting) TableName() string {
	return "user_settings"
}

func UserSettingFormJSON(data io.Reader) *UserSetting {
	var item *UserSetting
	_ = json.NewDecoder(data).Decode(&item)
	return item
}
