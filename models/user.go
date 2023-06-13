package models

type User struct {
	BaseModelLogicalDel `xorm:"extends"`
	Email               string `json:"email" xorm:"varchar(100) notnull unique"`
	Username            string `json:"username" xorm:"varchar(30) notnull unique"`
	Password            string `json:"password" xorm:"varchar(100) notnull"`
	Salt                string `json:"salt" xorm:"varchar(100) notnull"`
}

func (User) TableName() string {
	return "users"
}

func (selft *User) GetUserByUsername(username string) (bool, error) {
	return true, nil
}
