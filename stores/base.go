package stores

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"

	sq "github.com/Masterminds/squirrel"
)

type Params struct {
	DBType           string
	ConnectionString string
	TablePrefix      string
	DB               *sql.DB
}

func New(params Params) (*Store, error) {
	store := &Store{
		db:               params.DB,
		dbType:           params.DBType,
		tablePrefix:      params.TablePrefix,
		connectionString: params.ConnectionString,
	}

	var err error
	store.isBinaryParam, err = store.computeBinaryParam()
	if err != nil {
		return nil, err
	}

	store.schemaName, err = store.GetSchemaName()
	if err != nil {
		return nil, err
	}

	// TODO: Migration Code

	return store, nil
}

type Store struct {
	db               *sql.DB
	dbType           string
	tablePrefix      string
	connectionString string
	isBinaryParam    bool
	schemaName       string
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) DataDB() *sql.DB {
	return s.db
}

func (s *Store) DBType() string {
	return s.dbType
}

func (s *Store) getQueryBuilder(db sq.BaseRunner) sq.StatementBuilderType {
	builder := sq.StatementBuilder
	if s.dbType == PostgresDBType || s.dbType == SqliteDBType {
		builder = builder.PlaceholderFormat(sq.Dollar)
	}
	return builder.RunWith(db)
}

func (s *Store) computeBinaryParam() (bool, error) {
	if s.dbType != PostgresDBType {
		return false, nil
	}
	url, err := url.Parse(s.connectionString)
	if err != nil {
		return false, err
	}
	return url.Query().Get("binary_parameters") == "yes", nil
}

func (s *Store) escapeField(fieldName string) string { //nolint:unparam
	if s.dbType == MysqlDBType {
		return "`" + fieldName + "`"
	}
	if s.dbType == PostgresDBType || s.dbType == SqliteDBType {
		return "\"" + fieldName + "\""
	}
	return fieldName
}

func (s *Store) concatenationSelector(field string, delimiter string) string {
	if s.dbType == SqliteDBType {
		return fmt.Sprintf("group_concat(%s)", field)
	}
	if s.dbType == PostgresDBType {
		return fmt.Sprintf("string_agg(%s, '%s')", field, delimiter)
	}
	if s.dbType == MysqlDBType {
		return fmt.Sprintf("GROUP_CONCAT(%s SEPARATOR '%s')", field, delimiter)
	}
	return ""
}

func (s *Store) elementInColumn(column string) string {
	if s.dbType == SqliteDBType || s.dbType == MysqlDBType {
		return fmt.Sprintf("instr(%s, ?) > 0", column)
	}
	if s.dbType == PostgresDBType {
		return fmt.Sprintf("position(? in %s) > 0", column)
	}
	return ""
}

func (s *Store) DBVersion() string {
	var version string
	var row *sql.Row

	switch s.dbType {
	case MysqlDBType:
		row = s.db.QueryRow("SELECT VERSION()")
	case PostgresDBType:
		row = s.db.QueryRow("SHOW server_version")
	case SqliteDBType:
		row = s.db.QueryRow("SELECT sqlite_version()")
	default:
		return ""
	}

	if err := row.Scan(&version); err != nil {
		log.Printf("Unable to get database version: %v", err)
		return ""
	}

	return version
}
