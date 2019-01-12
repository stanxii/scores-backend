package sql

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/raphi011/scores"
)

var _ scores.UserRepository = &UserRepository{}

// UserRepository stores users.
type UserRepository struct {
	DB *sqlx.DB
}

// New persists a user and assigns a new id.
func (s *UserRepository) New(user *scores.User) (*scores.User, error) {
	err := insert(s.DB, "user/insert", user)

	return user, mapError(err)
}

// Update updates a user.
func (s *UserRepository) Update(user *scores.User) error {
	if user == nil || user.ID == 0 {
		return errors.New("user must exist")
	}

	result, err := s.DB.Exec(query("user/update"),
		user.ProfileImageURL,
		user.Email,
		user.VolleynetUserID,
		user.VolleynetLogin,
		user.Role,
		user.PasswordInfo.Salt,
		user.PasswordInfo.Hash,
		user.PasswordInfo.Iterations,
		user.ID)

	if err != nil {
		return mapError(err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return mapError(err)
	}

	if rowsAffected != 1 {
		return scores.ErrorNotFound
	}

	return nil
}

// All returns all user's, this is used mainly for testing.
func (s *UserRepository) All() (scores.Users, error) {
	users := scores.Users{}

	rows, err := s.DB.Query(query("user/select-all"))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		u, err := scanUser(rows)

		if err != nil {
			return nil, err
		}

		users = append(users, *u)
	}

	return users, nil
}

// ByID retrieves a user by his/her ID.
func (s *UserRepository) ByID(userID uint) (*scores.User, error) {
	row := s.DB.QueryRow(query("user/select-by-id"), userID)

	return scanUser(row)
}

// ByEmail retrieves a user by his/her email.
func (s *UserRepository) ByEmail(email string) (*scores.User, error) {
	row := s.DB.QueryRow(query("user/select-by-email"), email)

	return scanUser(row)
}

func scanUser(scanner scan) (*scores.User, error) {
	u := &scores.User{}

	err := scanner.Scan(
		&u.ID,
		&u.Email,
		&u.ProfileImageURL,
		&u.PlayerID,
		&u.CreatedAt,
		&u.PasswordInfo.Salt,
		&u.PasswordInfo.Hash,
		&u.PasswordInfo.Iterations,
		&u.VolleynetUserID,
		&u.VolleynetLogin,
		&u.Role,
	)

	return u, mapError(err)
}
