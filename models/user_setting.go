package models

type UserSetting struct {
	BaseModelLogicalDel `xorm:"extends"`
	UserId              int64 `json:"user_id" xorm:"index"`
	PrivateTweet        bool  `json:"private_tweet" xorm:"notnull default false"`
}

func (UserSetting) TableName() string {
	return "user_settings"
}
