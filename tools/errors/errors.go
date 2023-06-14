package errors

import (
	"database/sql"
	"errors"
)

var (
	ErrUnsupportedDatabaseType = errors.New("database type is unsupported")
)

func IsErrNotFound(err error) bool {
	if err == nil {
		return false
	}

	// check if this is a sql.ErrNotFound
	if errors.Is(err, sql.ErrNoRows) {
		return true
	}

	return false
}
