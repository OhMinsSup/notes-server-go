package stores

import "xorm.io/xorm"

func New(db *xorm.Engine) *Store {
	return &Store{
		db: db,
	}
}

type Store struct {
	db *xorm.Engine
}

func (store *Store) DataDB() *xorm.Engine {
	return store.db
}

func (store *Store) Clone() *Store {
	clone := *store

	return &clone
}
