package models

type UserAccount struct {
	BaseModelLogicalDel `xorm:"extends"`
	UserId              int64 `xorm:"index"`
}

func (UserAccount) TableName() string {
	return "user_accounts"
}
