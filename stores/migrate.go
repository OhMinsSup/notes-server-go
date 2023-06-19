package stores

import (
	"database/sql"
	"embed"
)

// go:embed migrations/*.sql
var Assets embed.FS

const (
	uniqueIDsMigrationRequiredVersion        = 14
)

func (s *Store) Migrate() error {
	return nil
}

func (s *Store) getMigrationConnection() (*sql.DB, error) {
	return nil, nil
}
