package users

import (
	"database/sql"
	"fmt"
)

type Repository interface {
	Get(userID int64) (*User, error)
}

type userRepository struct {
	db *sql.DB
}

func (r *userRepository) Get(userID int64) (*User, error) {
	if userID == 0 {
		return nil, fmt.Errorf("empty user ID supplied")
	}

	user := User{}
	err := r.db.QueryRow("SELECT id, email FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func New(db *sql.DB) Repository {
	return &userRepository{
		db: db,
	}
}
