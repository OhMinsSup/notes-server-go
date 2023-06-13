package services

import (
	"github.com/OhMinsSup/notes-server-go/stores"
)

type ServiceConfig struct {
	store *stores.Store
}

func New(store *stores.Store) *Service {
	return &Service{
		store: store,
	}
}

type Service struct {
	store *stores.Store
}
