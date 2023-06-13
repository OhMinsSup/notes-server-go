package models

import "time"

type BaseModelLogicalDel struct {
	Id        int64
	CreatedAt time.Time  `xorm:"created"`
	UpdatedAt time.Time  `xorm:"updated"`
	DeletedAt *time.Time `xorm:"deleted"`
}
