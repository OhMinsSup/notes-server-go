package models

import "database/sql"

type UserSocial struct {
	BaseModelLogicalDel `xorm:"extends"`
	UserId              int64          `json:"user_id" xorm:"index"`
	WebSite             sql.NullString `json:"website" xorm:"varchar(100) default null"`
}

func (UserSocial) TableName() string {
	return "user_socials"
}
