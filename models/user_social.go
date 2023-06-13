package models

import "database/sql"

type UserSocial struct {
	BaseModelLogicalDel `xorm:"extends"`
	UserId              int64          `xorm:"index"`
	WebSite             sql.NullString `xorm:"varchar(100) null"`
}

func (UserSocial) TableName() string {
	return "user_socials"
}
