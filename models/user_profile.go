package models

import "database/sql"

type UserProfile struct {
	BaseModelLogicalDel `xorm:"extends"`
	UserId              int64          `xorm:"index"`
	Bio                 sql.NullString `xorm:"varchar(100) null"`
	Location            sql.NullString `xorm:"varchar(30) null"`
	Birthday            sql.NullTime   `xorm:"datetime null"`
	Country             sql.NullString `xorm:"varchar(30) null"`
	Gender              sql.NullString `xorm:"varchar(10) null"`
}

func (UserProfile) TableName() string {
	return "user_profiles"
}
