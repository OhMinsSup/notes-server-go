package models

type User struct {
	BaseModelLogicalDel `xorm:"extends"`
	Email               string `xorm:"varchar(100) notnull unique"`
	Username            string `xorm:"varchar(30) notnull unique"`
	Password            string `xorm:"varchar(100) notnull"`
	Salt                string `xorm:"varchar(100) notnull"`
}

func (User) TableName() string {
	return "users"
}
