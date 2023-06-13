package models

import "time"

type BaseModelLogicalDel struct {
	Id        int64      `json:"id"`
	CreatedAt time.Time  `json:"created_at" xorm:"created"`
	UpdatedAt time.Time  `json:"update_at" xorm:"updated"`
	DeletedAt *time.Time `json:"delete_at" xorm:"deleted"`
}
