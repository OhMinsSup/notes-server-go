package stores

import (
	"database/sql"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/OhMinsSup/notes-server-go/tools/errors"
)

func (s *Store) GetSchemaName() (string, error) {
	var query sq.SelectBuilder

	switch s.dbType {
	case MysqlDBType:
		query = s.getQueryBuilder(s.db).Select("DATABASE()")
	case PostgresDBType:
		query = s.getQueryBuilder(s.db).Select("current_schema()")
	case SqliteDBType:
		return "", nil
	default:
		return "", errors.ErrUnsupportedDatabaseType
	}

	scanner := query.QueryRow()

	var result string
	err := scanner.Scan(&result)
	if err != nil && !errors.IsErrNotFound(err) {
		return "", err
	}
	return result, nil
}

func (s *Store) CloseRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		log.Println("error closing row set", err)
	}
}
