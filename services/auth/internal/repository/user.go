package repository

import (
	"database/sql"
	"errors"

	"syncpad/services/auth/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	_, err := r.db.Exec(
		`INSERT INTO users (id, email, password_hash)
		 VALUES ($1, $2, $3)`,
		user.ID,
		user.Email,
		user.PasswordHash,
	)
	return err
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	row := r.db.QueryRow(
		`SELECT id, email, password_hash, created_at
		 FROM users WHERE email = $1`,
		email,
	)

	var u model.User
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &u, nil
}
