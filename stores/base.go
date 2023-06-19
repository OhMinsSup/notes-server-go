package stores

import (
	"context"

	sqlstore "github.com/OhMinsSup/notes-server-go"
)

type Params struct {
	DB  *sqlstore.PrismaClient
	Ctx context.Context
}

func New(params Params) (*Store, error) {
	store := &Store{
		db:  params.DB,
		ctx: params.Ctx,
	}

	return store, nil
}

func (s *Store) Close() error {
	return s.db.Disconnect()
}

type Store struct {
	db  *sqlstore.PrismaClient
	ctx context.Context
}
