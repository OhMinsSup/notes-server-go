package models

type BaseModelLogicalDel struct {
	Id        int64 `json:"id"`
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"update_at"`
	DeletedAt int64 `json:"delete_at"`
}
