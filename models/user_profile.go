package models

import "database/sql"

type UserProfile struct {
	BaseModelLogicalDel `xorm:"extends"`
	UserId              int64          `json:"user_id" xorm:"index"`
	Bio                 sql.NullString `json:"bio" xorm:"varchar(100) default null"`
	Location            sql.NullString `json:"location" xorm:"varchar(30) default null"`
	Birthday            sql.NullTime   `json:"birthday" xorm:"datetime default null"`
	Country             sql.NullString `json:"country" xorm:"varchar(30) default null"`
	Gender              sql.NullString `json:"gender" xorm:"varchar(10) default null"`
}

func (UserProfile) TableName() string {
	return "user_profiles"
}
