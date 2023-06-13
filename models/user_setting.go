package models

type UserSetting struct {
	BaseModelLogicalDel `xorm:"extends"`
	UserId              int64 `xorm:"index"`
	PrivateTweet        bool  `xorm:"notnull default false"`
}

func (UserSetting) TableName() string {
	return "user_settings"
}
