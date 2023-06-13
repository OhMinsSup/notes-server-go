package models

type UserAccount struct {
	BaseModelLogicalDel `xorm:"extends"`
	UserId              int64 `json:"user_id" xorm:"index"`
}

func (UserAccount) TableName() string {
	return "user_accounts"
}
