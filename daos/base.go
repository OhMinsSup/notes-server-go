package daos

import "xorm.io/xorm"

func New(db *xorm.Engine) *Dao {
	return &Dao{
		db: db,
	}
}

type Dao struct {
	db *xorm.Engine
}

func (dao *Dao) DataDB() *xorm.Engine {
	return dao.db
}

func (dao *Dao) Clone() *Dao {
	clone := *dao

	return &clone
}
