package stores

import (
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/OhMinsSup/notes-server-go/models"
)

func (s *Store) GetUserByEmail(email string) (*models.User, error) {
	return s.getUserByCondition(s.db, sq.Eq{"email": email})
}

func (s *Store) GetUserByUsername(username string) (*models.User, error) {
	return s.getUserByCondition(s.db, sq.Eq{"username": username})
}

func (s *Store) CreateUser(user *models.User) (*models.User, error) {
	return s.createUser(s.db, user)
}

func (s *Store) getUserByCondition(db sq.BaseRunner, condition sq.Eq) (*models.User, error) {
	users, err := s.getUsersByCondition(db, condition, 0)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, sql.ErrNoRows
	}

	return users[0], nil
}

func (s *Store) getUsersByCondition(db sq.BaseRunner, condition interface{}, limit uint64) ([]*models.User, error) {
	query := s.getQueryBuilder(db).
		Select(
			"id",
			"email",
			"username",
			"password",
			"salt",
		).
		From(s.tablePrefix + "users").
		Where(sq.Eq{"deleted_at": 0}).
		Where(condition)

	if limit != 0 {
		query = query.Limit(limit)
	}

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer s.CloseRows(rows)

	users, err := s.usersFromRows(rows)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, sql.ErrNoRows
	}

	return users, nil

}

func (s *Store) createUser(db sq.BaseRunner, user *models.User) (*models.User, error) {
	now := time.Now().UnixNano()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.DeletedAt = 0

	query := s.getQueryBuilder(db).Insert(s.tablePrefix+"users").
		Columns("id", "username", "email", "password", "salt", "create_at", "update_at", "delete_at").
		Values(user.Id, user.Username, user.Email, user.Password, user.Salt, user.CreatedAt, user.UpdatedAt, user.DeletedAt)

	_, err := query.Exec()
	return user, err
}

func (s *Store) usersFromRows(rows *sql.Rows) ([]*models.User, error) {
	users := []*models.User{}

	for rows.Next() {
		var user models.User

		err := rows.Scan(
			&user.Id,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.Salt,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}
